package etcdapi

import (
	"encoding/json"
	"fmt"
)

// Copyright 2017 guangbo. All rights reserved.

//节点配置，节点启动读取
//key: node:nodeName
//节点通用配置，用于服务注册服务发现
//key: nodeRegister

const (
	NodeRegisterKey = "nodeRegister:"
)

type FunctionInfo struct {
	Name string `json:"name,omitempty"`
	Cmd  uint16 `json:"cmd,omitempty"`
}

type ServiceInfo struct {
	Name string          `json:"name,omitempty"`
	Func []*FunctionInfo `json:"func,omitempty"`
}

// 节点注册函数
type NodeInfo struct {
	Id      uint64         `json:"id,omitempty"`
	Name    string         `json:"name,omitempty"`    //
	Type    string         `json:"type,omitempty"`    //
	Address string         `json:"address,omitempty"` //ip:port
	Service []*ServiceInfo `json:"service,omitempty"` //服务列表
	Region  uint32         `json:"region,omitempty"`
	Set     string         `json:"set,omitempty"`
}

func GetAllRegisterNode(client *EtcdCli) map[string]*NodeInfo {
	m := make(map[string]*NodeInfo)
	values, err := client.KeyPrefix(NodeRegisterKey)
	if err != nil {
		return m
	}
	for i := 0; i < len(values); i = i + 2 {
		info := &NodeInfo{}
		json.Unmarshal([]byte(values[i+1]), info)
		m[info.Name] = info
	}
	return m
}

// RegisterNode 注册节点和服务
func RegisterNode(client *EtcdCli, info *NodeInfo) {
	key := fmt.Sprintf("%v:%v", NodeRegisterKey, info.Name)
	buff, _ := json.Marshal(info)
	client.Put(key, string(buff))
}

func DelRegisterNode(client *EtcdCli, info *NodeInfo) {
	key := fmt.Sprintf("%v:%v", NodeRegisterKey, info.Name)
	client.Delete(key)
}

func GetRegisterNode(client *EtcdCli, nodeName, set string) *NodeInfo {
	key := fmt.Sprintf("%v:%v", NodeRegisterKey, nodeName)
	value, err := client.Get(key)
	if value == "" || err != nil {
		return nil
	}
	info := &NodeInfo{}
	json.Unmarshal([]byte(value), info)
	return info
}
