// Copyright 2009 The Go Authors. All rights reserved.

// 基于gorpc代码简单修改，支持内部调用和recover处理
package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/jkkkls/hjing/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfUint16 = reflect.TypeOf(uint16(0))

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	numCalls   uint
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

// Request is a header written before every RPC call. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	ServiceMethod string // format: "Service.Method"
	Seq           uint64 // sequence number chosen by client
	NoResp        bool
	Conn          *Context //连接信息
	next          *Request // for free list in Server
	Raw           int
}

// Response is a header written before every RPC return. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Response struct {
	ServiceMethod string // echoes that of the Request
	Seq           uint64 // echoes that of the request
	Error         string // error, if any.
	Ret           uint16
	next          *Response // for free list in Server
	Raw           int
}

type RpcCallBack func(conn *Context, methon string, req, rsp proto.Message, ret uint16, err error, cost time.Duration)

// Server represents an RPC Server.
type Server struct {
	serviceMap sync.Map   // map[string]*service
	reqLock    sync.Mutex // protects freeReq
	freeReq    *Request
	respLock   sync.Mutex // protects freeResp
	freeResp   *Response

	RpcCallBack RpcCallBack
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//   - exported method of exported type
//   - two arguments, both of exported type
//   - the second argument is a pointer
//   - one return value, of type error
//
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (server *Server) Register(rcvr interface{}) error {
	return server.register(rcvr, "", false)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (server *Server) RegisterName(name string, rcvr interface{}) error {
	return server.register(rcvr, name, true)
}

func (server *Server) register(rcvr interface{}, name string, useName bool) error {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		return errors.New(s)
	}
	if !isExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		return errors.New(s)
	}
	s.name = sname

	// Install the methods
	s.method = suitableMethods(name, s.typ, false)

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(name, reflect.PtrTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(name string, typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 4 {
			if reportErr {
				log.Println("method", mname, "has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Println(mname, "argument type not exported:", argType)
			}
			continue
		}
		// Thrid arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr && replyType.Kind() != reflect.Interface {
			if reportErr {
				log.Println("method", mname, "reply type not a pointer:", replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Println("method", mname, "reply type not exported:", replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 2 {
			if reportErr {
				log.Println("method", mname, "has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(1); returnType != typeOfError {
			if reportErr {
				log.Println("method", mname, "returns", returnType.String(), "not error")
			}
			continue
		}
		if returnType := mtype.Out(0); returnType != typeOfUint16 {
			if reportErr {
				log.Println("method", mname, "returns", returnType.String(), "not error")
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
		log.Printf("[%v]register method: %v\n", name, mname)
	}
	return methods
}

// A value sent as a placeholder for the server's response value when the server
// receives an invalid request. It is never decoded by the client since the Response
// contains an error when it is used.
var invalidRequest = &NullMessage{}

func (server *Server) sendResponse(sending *sync.Mutex, req *Request, reply interface{}, codec ServerCodec, ret uint16, errmsg string) {
	resp := server.getResponse()
	resp.Ret = ret
	// Encode the response header
	resp.ServiceMethod = req.ServiceMethod
	if errmsg != "" {
		resp.Error = errmsg
		reply = invalidRequest
	}
	resp.Seq = req.Seq
	resp.Raw = req.Raw
	sending.Lock()
	err := codec.WriteResponse(resp, reply)
	if err != nil {
		log.Println("rpc: writing response:", err)
	}
	sending.Unlock()
	server.freeResponse(resp)
}

func (m *methodType) NumCalls() (n uint) {
	m.Lock()
	n = m.numCalls
	m.Unlock()
	return n
}

func splitServiceMethod(serviceMethod string) (string, string) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		return "", ""
	}
	return serviceMethod[:dot], serviceMethod[dot+1:]
}

func (server *Server) RawCall(conn *Context, serviceMethod string, reqBuff []byte, js bool) (uint16, []byte, error) {
	serviceName, methodName := splitServiceMethod(serviceMethod)
	if serviceName == "" || methodName == "" {
		err := errors.New("rpc: service/method request ill-formed: " + serviceMethod)
		return 0, nil, err
	}

	// Look up the request.
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		return 0, nil, errors.New("rpc: can't find service " + serviceName)
	}
	svc := svci.(*service)
	mtype := svc.method[methodName]
	if mtype == nil {
		return 0, nil, errors.New("rpc1: can't find method " + methodName)
	}
	argv := reflect.New(mtype.ArgType.Elem())
	replyv := reflect.New(mtype.ReplyType.Elem())

	if pb, ok := argv.Interface().(proto.Message); ok {
		if js {
			protojson.Unmarshal(reqBuff, pb)
		} else {
			proto.Unmarshal(reqBuff, pb)
		}
	} else {
		return 0, nil, fmt.Errorf("does not implement proto.Message")
	}

	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	function := mtype.method.Func

	// Invoke the method, providing a new value for the reply.
	start := time.Now()
	returnValues := function.Call([]reflect.Value{svc.rcvr, reflect.ValueOf(conn), argv, replyv})
	errInter := returnValues[1].Interface()
	ret := returnValues[0].Interface().(uint16)

	if server.RpcCallBack != nil {
		utils.Submit(func() {
			server.RpcCallBack(conn, serviceMethod, argv.Interface().(proto.Message), replyv.Interface().(proto.Message),
				ret, nil, time.Since(start))
		})
	}

	if errInter != nil {
		return 0, nil, errInter.(error)
	}

	if pb, ok := replyv.Interface().(proto.Message); ok {
		if js {
			buf, err := json.Marshal(pb)
			if err != nil {
				return 0, nil, err
			}
			return ret, buf, nil
		}

		buf, err := proto.Marshal(pb)
		if err != nil {
			return 0, nil, err
		}
		return ret, buf, nil
	}

	return 0, nil, fmt.Errorf("does not implement proto.Message")
}

// InternalCall 内部调用
func (server *Server) InternalCall(conn *Context, serviceMethod string, req proto.Message, rsp proto.Message) (uint16, error) {
	serviceName, methodName := splitServiceMethod(serviceMethod)
	if serviceName == "" || methodName == "" {
		err := errors.New("rpc: service/method request ill-formed: " + serviceMethod)
		return 0, err
	}

	if rsp == nil {
		rsp = &NullMessage{}
	}

	// Look up the request.
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		return 0, errors.New("rpc: can't find service " + serviceName)
	}
	svc := svci.(*service)
	mtype := svc.method[methodName]
	if mtype == nil {
		return 0, errors.New("rpc2: can't find method " + methodName)
	}

	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	function := mtype.method.Func

	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{svc.rcvr, reflect.ValueOf(conn), reflect.ValueOf(req), reflect.ValueOf(rsp)})
	errInter := returnValues[1].Interface()
	if errInter != nil {
		return 0, errInter.(error)
	}
	ret := returnValues[0].Interface().(uint16)
	return ret, nil
}

func (s *service) call(server *Server, sending *sync.Mutex, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {
	var ret uint16 = 1
	errmsg := "server exception"
	defer utils.Recover()
	defer func() {

		if req.NoResp {
			server.freeRequest(req)
			return
		}

		server.sendResponse(sending, req, replyv.Interface(), codec, ret, errmsg)
		server.freeRequest(req)
	}()

	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	function := mtype.method.Func

	//
	if replyv.Kind() == reflect.Invalid {
		replyv = reflect.ValueOf(&NullMessage{})
	}

	// Invoke the method, providing a new value for the reply.
	start := time.Now()
	utils.SetCallId(req.Conn.GetUid(), req.Conn.GetCallId())
	returnValues := function.Call([]reflect.Value{s.rcvr, reflect.ValueOf(req.Conn), argv, replyv})
	utils.UnLockCallId()

	// The return value for the method is an error.
	errInter := returnValues[1].Interface()
	if errInter != nil {
		errmsg = errInter.(error).Error()
	} else {
		errmsg = ""
	}
	ret = returnValues[0].Interface().(uint16)

	if server.RpcCallBack != nil {
		utils.Submit(func() {
			server.RpcCallBack(req.Conn, req.ServiceMethod, argv.Interface().(proto.Message), replyv.Interface().(proto.Message),
				ret, errors.New(errmsg), time.Since(start))
		})
	}
}

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection. To use an alternate codec, use ServeCodec.
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	server.ServeCodec(NewPbServerCodec(conn))
}

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
func (server *Server) ServeCodec(codec ServerCodec) {
	sending := new(sync.Mutex)
	for {
		service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
		if err != nil {
			// if err != io.EOF {
			// 	log.Println("rpc4:", err)
			// }
			if !keepReading {
				break
			}
			// send a response if we actually managed to read a header.
			if req != nil {
				if !req.NoResp {
					server.sendResponse(sending, req, invalidRequest, codec, 0, err.Error())
				}
				server.freeRequest(req)
			}
			continue
		}
		utils.Submit(func() {
			service.call(server, sending, mtype, req, argv, replyv, codec)
		})
	}
	codec.Close()
}

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.
func (server *Server) ServeRequest(codec ServerCodec) error {
	sending := new(sync.Mutex)
	service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
	if err != nil {
		if !keepReading {
			return err
		}
		// send a response if we actually managed to read a header.
		if req != nil {
			if !req.NoResp {
				server.sendResponse(sending, req, invalidRequest, codec, 0, err.Error())
			}
			server.freeRequest(req)
		}
		return err
	}
	service.call(server, sending, mtype, req, argv, replyv, codec)
	return nil
}

func (server *Server) getRequest() *Request {
	server.reqLock.Lock()
	req := server.freeReq
	if req == nil {
		req = new(Request)
	} else {
		server.freeReq = req.next
		*req = Request{}
	}
	server.reqLock.Unlock()
	return req
}

func (server *Server) freeRequest(req *Request) {
	server.reqLock.Lock()
	req.next = server.freeReq
	server.freeReq = req
	server.reqLock.Unlock()
}

func (server *Server) getResponse() *Response {
	server.respLock.Lock()
	resp := server.freeResp
	if resp == nil {
		resp = new(Response)
	} else {
		server.freeResp = resp.next
		*resp = Response{}
	}
	server.respLock.Unlock()
	return resp
}

func (server *Server) freeResponse(resp *Response) {
	server.respLock.Lock()
	resp.next = server.freeResp
	server.freeResp = resp
	server.respLock.Unlock()
}

func (server *Server) readRequest(codec ServerCodec) (service *service, mtype *methodType, req *Request, argv, replyv reflect.Value, keepReading bool, err error) {
	service, mtype, req, keepReading, err = server.readRequestHeader(codec)
	if err != nil {
		if !keepReading {
			return
		}

		// discard body
		codec.ReadRequestBody(int32(req.Raw), nil)
		return
	}

	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = codec.ReadRequestBody(int32(req.Raw), argv.Interface()); err != nil {
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}

	if mtype.ReplyType.Kind() != reflect.Interface {
		replyv = reflect.New(mtype.ReplyType.Elem())
		switch mtype.ReplyType.Elem().Kind() {
		case reflect.Map:
			replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
		case reflect.Slice:
			replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
		}
	}

	return
}

func (server *Server) readRequestHeader(codec ServerCodec) (svc *service, mtype *methodType, req *Request, keepReading bool, err error) {
	// Grab the request header.
	req = server.getRequest()
	err = codec.ReadRequestHeader(req)
	if err != nil {
		req = nil
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		err = errors.New("rpc: server cannot decode request: " + err.Error())
		return
	}

	// We read the header successfully. If we see an error now,
	// we can still recover and move on to the next request.
	keepReading = true

	dot := strings.LastIndex(req.ServiceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc: service/method request ill-formed: " + req.ServiceMethod)
		return
	}
	serviceName := req.ServiceMethod[:dot]
	methodName := req.ServiceMethod[dot+1:]

	// Look up the request.
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc: can't find service " + req.ServiceMethod)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc3: can't find method " + req.ServiceMethod)
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection. Accept blocks until the listener
// returns a non-nil error. The caller typically invokes Accept in a
// go statement.
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			return
		}
		// log.Println("new rpc client, remote:", conn.RemoteAddr().String())
		utils.Submit(func() {
			server.ServeConn(conn)
		})
	}
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

// A ServerCodec implements reading of RPC requests and writing of
// RPC responses for the server side of an RPC session.
// The server calls ReadRequestHeader and ReadRequestBody in pairs
// to read requests from the connection, and it calls WriteResponse to
// write a response back. The server calls Close when finished with the
// connection. ReadRequestBody may be called with a nil
// argument to force the body of the request to be read and discarded.
type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(int32, interface{}) error
	// WriteResponse must be safe for concurrent use by multiple goroutines.
	WriteResponse(*Response, interface{}) error

	Close() error
}

// ServeConn runs the DefaultServer on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection. To use an alternate codec, use ServeCodec.
func ServeConn(conn io.ReadWriteCloser) {
	DefaultServer.ServeConn(conn)
}

// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
func ServeCodec(codec ServerCodec) {
	DefaultServer.ServeCodec(codec)
}

// ServeRequest is like ServeCodec but synchronously serves a single request.
// It does not close the codec upon completion.
func ServeRequest(codec ServerCodec) error {
	return DefaultServer.ServeRequest(codec)
}

// Accept accepts connections on the listener and serves requests
// to DefaultServer for each incoming connection.
// Accept blocks; the caller typically invokes it in a go statement.
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }
