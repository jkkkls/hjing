package etcdapi

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdCli struct {
	Cli *clientv3.Client
}

func ConnEtcd(addrs ...string) (*EtcdCli, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdCli{
		Cli: cli,
	}, nil
}

func (ec *EtcdCli) Start() {

}

func (ec *EtcdCli) Get(key string) (string, error) {
	rsp, err := ec.Cli.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	if len(rsp.Kvs) == 0 {
		return "", nil
	}

	return string(rsp.Kvs[0].Value), nil
}

func (ec *EtcdCli) KeyPrefix(prefix string) ([]string, error) {
	rsp, err := ec.Cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range rsp.Kvs {
		ret = append(ret, string(v.Key), string(v.Value))
	}
	return ret, nil
}

func (ec *EtcdCli) Put(key, value string) error {
	_, err := ec.Cli.Put(context.TODO(), key, value)
	return err
}

func (ec *EtcdCli) Delete(key string) error {
	_, err := ec.Cli.Delete(context.TODO(), key)
	return err
}
