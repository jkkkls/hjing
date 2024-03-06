// Copyright 2009 The Go Authors. All rights reserved.

// 基于gorpc代码简单修改，支持内部调用和recover处理
package rpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/jkkkls/hjing/utils"
)

// ServerError represents an error that has been returned from
// the remote side of the RPC connection.
type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

var ErrShutdown = errors.New("connection is shut down")

// CallInfo represents an active RPC.
type CallInfo struct {
	ServiceMethod string         // The name of the service and method to call.
	Args          interface{}    // The argument to the function (*struct).
	Reply         interface{}    // The reply from the function (*struct).
	Error         error          // After completion, the error status.
	Done          chan *CallInfo // Strobes when call is complete.
	Raw           int            // 是否字节流调用.0-否 1-是 2-json
	Ret           uint16         // 返回值
	Conn          *Context       // 连接信息
	NoResp        bool           // 是否不需要回复
}

type ClientOption func(*Client)

func WithName(name string, CloseCallback func(string, error)) ClientOption {
	return func(cli *Client) {
		cli.Name = name
		cli.CloseCallback = CloseCallback
	}
}

func WithTimeout(t int) ClientOption {
	return func(cli *Client) {
		cli.TimeoutSec = t
	}
}

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	Id    int
	codec ClientCodec

	// reqMutex sync.Mutex // protects following
	// request  Request

	mutex    sync.Mutex // protects following
	seq      uint64
	pending  map[uint64]*CallInfo
	closing  bool // user has called Close
	shutdown bool // server has told us to stop

	Name          string
	CloseCallback func(string, error)
	Region        uint32
	TimeoutSec    int
}

// A ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses. The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
type ClientCodec interface {
	// WriteRequest must be safe for concurrent use by multiple goroutines.
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(int32, interface{}) error

	WriteByteRequest(*Request, []byte) error
	ReadByteResponseBody() ([]byte, error)

	Close() error
}

func (client *Client) send(call *CallInfo) {
	// Register this call.
	client.mutex.Lock()
	if client.shutdown || client.closing {
		call.Error = ErrShutdown
		client.mutex.Unlock()
		call.done()
		return
	}
	seq := client.seq
	client.seq++
	if !call.NoResp {
		client.pending[seq] = call
	}
	client.mutex.Unlock()

	// Encode and send the request.
	request := &Request{}
	request.Seq = seq
	request.ServiceMethod = call.ServiceMethod
	request.Conn = call.Conn
	request.NoResp = call.NoResp
	request.Raw = call.Raw

	var err error
	if call.Raw == 0 {
		err = client.codec.WriteRequest(request, call.Args)
	} else {
		// 包括 byte和json格式
		err = client.codec.WriteByteRequest(request, call.Args.([]byte))
	}

	if call.NoResp {
		return
	}
	if err != nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

func (client *Client) input() {
	var err error
	var response Response
	for err == nil {
		response = Response{}
		err = client.codec.ReadResponseHeader(&response)
		if err != nil {
			break
		}
		seq := response.Seq
		raw := response.Raw
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
			// We've got no pending call. That usually means that
			// WriteRequest partially failed, and call was already
			// removed; response is a server telling us about an
			// error reading request body. We should still attempt
			// to read error body, but there's no one to give it to.
			err = client.codec.ReadResponseBody(int32(raw), nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case response.Error != "":
			// We've got an error response. Give this to the request;
			// any subsequent requests will get the ReadResponseBody
			// error if there is one.
			call.Error = ServerError(response.Error)
			err = client.codec.ReadResponseBody(0, nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
			call.done()
		default:
			call.Ret = response.Ret
			if call.Raw == 0 {
				err = client.codec.ReadResponseBody(int32(call.Raw), call.Reply)
			} else {
				// 包括 byte和json格式
				call.Reply, err = client.codec.ReadByteResponseBody()
			}
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// Terminate pending calls.
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()
	if err != io.EOF && closing {
		if client.CloseCallback != nil {
			client.CloseCallback(client.Name, err)
		}
	}
}

func (call *CallInfo) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

// NewClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
// It adds a buffer to the write side of the connection so
// the header and payload are sent as a unit.
func NewClient(conn io.ReadWriteCloser) *Client {
	client := NewPbClientCodec(conn)
	return NewClientWithCodec(client)
}

// NewClientWithCodec is like NewClient but uses the specified
// codec to encode requests and decode responses.
func NewClientWithCodec(codec ClientCodec) *Client {
	client := &Client{
		codec:      codec,
		pending:    make(map[uint64]*CallInfo),
		TimeoutSec: 300,
	}
	utils.Go(func() {
		client.input()
	})

	return client
}

// Dial connects to an RPC server at the specified network address.
func Dial(network, address string, options ...ClientOption) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	client := NewClient(conn)

	for _, o := range options {
		o(client)
	}

	return client, nil
}

// Close calls the underlying codec's Close method. If the connection is already
// shutting down, ErrShutdown is returned.
func (client *Client) Close() error {
	client.mutex.Lock()
	if client.closing {
		client.mutex.Unlock()
		return ErrShutdown
	}
	client.closing = true
	client.mutex.Unlock()
	return client.codec.Close()
}

func (client *Client) IsClose() bool {
	return client.shutdown || client.closing
}

// Go invokes the function asynchronously. It returns the CallInfo structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same CallInfo object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (client *Client) Go(conn *Context, serviceMethod string, args interface{}, reply interface{}, done chan *CallInfo, noResp bool) *CallInfo {
	call := new(CallInfo)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.Raw = 0
	call.Conn = conn
	call.NoResp = noResp
	if done == nil {
		done = make(chan *CallInfo, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			utils.Error("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.send(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (client *Client) Call(conn *Context, serviceMethod string, args interface{}, reply interface{}) (uint16, error) {
	for {
		select {
		case call := <-client.Go(conn, serviceMethod, args, reply, make(chan *CallInfo, 1), false).Done:
			return call.Ret, call.Error
		case <-time.After(time.Duration(client.TimeoutSec) * time.Second):
			return 1, fmt.Errorf("methon[%v] RawCall timeout", serviceMethod)
		}
	}
}

// Send 异步调用，忽略返回值
func (client *Client) Send(conn *Context, serviceMethod string, args interface{}) {
	client.Go(conn, serviceMethod, args, nil, make(chan *CallInfo, 1), true)
}

func (client *Client) RawGo(conn *Context, serviceMethod string, args []byte, done chan *CallInfo, noResp bool, raw int) *CallInfo {
	call := new(CallInfo)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Raw = raw
	call.Conn = conn
	call.NoResp = noResp
	if done == nil {
		done = make(chan *CallInfo, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			utils.Error("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	client.send(call)
	return call
}

// RawCall 传输字节流，用于rpc客户端本身就拿到的是字节流数据
func (client *Client) RawCall(conn *Context, serviceMethod string, args []byte) (uint16, []byte, error) {
	for {
		select {
		case call := <-client.RawGo(conn, serviceMethod, args, make(chan *CallInfo, 1), false, 1).Done:
			if call.Error != nil {
				return call.Ret, nil, call.Error
			}
			return call.Ret, call.Reply.([]byte), nil
		case <-time.After(time.Duration(client.TimeoutSec) * time.Second):
			return 1, nil, fmt.Errorf("methon[%v] RawCall timeout", serviceMethod)
		}
	}
}

// RawSend 传输字节流，用于rpc客户端本身就拿到的是字节流数据
func (client *Client) RawSend(conn *Context, serviceMethod string, args []byte) {
	client.RawGo(conn, serviceMethod, args, make(chan *CallInfo, 1), true, 1)
}

// JsonCall 传输字节流，用于rpc客户端本身就拿到的是字节流数据
func (client *Client) JsonCall(conn *Context, serviceMethod string, args []byte) (uint16, []byte, error) {
	for {
		select {
		case call := <-client.RawGo(conn, serviceMethod, args, make(chan *CallInfo, 1), false, 2).Done:
			if call.Error != nil {
				return call.Ret, nil, call.Error
			}
			return call.Ret, call.Reply.([]byte), nil
		case <-time.After(time.Duration(client.TimeoutSec) * time.Second):
			return 1, nil, fmt.Errorf("methon[%v] JsonCall timeout", serviceMethod)
		}
	}
}

// JsonSend 传输字节流，用于rpc客户端本身就拿到的是字节流数据
func (client *Client) JsonSend(conn *Context, serviceMethod string, args []byte) {
	client.RawGo(conn, serviceMethod, args, make(chan *CallInfo, 1), true, 2)
}
