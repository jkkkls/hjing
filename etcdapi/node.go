package etcdapi

import (
	"encoding/json"
	"fmt"

	"github.com/jkkkls/hjing/config"
)

// Copyright 2017 guangbo. All rights reserved.

//节点配置，节点启动读取
//key: node:nodeName
//节点通用配置，用于服务注册服务发现
//key: nodeRegister

func GetAllRegisterNode(client *EtcdCli) map[string]*config.NodeInfo {
	m := make(map[string]*config.NodeInfo)
	values, err := client.KeyPrefix(config.NodeRegisterKey + ":")
	if err != nil {
		return m
	}
	for i := 0; i < len(values); i = i + 2 {
		info := &config.NodeInfo{}
		json.Unmarshal([]byte(values[i+1]), info)
		m[info.Name] = info
	}
	return m
}

// RegisterNode 注册节点和服务
func RegisterNode(client *EtcdCli, info *config.NodeInfo) {
	key := fmt.Sprintf("%v:%v", config.NodeRegisterKey, info.Name)
	buff, _ := json.Marshal(info)
	client.Put(key, string(buff))
}

func DelRegisterNode(client *EtcdCli, name string) {
	key := fmt.Sprintf("%v:%v", config.NodeRegisterKey, name)
	client.Delete(key)
}

func GetRegisterNode(client *EtcdCli, nodeName string) *config.NodeInfo {
	key := fmt.Sprintf("%v:%v", config.NodeRegisterKey, nodeName)
	value, err := client.Get(key)
	if value == "" || err != nil {
		return nil
	}
	info := &config.NodeInfo{}
	json.Unmarshal([]byte(value), info)
	return info
}
