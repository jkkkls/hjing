package net

import (
	"bufio"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jkkkls/hjing/rpc"
	"github.com/jkkkls/hjing/utils"
	"google.golang.org/protobuf/proto"
)

type Gater struct {
	Id        int
	IDCounter uint64 //连接id
	SyncCall  bool
	Count     int32
	Ids       utils.Map[uint64, *ClientConn]
	ICode     ICode

	NewFuncs []InterceptorNewFunc //新连接拦截器列表
	ReqFuncs []InterceptorReqFunc //新消息拦截器列表
	RspFuncs []InterceptorRspFunc //服务端回复拦截器列表
	DelFuncs []InterceptorNewFunc //连接断开拦截器列表
}

type GaterOption func(*Gater)

var (
	gater *Gater
)

// NewGater 初始化网关
func NewGater(syncCall bool, icode ICode, options ...GaterOption) *Gater {
	gater = &Gater{
		IDCounter: 1,
		SyncCall:  syncCall,
		ICode:     icode,
	}

	for _, o := range options {
		o(gater)
	}
	return gater
}

// RegisterNewFuncs 注册新连接拦截器
func (g *Gater) RegisterNewFuncs(f InterceptorNewFunc) {
	g.NewFuncs = append(g.NewFuncs, f)
}

// RegisterReqFuncs 注册新请求拦截器
func (g *Gater) RegisterReqFuncs(f InterceptorReqFunc) {
	g.ReqFuncs = append(g.ReqFuncs, f)
}

// RegisterRspFuncs 注册收到响应拦截器
func (g *Gater) RegisterRspFuncs(f InterceptorRspFunc) {
	g.RspFuncs = append(g.RspFuncs, f)
}

// RegisterDelFuncs 注册收到连接断开拦截器
func (g *Gater) RegisterDelFuncs(f InterceptorNewFunc) {
	g.DelFuncs = append(g.DelFuncs, f)
}

// newID 生成一个连接id，前两个字节保存节点id，后四个字节递增
func (g *Gater) newID() uint64 {
	for {
		id := atomic.AddUint64(&g.IDCounter, 1)
		id = id&0xFFFFFFFF | (rpc.NodeIntance.Config.Id&0xFFFF)<<32
		if _, ok := g.Ids.Load(id); !ok {
			return id
		}
	}
}

// GetConnByID 根据id返回连接
func (g *Gater) GetConnByID(id uint64) *ClientConn {
	v, ok := g.Ids.Load(id)
	if !ok {
		return nil
	}
	return v
}

// NewClientConn 初始化一个新链接
func (g *Gater) NewClientConn(remote string, rwc io.ReadWriteCloser) *ClientConn {
	conn := &ClientConn{
		Rwc: rwc,
		T:   time.NewTimer(30 * time.Second),
		Context: &rpc.Context{
			Remote:   remote,
			Id:       g.newID(),
			GateName: rpc.NodeIntance.Config.Nodename,
		},
	}

	atomic.AddInt32(&g.Count, 1)
	g.Ids.Store(conn.Context.Id, conn)
	return conn
}

// DelClientConn 销毁一个链接
func (g *Gater) DelClientConn(conn *ClientConn) {
	atomic.AddInt32(&g.Count, -1)
	g.Ids.Delete(conn.Context.Id)
}

// GetCount ...
func (g *Gater) GetCount() int32 {
	return atomic.LoadInt32(&g.Count)
}

// handleConn send back everything it received
func (g *Gater) handleConn(rwc io.ReadWriteCloser, remote string, connType string) {
	r := bufio.NewReader(rwc)
	remote = strings.Split(remote, ":")[0]
	conn := g.NewClientConn(remote, rwc)
	defer func() {
		utils.Trace("连接断开", "remote", remote, "connType", connType, "id", conn.Context.Id)
		g.DelClientConn(conn)
		rwc.Close()

		for _, v := range g.DelFuncs {
			v(conn)
		}
	}()
	utils.Trace("新连接", "remote", remote, "connType", connType, "id", conn.Context.Id)

	//新连接拦截器
	for _, v := range g.NewFuncs {
		if _, err := v(conn); err != nil {
			return
		}
	}

	//连接心跳定时器
	utils.Submit(func() {
		for range conn.T.C {
			rwc.Close()
			conn.T.Stop()
			return
		}
	})

	//读取数据逻辑
	for {
	LOOP:
		//读取消息
		msg, err := ReadMessage(r, g.ICode)
		if err != nil {
			utils.Debug("readMessage error", "remote", conn.Context.Remote, "error", err)
			return
		}

		//新请求拦截器
		for _, v := range g.ReqFuncs {
			if ret, err := v(conn, msg); err != nil {
				conn.SendMessage(g.ICode, ret)
				goto LOOP
			}
		}

		//重置定时器
		conn.T.Reset(120 * time.Second)
		//转发消息
		cmd := msg.Cmd
		nodeName := msg.NodeName
		serverMethod := rpc.FindServiceMethod("", cmd)
		if serverMethod != "" {
			f := func() {
				var (
					ret     uint16
					rspBuff []byte
					retErr  error
					rpcConn = proto.Clone(conn.Context).(*rpc.Context)
				)
				rpcConn.CallId = rpc.NewId(0)

				if nodeName != "" {
					if msg.IsJson {
						ret, rspBuff, retErr = rpc.NodeJsonCallWithConn(rpcConn, nodeName, serverMethod, msg.Buff)
					} else {
						ret, rspBuff, retErr = rpc.NodeRawCallWithConn(rpcConn, nodeName, serverMethod, msg.Buff)
					}
				} else {
					if msg.IsJson {
						ret, rspBuff, retErr = rpc.JsonCall(rpcConn, serverMethod, msg.Buff)
					} else {
						ret, rspBuff, retErr = rpc.RawCall(rpcConn, serverMethod, msg.Buff)
					}
				}

				if retErr != nil && ret == 0 {
					ret = 1
				}

				//响应拦截器
				for _, v := range g.RspFuncs {
					v(conn, msg, rspBuff, ret, retErr)
				}

				conn.SendMessage(g.ICode, &Message{
					Cmd:  msg.Cmd,
					Seq:  msg.Seq,
					Ret:  uint16(ret),
					Buff: rspBuff,
				})
			}
			if g.SyncCall {
				utils.ProtectCall(f, nil)
			} else {
				utils.Submit(f)
			}
		} else {
			conn.SendMessage(g.ICode, &Message{
				Cmd: cmd,
				Seq: msg.Seq,
				Ret: 2,
			})
			utils.Info("connect cmd error", "cmd", msg.Cmd)
		}
	}
}

// ForConn 遍历所有链接
func (g *Gater) ForConn(f func(conn *ClientConn)) {
	var ids []uint64
	g.Ids.Range(func(key uint64, _ *ClientConn) bool {
		ids = append(ids, key)
		return true
	})
	for i := 0; i < len(ids); i++ {
		wsconn := g.GetConnByID(ids[i])
		if wsconn == nil {
			continue
		}
		f(wsconn)
	}
}
