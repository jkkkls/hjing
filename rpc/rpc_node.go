// Copyright 2017 guangbo. All rights reserved.

//
//服务管理
//

package rpc

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/utils"

	"google.golang.org/protobuf/proto"
)

// GxService 服务接口
type GxService interface {
	GetName() string
	Run() error
	Exit()
	NodeConn(string)
	NodeClose(string)
	// 服务事件接口
	OnEvent(string, ...interface{})
}

type NodeConfig struct {
	Id       uint64           // 节点id
	Nodename string           // 节点名称
	Nodetype string           // 节点类型，rpc调用时候可以根据类型调用
	Set      string           // 节点集名称，有时候需要对节点进行分组
	Host     string           // 节点IP地址，需要提供内网地址，注册到Watch节点，以实现服务集群
	Port     int              // 节点rpc端口
	Region   uint32           // 地区
	Cmds     map[int32]string // 节点支持cmd的，需要使用protobuf定义的cmd
	HttpPort int              // 节点http端口
}

// GxNodeConn 微服务节点信息
type GxNodeConn struct {
	Info  *config.NodeInfo
	Conns []*Client

	UseTime []int64
	Mutex   sync.Mutex

	Close bool
}

// Rand2Num 生成两个不同的随机数
func Rand2Num(n int) (int, int) {
	if n == 1 {
		return 0, 0
	} else if n == 2 {
		return 0, 1
	}

	x := []int{}
	for i := 0; i < n; i++ {
		x = append(x, i)
	}
	index := rand.IntN(len(x))
	i1 := x[index]
	x = append(x[:index], x[index+1:]...)
	i2 := x[rand.IntN(len(x))]
	return i1, i2
}

// getConn 随机选择两个随机数，再选择一个负载少的连接
func (conn *GxNodeConn) getConn() *Client {
	n := len(conn.Conns)
	i1, i2 := Rand2Num(n)
	conn.Mutex.Lock()
	index := utils.If(conn.UseTime[i1] > conn.UseTime[i2], i2, i1)
	conn.UseTime[index]++

	// 重置
	if conn.UseTime[index] > 1>>30 {
		for i := 0; i < n; i++ {
			conn.UseTime[i] = 0
		}
	}
	conn.Mutex.Unlock()

	return conn.Conns[index]
}

// GxNode 本节点信息
type GxNode struct {
	IRegister IRegister
	Config    *NodeConfig
	Server    *Server                        // 自己rpc服务端
	Services  utils.Map[string, GxService]   // 自己的服务列表
	RpcClient utils.Map[string, *GxNodeConn] // 到其他节点的连接

	Mutex       sync.Mutex
	ServiceNode map[string]map[string]bool // 服务对应节点名

	// 全局变量
	IDGen *utils.IDGen

	//
	MockData    map[string]*MockRsp
	RpcCallBack RpcCallBack
	Node        *config.NodeInfo

	connNum int // 节点链接数
}

type MockRsp struct {
	Rsp   interface{}
	Ret   uint16
	Error error
}

// node 本节点实例
// var (
// 	node GxNode
// )

// GetRegisetr 获取master连接实例
func (node *GxNode) GetRegisetr() IRegister {
	return node.IRegister
}

// 获取节点地址，支持host1/host2:port和host:port格式
// host1-内网ip host2-外网ip
func (node *GxNode) getNodeAddress(nodeAddr string, nodeRegion uint32) string {
	arr := strings.Split(nodeAddr, ":")
	arr1 := strings.Split(arr[0], "/")

	if len(arr1) == 1 {
		return nodeAddr
	}

	i := 0
	if nodeRegion != node.Config.Region {
		i = 1
	}

	return fmt.Sprintf("%v:%v", arr1[i], arr[1])
}

// connectNode 连接到指定节点
func (node *GxNode) connectNode(info *config.NodeInfo) {
	if info.Name == node.Config.Nodename {
		return
	}

	nodeConn := &GxNodeConn{Info: info}

	address := node.getNodeAddress(info.Address, info.Region)
	for i := 0; i < node.connNum; i++ {
		conn, err := Dial("tcp", address, WithName(info.Name, CloseCallback))
		isClose := false
		if err != nil {
			isClose = true
			utils.Info(fmt.Sprintf("[%v --<xxx>-- %v]连接[%v]节点失败", node.Config.Nodename, info.Name, i), "name", info.Name, "address", address, "err", err)
		} else {
			utils.Info(fmt.Sprintf("[%v -->---<-- %v]连接[%v]节点成功", node.Config.Nodename, info.Name, i), "name", info.Name, "address", address, "region", info.Region, "isClose", isClose)
			conn.Region = info.Region
			conn.Id = i
			nodeConn.Conns = append(nodeConn.Conns, conn)
		}
	}
	nodeConn.UseTime = make([]int64, node.connNum)

	node.Mutex.Lock()
	// 保存节点的所有服务，可能多个节点都有同一个服务
	for i := 0; i < len(info.Services); i++ {
		name := info.Services[i]
		s, ok := node.ServiceNode[name]
		if ok {
			s[info.Name] = true
		} else {
			s1 := make(map[string]bool)
			s1[info.Name] = true
			node.ServiceNode[name] = s1
		}
	}
	node.Mutex.Unlock()

	node.RpcClient.Store(info.Name, nodeConn)
	if len(nodeConn.Conns) > 0 {
		node.Services.Range(func(key string, value GxService) bool {
			utils.Go(func() {
				value.NodeConn(info.Name)
			})
			return true
		})
	}
}

// Exit 退出处理
// 信号处理，程序退出统一使用kill -2
func (node *GxNode) Exit() {
	node.IRegister.DelNode(node.Config.Nodename)
	node.Services.Range(func(key string, value GxService) bool {
		utils.Go(func() {
			value.Exit()
		})
		return true
	})
}

type GxNodeParam func(*GxNode)

// WithExitTime 注册rpc回调，一般用于监控
func WithRpcCallBack(cb RpcCallBack) GxNodeParam {
	return func(node *GxNode) {
		node.RpcCallBack = cb
	}
}

// NodeConfig 节点配置
var localNode *GxNode

// InitNode 初始化服务节点
// @client 服务管理连接
// @id 节点实例id
// @nodeName 节点名
// @set 分组，空表示全局组
// @host 节点地址
// @port 节点端口
// @region 所属区域
// @cmds 注册和uint16的cmd绑定接口，用于游戏网关。原来是通过服务接口最后四个字符标识cmd，现在通过proto文件获取
func InitNode(iRegisetr IRegister, nodeConfig *NodeConfig, params ...GxNodeParam) (*GxNode, error) {
	utils.Info("初始化服务节点", "id", nodeConfig.Id, "name", nodeConfig.Nodename, "host", nodeConfig.Host, "port", nodeConfig.Port, "region", nodeConfig.Region)

	// 处理cmd
	newCmds := make(map[int32]string)
	for k, v := range nodeConfig.Cmds {
		newCmds[k] = strings.Replace(v, "_", ".", 1)
	}
	nodeConfig.Cmds = newCmds

	node := &GxNode{}
	node.IRegister = iRegisetr
	node.Config = nodeConfig
	node.connNum = 4
	node.ServiceNode = make(map[string]map[string]bool)
	node.Server = NewServer()
	node.Server.RpcCallBack = node.RpcCallBack

	// 初始化一些全局变量
	node.IDGen = utils.NewIDGen(nodeConfig.Id)
	addr := fmt.Sprintf("%v:%v", nodeConfig.Host, nodeConfig.Port)
	node.Node = &config.NodeInfo{
		Id:      nodeConfig.Id,
		Name:    nodeConfig.Nodename,
		Type:    nodeConfig.Nodetype,
		Address: addr,
		Region:  nodeConfig.Region,
		Set:     nodeConfig.Set,
	}

	for _, v := range params {
		v(node)
	}

	localNode = node
	return node, nil
}

func (node *GxNode) Start() {
	node.Services.Range(func(serviceName string, service GxService) bool {
		node.Server.RegisterName(serviceName, service)
		utils.Go(func() {
			err := service.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		})
		return true
	})

	nodeConfig := node.Config
	// 注册节点信息
	node.Services.Range(func(name string, service GxService) bool {
		node.Node.Services = append(node.Node.Services, name)
		return true
	})
	node.IRegister.RegNode(node.Node)
	log.Println("注册节点", nodeConfig.Nodename, node.Node.Address, node.Node.Services)

	utils.Go(func() {
		// 拉去公共节点
		nodes, _ := node.IRegister.QueryNodes()
		for _, v := range nodes {
			if v.Set != "" && v.Set != nodeConfig.Set {
				continue
			}
			log.Println("Get All global Node", v.Name, v.Address)
			node.connectNode(v)
		}
	})

	// 内部
	if nodeConfig.HttpPort != 0 {
		go func() {
			err := RunHttpGateway(nodeConfig.HttpPort)
			if err != nil {
				fmt.Println("listen error", err)
				os.Exit(0)
			}
		}()
	}

	listenPort := fmt.Sprintf(":%v", nodeConfig.Port)
	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		fmt.Println("listen error", err)
		return
	}

	utils.Go(func() {
		node.Server.Accept(l)
	})
}

// RangeNode 遍历节点
func (node *GxNode) RangeNode(f func(node *GxNodeConn) bool) {
	node.RpcClient.Range(func(key string, nodeInfo *GxNodeConn) bool {
		return f(nodeInfo)
	})
}

// QueryNodeStatus 查询当前节点连接状态
func (node *GxNode) QueryNodeStatus() []*GxNodeConn {
	var cs []*GxNodeConn
	node.RpcClient.Range(func(key string, nodeInfo *GxNodeConn) bool {
		cs = append(cs, &GxNodeConn{
			Info: &config.NodeInfo{
				Id:      nodeInfo.Info.Id,
				Name:    nodeInfo.Info.Name,
				Address: nodeInfo.Info.Address,
			},
			Close: nodeInfo.Close,
		})
		return true
	})

	return cs
}

// ConnectNewNode 尝试连接新节点
func (node *GxNode) ConnectNewNode(name string) {
	info, _ := node.IRegister.QueryNode(name)
	if info != nil && (info.Set == "" || info.Set == node.Config.Set) {
		node.connectNode(info)
	}
}

// DisconnectNode 断开节点连接
func (node *GxNode) DisconnectNode(name string) {
	nodeInfo, ok := node.RpcClient.Load(name)
	if ok {
		for _, v := range nodeInfo.Conns {
			if v != nil {
				v.Close()
			}
		}
		nodeInfo.Close = true
		node.RpcClient.Delete(name)

		utils.Info("删除节点", "nodeName", name)
	}
}

// FindRpcConnByService 多节点模式下，返回提供服务的所有节点，自己处理，例如
// key := fmt.Sprintf("%v:%v", appid, username)
// verifies := rpc.FindRpcConnByService(serviceName)
//
//	if len(verifies) == 0 {
//		return static.RetServiceStop, nil
//	}
//
// ring := ketama.NewRing(200)
//
//	for k, _ := range verifies {
//		ring.AddNode(k, 100)
//		ring.Bake()
//	}
//
// name := ring.Hash(key)
// ret, err := verifies[name].Call(funcName, &req, &rsp)
// return uint16(ret), err
func FindRpcConnByService(serviceName string) map[string]*Client {
	localNode.Mutex.Lock()
	defer localNode.Mutex.Unlock()

	NodeNames, ok := localNode.ServiceNode[serviceName]
	if !ok || len(NodeNames) == 0 {
		return nil
	}

	m := make(map[string]*Client)
	for k := range NodeNames {
		conn := getNode(k)
		if conn == nil {
			continue
		}
		m[k] = conn
	}

	return m
}

func getNode(nodeName string) *Client {
	nc, ok2 := localNode.RpcClient.Load(nodeName)
	if !ok2 {
		return nil
	}

	// TODO: maybe nil connect or closed connect
	// var cli *Client
	// for _, v := range nc.Conns {
	// 	if cli == nil || cli.UseTime < v.UseTime {
	// 		cli = v
	// 	}
	// }

	return nc.getConn()
}

// GetNode 获取指定节点的rpc连接实例
func (node *GxNode) GetNode(nodeName string) *Client {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	return getNode(nodeName)
}

// RegisterService 注册服务
func (node *GxNode) RegisterService(service GxService) {
	node.Services.Store(service.GetName(), service)
}

// NodeCall 节点rpc调用
func NodeCall(nodeName string, serviceMethod string, req proto.Message, rsp proto.Message) (uint16, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		buff, _ := proto.Marshal(mockRsp.(proto.Message))
		proto.Unmarshal(buff, rsp)
		return mockRet, mockErr
	}

	if nodeName == localNode.Config.Nodename {
		return localNode.Server.InternalCall(EmptyContext(), serviceMethod, req, rsp)
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		return node.Call(EmptyContext(), serviceMethod, req, rsp)
	}

	return 1, fmt.Errorf("node %v not exists", nodeName)
}

// NodeJsonCallWithConn 节点rpc调用
func NodeJsonCallWithConn(context *Context, nodeName string, serviceMethod string, reqBuff []byte) (uint16, []byte, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		return mockRet, mockRsp.([]byte), mockErr
	}
	if nodeName == localNode.Config.Nodename {
		return localNode.Server.RawCall(context, serviceMethod, reqBuff, true)
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		return node.JsonCall(context, serviceMethod, reqBuff)
	}

	return 1, nil, fmt.Errorf("node %v not exists", nodeName)
}

// NodeRawCallWithConn 节点rpc调用
func NodeRawCallWithConn(context *Context, nodeName string, serviceMethod string, reqBuff []byte) (uint16, []byte, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		return mockRet, mockRsp.([]byte), mockErr
	}
	if nodeName == localNode.Config.Nodename {
		return localNode.Server.RawCall(context, serviceMethod, reqBuff, false)
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		return node.RawCall(context, serviceMethod, reqBuff)
	}

	return 1, nil, fmt.Errorf("node %v not exists", nodeName)
}

// NodeSend 向指定节点异步发送消息
func NodeSend(nodeName string, serviceMethod string, req proto.Message) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	if nodeName == localNode.Config.Nodename {
		utils.Go(func() { localNode.Server.InternalCall(EmptyContext(), serviceMethod, req, nil) })
		return nil
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		node.Send(EmptyContext(), serviceMethod, req)
		return nil
	}

	return fmt.Errorf("node %v not exists", nodeName)
}

// NodeCallWithConn 调用玩家所属网关接口
func NodeCallWithConn(context *Context, nodeName string, serviceMethod string, req proto.Message, rsp proto.Message) (uint16, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		buff, _ := proto.Marshal(mockRsp.(proto.Message))
		proto.Unmarshal(buff, rsp)
		return mockRet, mockErr
	}

	if nodeName == localNode.Config.Nodename {
		return localNode.Server.InternalCall(context, serviceMethod, req, rsp)
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		return node.Call(context, serviceMethod, req, rsp)
	}

	return 1, fmt.Errorf("gate node %v not exists", context.GateName)
}

// NodeSendWithConn ...
func NodeSendWithConn(context *Context, nodeName string, serviceMethod string, req proto.Message) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	if nodeName == localNode.Config.Nodename {
		utils.Go(func() {
			localNode.Server.InternalCall(context, serviceMethod, req, nil)
		})
		return nil
	}

	node := localNode.GetNode(nodeName)
	if node != nil {
		node.Send(context, serviceMethod, req)
		return nil
	}

	return fmt.Errorf("gate node %v not exists", context.GateName)
}

// Call 服务之间的rpc调用
func Call(context *Context, serviceMethod string, req proto.Message, rsp proto.Message) (uint16, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		buff, _ := proto.Marshal(mockRsp.(proto.Message))
		proto.Unmarshal(buff, rsp)
		return mockRet, mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return 1, fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			return client.Call(context, serviceMethod, req, rsp)
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		return localNode.Server.InternalCall(context, serviceMethod, req, rsp)
	} else {
		client := getClient(serviceName)
		if client == nil {
			return 1, fmt.Errorf("Call not support node rpc")
		}

		return client.Call(context, serviceMethod, req, rsp)
	}
}

// Send 服务之间的异步调用
func Send(context *Context, serviceMethod string, req proto.Message) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			client.Send(context, serviceMethod, req)
			return nil
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		utils.Go(func() {
			localNode.Server.InternalCall(context, serviceMethod, req, nil)
		})
		return nil
	} else {
		client := getClient(serviceName)
		if client == nil {
			return fmt.Errorf("not support node rpc")
		}

		client.Send(context, serviceMethod, req)
		return nil
	}
}

// Broadcast 服务广播, 消息会发送到所有注册了该服务的节点
func Broadcast(serviceMethod string, req proto.Message) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		utils.Go(func() {
			localNode.Server.InternalCall(EmptyContext(), serviceMethod, req, nil)
		})
	}

	clients := FindRpcConnByService(serviceName)
	if len(clients) == 0 {
		return nil
	}

	for name, client := range clients {
		if name == localNode.Config.Nodename {
			continue
		}

		client.Send(EmptyContext(), serviceMethod, req)
	}

	return nil
}

// BroadcastCall 顺序调用
func BroadcastCall(serviceMethod string, req proto.Message, rsp proto.Message, f func(nodeName string) bool) (uint16, error) {
	if ok, _, ret, mockErr := callMock(serviceMethod); ok {
		return ret, mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		localNode.Server.InternalCall(EmptyContext(), serviceMethod, req, rsp)
		if !f(localNode.Config.Nodename) {
			return 0, nil
		}
	}

	clients := FindRpcConnByService(serviceName)
	if len(clients) == 0 {
		return 0, nil
	}

	for name, client := range clients {
		if name == localNode.Config.Nodename {
			continue
		}

		client.Call(EmptyContext(), serviceMethod, req, rsp)
		if !f(name) {
			return 0, nil
		}
	}

	return 0, nil
}

// JsonCall ...
func JsonCall(context *Context, serviceMethod string, reqBuff []byte) (uint16, []byte, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		return mockRet, mockRsp.([]byte), mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return 1, nil, fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			return client.JsonCall(context, serviceMethod, reqBuff)
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		return localNode.Server.RawCall(context, serviceMethod, reqBuff, true)
	} else {
		client := getClient(serviceName)
		if client == nil {
			return 1, nil, fmt.Errorf("not support node rpc")
		}
		return client.JsonCall(context, serviceMethod, reqBuff)
	}
}

// JsonSend ...
func JsonSend(context *Context, serviceMethod string, reqBuff []byte) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			client.JsonSend(context, serviceMethod, reqBuff)
			return nil
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		utils.Go(func() {
			localNode.Server.RawCall(context, serviceMethod, reqBuff, true)
		})
		return nil
	} else {
		client := getClient(serviceName)
		if client == nil {
			return fmt.Errorf("not support node rpc")
		}

		client.JsonSend(context, serviceMethod, reqBuff)
		return nil
	}
}

// RawCall ...
func RawCall(context *Context, serviceMethod string, reqBuff []byte) (uint16, []byte, error) {
	if ok, mockRsp, mockRet, mockErr := callMock(serviceMethod); ok {
		return mockRet, mockRsp.([]byte), mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return 1, nil, fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			return client.RawCall(context, serviceMethod, reqBuff)
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		return localNode.Server.RawCall(context, serviceMethod, reqBuff, false)
	} else {
		client := getClient(serviceName)
		if client == nil {
			return 1, nil, fmt.Errorf("not support node rpc")
		}

		return client.RawCall(context, serviceMethod, reqBuff)
	}
}

// RawSend ...
func RawSend(context *Context, serviceMethod string, reqBuff []byte) error {
	if ok, _, _, mockErr := callMock(serviceMethod); ok {
		return mockErr
	}

	serviceName, _ := splitServiceMethod(serviceMethod)

	// 根据路由转发
	for i := 0; i < len(context.Nodes); i++ {
		if serviceName == context.Nodes[i].ServiceName {
			client := getNode(context.Nodes[i].NodeName)
			if client == nil {
				return fmt.Errorf("node[%v] not exist", context.Nodes[i].NodeName)
			}
			client.RawSend(context, serviceMethod, reqBuff)
			return nil
		}
	}

	_, ok := localNode.Services.Load(serviceName)
	if ok {
		// 内部调用
		utils.Go(func() {
			localNode.Server.RawCall(context, serviceMethod, reqBuff, false)
		})
		return nil
	} else {
		client := getClient(serviceName)
		if client == nil {
			return fmt.Errorf("not support node rpc")
		}

		client.RawSend(context, serviceMethod, reqBuff)
		return nil
	}
}

var emptyContext = &Context{
	Remote: "empty",
}

// EmptyContext 返回一个空连接
func EmptyContext() *Context {
	return emptyContext
}

// CloseCallback rpc连接断开回调
func CloseCallback(name string, err error) {
	utils.Info(fmt.Sprintf("[%v --<xxx>-- %v]节点断开连接", localNode.Node.Name, name), "name", name, "err", err)
	localNode.RpcClient.Delete(name)

	localNode.Services.Range(func(key string, value GxService) bool {
		utils.Go(func() {
			value.NodeClose(name)
		})
		return true
	})
}

// NewId 生成一个新id
func NewId(moduleId uint64) uint64 {
	return localNode.IDGen.NewID(moduleId)
}

// getClient 寻找和本节点匹配的节点
func getClient(serviceName string) *Client {
	clients := FindRpcConnByService(serviceName)
	if len(clients) == 0 {
		return nil
	}

	// 优先找区域匹配节点，如果找不到就随便找一个
	var client, client2 *Client
	for _, c := range clients {
		if client2 == nil {
			client2 = c
		}
		if client == nil && localNode.Config.Region == c.Region {
			client = c
		}
	}
	if client == nil {
		client = client2
	}
	return client
}

func InitMock() {
	localNode.MockData = make(map[string]*MockRsp)
}

func InsertMock(serviceMethon string, rsp *MockRsp) {
	localNode.MockData[serviceMethon] = rsp
}

func callMock(serviceMethon string) (bool, interface{}, uint16, error) {
	if localNode.MockData == nil {
		return false, nil, 0, nil
	}
	data, ok := localNode.MockData[serviceMethon]
	if !ok {
		return ok, nil, 0, nil
	}

	return true, data.Rsp, data.Ret, data.Error
}

func RemoveMock() {
	localNode.MockData = nil
}

func SubmitEvent(serviceName, eventName string, args ...interface{}) {
	if serviceName != "" {
		v, ok := localNode.Services.Load(serviceName)
		if ok {
			utils.Go(func() {
				v.OnEvent(eventName, args...)
			})
			return
		}
	}
	localNode.Services.Range(func(key string, value GxService) bool {
		utils.Go(func() {
			value.OnEvent(eventName, args...)
		})
		return true
	})
}

// FindServiceMethod 查询服务方法
func FindServiceMethod(cmd uint16) string {
	i, ok := localNode.Config.Cmds[int32(cmd)]
	if !ok {
		return ""
	}

	return i
}

func GetNodeConfig() *NodeConfig {
	return localNode.Config
}
