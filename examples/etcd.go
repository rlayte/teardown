package main

import (
	"fmt"
	"log"

	"github.com/rlayte/teardown"
	"github.com/rlayte/teardown/adapters/etcd"
)

type EtcdTests struct {
	addresses []string
	requests  []teardown.Request
	client    *etcd.Client
	count     int
}

func NewEtcdTests(addresses []string) *EtcdTests {
	t := EtcdTests{}

	t.addresses = addresses
	t.client = etcd.NewClient(addresses)
	t.requests = []Request{}

	return &t
}

func (t *EtcdTests) Step() error {
	request := Request{
		key:   fmt.Sprintf("/%d", t.count),
		value: "hi!",
	}

	response, err := t.client.Set(request.key, request.value, 0)

	if err != nil {
		log.Fatal("Set", err)
	}

	t.requests = append(t.requests, request)
	t.count++

	getResp, err := t.client.Get(request.key, false, false)

	log.Println("Response", response)
	log.Println("Error", err)
	log.Println("Get response", getResp)
}

func (t *EtcdTests) Finalize() {
}

func main() {
	cluster := etcd.Cluster()
	tests := NewEtcdTests(cluster.Addresses())
	teardown.RunTests(cluster, tests)
}
