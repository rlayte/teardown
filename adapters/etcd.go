// adapters/etcd.go

package etcd

import "os/exec"

type EtcdAdapter struct {
	internal_addresses []string
	Addresses          []string // for clients
}

func (c *EtcdAdapter) Setup() error {
	for _, address := range c.addresses {
		exec.Cmd("etcd", "--initail-cluster", cluster_string)
	}
}

func (c *EtcdAdapter) Teardown() error {
}

func Cluster() *EtcdAdapter {
	return EtcdAdapter{[]string{
		"127.0.0.2",
		"127.0.0.3",
		"127.0.0.4",
		"127.0.0.5",
		"127.0.0.6",
	}}
}
