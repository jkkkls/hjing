package net

import (
	"bufio"
	"io"
)

type ICode interface {
	Encode(io.Writer, *Message) error
	Decode(io.Reader, *Message) error
}

type Message struct {
	Cmd      uint16 // 消息类型Id
	Seq      uint16 // 消息序号
	Ret      uint16 // 服务端返回结果
	NodeName string // 节点名称
	Params   map[string]any
	Disconn  bool
	Buff     []byte // 消息内容
	IsJson   bool
}

// ReadMessage 通用读取请求接口
func ReadMessage(r *bufio.Reader, ICode ICode) (*Message, error) {
	msg := &Message{}
	err := ICode.Decode(r, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// WriteMessage 通用写请求接口
func WriteMessage(conn io.ReadWriteCloser, ICode ICode, msg *Message) error {
	err := ICode.Encode(conn, msg)
	if err != nil {
		return err
	}
	return nil
}
