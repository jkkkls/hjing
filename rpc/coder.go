package rpc

// Copyright 2017 guangbo. All rights reserved.

// rpc编码解码模块

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/jkkkls/hjing/utils/bytes_cache"
	"google.golang.org/protobuf/proto"
)

// tooBig 内部通讯最大消息长度 1G
const tooBig = 1 << 30

var errBadCount = errors.New("invalid message length")

func writeFrame(w *bytes.Buffer, buf []byte) error {
	l := len(buf)
	if l >= tooBig {
		return errBadCount
	}

	lenBuf := bytes_cache.Get(4)
	defer bytes_cache.Put(lenBuf)
	binary.BigEndian.PutUint32(lenBuf[:], uint32(l))
	_, err := w.Write(lenBuf[:])
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func readFrame(r *bufio.Reader) ([]byte, error) {
	header := bytes_cache.Get(4)
	defer bytes_cache.Put(header)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}

	l := binary.BigEndian.Uint32(header)
	if l >= tooBig {
		return nil, errBadCount
	}

	buff := bytes_cache.Get(int(l))
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func encode(w *bytes.Buffer, raw int32, m interface{}) error {
	if pb, ok := m.(proto.Message); ok {
		if raw == 2 {
			buf, err := json.Marshal(pb)
			if err != nil {
				return err
			}
			return writeFrame(w, buf)
		}

		buf, err := proto.Marshal(pb)
		if err != nil {
			return err
		}
		return writeFrame(w, buf)
	}
	return fmt.Errorf("%T does not implement proto.Message", m)
}

func decode(r *bufio.Reader, raw int32, m interface{}) error {
	buff, err := readFrame(r)
	if err != nil {
		return err
	}

	if m == nil {
		return nil
	}

	if raw == 2 {
		return json.Unmarshal(buff, m)
		// return jsonpb.UnmarshalString(string(buff), m.(proto.Message))
	}

	return proto.Unmarshal(buff, m.(proto.Message))
}

type PbClientCodec struct {
	c  io.ReadWriteCloser
	w  *bufio.Writer
	r  *bufio.Reader
	mu sync.Mutex
}

func NewPbClientCodec(rwc io.ReadWriteCloser) ClientCodec {
	return &PbClientCodec{
		r: bufio.NewReaderSize(rwc, 4096),
		w: bufio.NewWriterSize(rwc, 4096),
		c: rwc,
	}
}

func (c *PbClientCodec) WriteRequest(r *Request, body interface{}) error {
	req := reqHeaderPool.Get().(*ReqHeader)
	defer reqHeaderPool.Put(req)
	req.Reset()
	req.Method = r.ServiceMethod
	req.Seq = r.Seq
	req.NoResp = r.NoResp
	req.Raw = int32(r.Raw)
	if r.Conn != nil {
		req.Context = &Context{
			GateName: r.Conn.GateName,
			Remote:   r.Conn.Remote,
			Id:       r.Conn.Id,
			CallId:   r.Conn.CallId,
			Kvs:      r.Conn.Kvs,
			Ps:       r.Conn.Ps,
		}
	}

	var buff bytes.Buffer
	err := encode(&buff, 0, req)
	if err != nil {
		return err
	}
	if err = encode(&buff, 0, body); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.w.Write(buff.Bytes())
	if err != nil {
		return err
	}
	return c.w.Flush()
}

func (c *PbClientCodec) ReadResponseHeader(r *Response) error {
	resp := rspHeaderPool.Get().(*RspHeader)
	defer rspHeaderPool.Put(resp)
	resp.Reset()
	err := decode(c.r, 0, resp)
	if err != nil {
		return err
	}
	r.ServiceMethod = resp.Method
	r.Seq = resp.Seq
	r.Error = resp.Error
	r.Ret = uint16(resp.Ret)
	return nil
}

func (c *PbClientCodec) ReadResponseBody(raw int32, body interface{}) error {
	return decode(c.r, raw, body)
}

func (c *PbClientCodec) WriteByteRequest(r *Request, buf []byte) error {
	req := reqHeaderPool.Get().(*ReqHeader)
	defer reqHeaderPool.Put(req)
	req.Reset()
	req.Method = r.ServiceMethod
	req.Seq = r.Seq
	req.Raw = int32(r.Raw)
	if r.Conn != nil {
		req.Context = &Context{
			GateName: r.Conn.GateName,
			Remote:   r.Conn.Remote,
			Id:       r.Conn.Id,
			CallId:   r.Conn.CallId,
			Kvs:      r.Conn.Kvs,
			Ps:       r.Conn.Ps,
		}
	}

	var buff bytes.Buffer
	err := encode(&buff, 0, req)
	if err != nil {
		return err
	}
	if err = writeFrame(&buff, buf); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.w.Write(buff.Bytes())
	if err != nil {
		return err
	}
	return c.w.Flush()
}

func (c *PbClientCodec) ReadByteResponseBody() ([]byte, error) {
	return readFrame(c.r)
}

func (c *PbClientCodec) Close() error {
	return c.c.Close()
}

type PbServerCodec struct {
	mu   sync.Mutex // exclusive writer lock
	req  ReqHeader
	resp RspHeader
	w    *bufio.Writer
	r    *bufio.Reader

	c io.Closer
}

func NewPbServerCodec(rwc io.ReadWriteCloser) ServerCodec {
	return &PbServerCodec{
		r: bufio.NewReaderSize(rwc, 4096),
		w: bufio.NewWriterSize(rwc, 4096),
		c: rwc,
	}
}

func (c *PbServerCodec) WriteResponse(resp *Response, body interface{}) error {
	var cresp RspHeader
	cresp.Method = resp.ServiceMethod
	cresp.Seq = resp.Seq
	cresp.Error = resp.Error
	cresp.Ret = uint32(resp.Ret)
	cresp.Raw = int32(resp.Raw)

	var buff bytes.Buffer
	err := encode(&buff, 0, &cresp)
	if err != nil {
		return err
	}
	if err = encode(&buff, cresp.Raw, body); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.w.Write(buff.Bytes())
	if err != nil {
		return err
	}
	return c.w.Flush()
}

func (c *PbServerCodec) ReadRequestHeader(req *Request) error {
	c.req.Reset()

	err := decode(c.r, 0, &c.req)
	if err != nil {
		return err
	}

	req.ServiceMethod = c.req.Method
	req.Seq = c.req.Seq
	req.NoResp = c.req.NoResp
	req.Raw = int(c.req.Raw)
	if c.req.Context != nil {
		req.Conn = &Context{
			GateName: c.req.Context.GateName,
			Remote:   c.req.Context.Remote,
			Id:       c.req.Context.Id,
			CallId:   c.req.Context.CallId,
			Kvs:      c.req.Context.Kvs,
			Ps:       c.req.Context.Ps,
		}
	}
	return nil
}

func (c *PbServerCodec) ReadRequestBody(raw int32, body interface{}) error {
	return decode(c.r, raw, body)
}

func (c *PbServerCodec) Close() error { return c.c.Close() }
