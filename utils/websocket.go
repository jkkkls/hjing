package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/seveye/goms/util/bytes_cache"
	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

// Message 请求，协议头固定12字节：
// 请求，协议头固定12字节：
//     是否加密（1字节）| 长度（2字节）| cmd（2字节）| 序列号（2字节）| 俱乐部id（4字节）| 预留（1字节）| 加密数据（n字节）

// 响应，协议头固定12字节：
//
//	是否加密（1字节）| 长度（2字节）| cmd（2字节）| 序列号（2字节）| 返回值（4字节）| 预留（1字节）| 加密数据（n字节）
type RequestMessage struct {
	Crype       uint8
	Length      uint16
	Cmd         uint16
	Seq         uint16
	RetOrClubId uint32
	Buff        []byte
	Pointer     proto.Message
	Disconn     bool
}

// 通讯密钥
var AesKey = []byte("12345678901234567890123456789012")

// WsConn websocket封装
type WsConn struct {
	Conn *websocket.Conn
	Buff bytes.Buffer
}

func (c *WsConn) Write(b []byte) (int, error) {
	err := c.Conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *WsConn) Read(b []byte) (int, error) {
	if c.Buff.Len() > 0 {

		return c.Buff.Read(b)
	}
	_, buff, err := c.Conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	c.Buff.Write(buff)
	return c.Buff.Read(b)
}

func (c *WsConn) Close() error {
	return c.Conn.Close()
}

func (s *WsConn) SetWriteDeadline(t time.Time) error {
	return s.Conn.SetWriteDeadline(t)
}

// ReadMessage 通用读取请求接口
func ReadMessage(r *bufio.Reader) (*RequestMessage, error) {
	req := &RequestMessage{}
	var (
		header [12]byte
	)
	_, err := io.ReadFull(r, header[:])
	if err != nil {
		return nil, err
	}

	//是否加密（1字节）| 长度（2字节）| cmd（2字节）| 序列号（2字节）| 俱乐部id（4字节）| 预留（1字节）| 加密数据（n字节）
	req.Crype = uint8(header[0])
	req.Length = binary.BigEndian.Uint16(header[1:3])
	req.Cmd = binary.BigEndian.Uint16(header[3:5])
	req.Seq = binary.BigEndian.Uint16(header[5:7])
	req.RetOrClubId = binary.BigEndian.Uint32(header[7:11])

	if req.Length == 0 {
		return req, nil
	}

	buff := bytes_cache.Get(int(req.Length))
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return nil, err
	}

	//解密
	if req.Crype == 1 {
		req.Buff, err = AesDecrypt(buff, AesKey)
		if err != nil {
			return nil, err
		}
	} else {
		req.Buff = buff
	}

	return req, nil
}

// WriteMessage 通用写请求接口
func WriteMessage(conn io.ReadWriteCloser, req *RequestMessage) error {
	var err error
	if len(req.Buff) == 0 && req.Pointer != nil {
		msg, ok := req.Pointer.(proto.Message)
		if !ok {
			return fmt.Errorf("%T does not implement proto.Message", req.Pointer)
		}
		req.Buff, err = proto.Marshal(msg)
		if err != nil {
			return err
		}
	}

	buff := req.Buff
	if req.Crype == 1 {
		buff, err = AesEncrypt(req.Buff, AesKey)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer
	b.WriteByte(byte(req.Crype))
	binary.Write(&b, binary.BigEndian, uint16(len(buff)))
	binary.Write(&b, binary.BigEndian, req.Cmd)
	binary.Write(&b, binary.BigEndian, req.Seq)
	binary.Write(&b, binary.BigEndian, uint32(req.RetOrClubId))
	b.WriteByte(0) //预留
	b.Write(buff)

	_, err = conn.Write(b.Bytes())
	return err
}
