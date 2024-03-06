package net

import (
	"io"
	"sync"
	"time"

	"github.com/jkkkls/hjing/rpc"
	"github.com/jkkkls/hjing/utils"
)

// ClientConn 客户端连接信息
type ClientConn struct {
	Mutex      sync.Mutex
	Rwc        io.ReadWriteCloser
	Context    *rpc.Context
	Uid        string
	CommonNode string // 通用逻辑节点
	GameNode   string // 游戏逻辑节点
	T          *time.Timer
	App        string
	Mask       uint8
	RpcClient  *rpc.Client
	Key        []byte
}

type (
	InterceptorNewFunc func(*ClientConn) (*Message, error)                                  // 新连接
	InterceptorReqFunc func(*ClientConn, *Message) (*Message, error)                        // 新消息
	InterceptorRspFunc func(*ClientConn, *Message, []byte, uint16, error) (*Message, error) // 服务端回复
)

func (conn *ClientConn) SendMessage(ICode ICode, msg *Message) {
	utils.Go(
		func() {
			conn.Mutex.Lock()
			defer conn.Mutex.Unlock()
			err := WriteMessage(conn.Rwc, ICode, msg)
			if err != nil {
				// api.Debug("writeMessage error", "remote", conn.Context.Remote, "uid", conn.ConnInfo.Uid, "err", err)
				conn.Rwc.Close()
				return

			}
			if msg.Disconn {
				utils.Debug("断开连接", "remote", conn.Context.Remote)
				conn.Rwc.Close()
			}
		},
	)
}
