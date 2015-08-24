// client.go

package client

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/rlayte/teardown"
)

type EtcdClient struct {
	client *etcd.Client
}

func (c EtcdClient) Get(key string) (string, error) {
	_, err := c.client.Get(key, false, false)
	return "", err
}

func (c EtcdClient) Put(key string, value string) error {
	_, err := c.client.Set(key, value, 0)
	return err
}

func New(cluster teardown.Cluster) teardown.Client {
	return EtcdClient{etcd.NewClient(cluster.Addresses())}
}
