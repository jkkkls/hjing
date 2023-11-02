package rpc

import (
	"os"

	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/etcdapi"
	"github.com/jkkkls/hjing/utils"
	"gopkg.in/yaml.v2"
)

// 服务发现
type IRegister interface {
	DelNode(string) error
	RegNode(*config.NodeInfo) error
	WatchNode(func(string, *config.NodeInfo)) error
	QueryNodes() (map[string]*config.NodeInfo, error)
	QueryNode(string) (*config.NodeInfo, error)
}

// EtcdRegister 基于etcd的服务发现
type EtcdRegister struct {
	client *etcdapi.EtcdCli
}

func (r *EtcdRegister) DelNode(name string) error {
	etcdapi.DelRegisterNode(r.client, name)
	return nil
}

func (r *EtcdRegister) RegNode(info *config.NodeInfo) error {
	etcdapi.RegisterNode(r.client, info)
	return nil
}

func (r *EtcdRegister) WatchNode(f func(string, *config.NodeInfo)) error {
	r.client.Watch(config.NodeRegisterKey, f)
	return nil
}

func (r *EtcdRegister) QueryNodes() (map[string]*config.NodeInfo, error) {
	return etcdapi.GetAllRegisterNode(r.client), nil
}

func (r *EtcdRegister) QueryNode(name string) (*config.NodeInfo, error) {
	return etcdapi.GetRegisterNode(r.client, name), nil
}

func NewEtcdRegister(addrs ...string) (IRegister, error) {
	client, err := etcdapi.ConnEtcd(addrs...)
	if err != nil {
		return nil, err
	}

	return &EtcdRegister{client: client}, nil
}

// YamlRegister 基于yaml的服务发现
type YamlRegister struct {
	configName string
	nodes      map[string]*config.NodeInfo
}

func (r *YamlRegister) DelNode(name string) error {
	return nil
}

func (r *YamlRegister) RegNode(info *config.NodeInfo) error {
	return nil
}

func (r *YamlRegister) WatchNode(f func(string, *config.NodeInfo)) error {
	utils.InitWatch(func() {
		var add, del []string
		newNodes, err := loadYaml(r.configName)
		if err != nil {
			utils.Warn("load yaml config error", "err", err)
			return
		}
		for k := range newNodes {
			if _, ok := r.nodes[k]; !ok {
				add = append(add, k)
			}
		}
		for k := range r.nodes {
			if _, ok := newNodes[k]; !ok {
				del = append(del, k)
			}
		}

		for _, v := range add {
			f(v, newNodes[v])
		}
		for _, v := range del {
			f(v, nil)
		}
	}, r.configName)
	return nil
}

func (r *YamlRegister) QueryNodes() (map[string]*config.NodeInfo, error) {
	return r.nodes, nil
}

func (r *YamlRegister) QueryNode(name string) (*config.NodeInfo, error) {
	return r.nodes[name], nil
}

func loadYaml(configName string) (map[string]*config.NodeInfo, error) {
	buff, err := os.ReadFile(configName)
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]*config.NodeInfo)
	err = yaml.Unmarshal(buff, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func NewYamlRegister(configName string) (IRegister, error) {
	nodes, err := loadYaml(configName)
	if err != nil {
		return nil, err
	}
	r := &YamlRegister{configName: configName, nodes: nodes}
	return r, nil
}
