package net

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"
)

// Handshake 客户端dh算法握手逻辑
func ClientHandshake(conn io.ReadWriteCloser) ([]byte, error) {
	var err error
	// 生成p,g,a
	p, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	g, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	//私钥a不需要发送
	a, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	// 计算A
	A := big.NewInt(0)
	A.Exp(g, a, p)

	// 发送p,g,A
	err = writeBigInts(conn, p, g, A)
	if err != nil {
		return nil, err
	}

	// 读取返回的B
	b, err := readBigInt(conn)
	if err != nil {
		return nil, err
	}

	// 计算s
	s := big.NewInt(0)
	s.Exp(b, a, p)
	key := sha256.Sum256(s.Bytes())
	return key[:], nil
}

// Handshake 服务端dh算法握手逻辑
func Handshake(conn io.ReadWriteCloser) ([]byte, error) {

	var (
		err error
		p   = big.NewInt(0)
		g   = big.NewInt(0)
		a   = big.NewInt(0)
		b   = big.NewInt(0)
		s   = big.NewInt(0)
	)

	// 读取客户端p,g,a
	if p, err = readBigInt(conn); err != nil {
		return nil, err
	}

	if g, err = readBigInt(conn); err != nil {
		return nil, err
	}

	if a, err = readBigInt(conn); err != nil {
		return nil, err
	}

	// 生成b
	if b, err = rand.Prime(rand.Reader, 1024); err != nil {
		return nil, err
	}

	// 计算返回的B
	B := big.NewInt(0)
	B.Exp(g, b, p)
	err = writeBigInt(conn, B)
	if err != nil {
		return nil, err
	}

	// 计算s
	s.Exp(a, b, p)
	key := sha256.Sum256(s.Bytes())
	return key[:], nil
}

func writeBigInts(conn io.ReadWriteCloser, is ...*big.Int) error {
	for _, i := range is {
		buff := i.Bytes()
		if err := writeBytes(conn, buff); err != nil {
			return err
		}
	}
	return nil
}

func writeBigInt(conn io.ReadWriteCloser, i *big.Int) error {
	buff := i.Bytes()
	if err := writeBytes(conn, buff); err != nil {
		return err
	}
	return nil
}

func writeBytes(conn io.ReadWriteCloser, buff []byte) error {
	if _, err := conn.Write(buff); err != nil {
		return err
	}
	return nil
}

func readBigInt(conn io.ReadWriteCloser) (*big.Int, error) {
	buff := make([]byte, 128)
	if _, err := io.ReadFull(conn, buff); err != nil {
		return nil, err
	}
	i := new(big.Int)
	i.SetBytes(buff)
	return i, nil
}
