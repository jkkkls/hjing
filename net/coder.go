package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// ConnCoder 简单的长链接消息编解码器
// 消息格式: G + Cmd + BuffLen + Buff
// 长度：1 + 2 + 2 + BuffLen
type ConnCoder struct{}

var ErrInvalidMessage = errors.New("error message flag")

func (c *ConnCoder) Encode(w io.Writer, msg *Message) error {
	var b bytes.Buffer

	b.WriteByte('G')
	binary.Write(&b, binary.BigEndian, msg.Cmd)
	binary.Write(&b, binary.BigEndian, uint16(len(msg.Buff)))
	b.Write(msg.Buff)

	_, err := w.Write(b.Bytes())

	return err
}

func (c *ConnCoder) Decode(r io.Reader, msg *Message) error {
	var head [5]byte

	_, err := io.ReadFull(r, head[:])
	if err != nil {
		return err
	}

	if head[0] != 'G' {
		return ErrInvalidMessage
	}

	msg.Cmd = binary.BigEndian.Uint16(head[1:3])
	buffLen := binary.BigEndian.Uint16(head[3:5])

	msg.Buff = make([]byte, buffLen)
	_, err = io.ReadFull(r, msg.Buff)
	if err != nil {
		return err
	}

	return nil
}
