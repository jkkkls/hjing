package rpc

import (
	"github.com/jkkkls/hjing/config"
	"github.com/jkkkls/hjing/etcdapi"
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
