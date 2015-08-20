package main

import (
	"fmt"
	"log"

	etcdclient "github.com/coreos/go-etcd/etcd"
	"github.com/rlayte/teardown"
	"github.com/rlayte/teardown/adapters/etcd"
)

type EtcdTests struct {
	addresses []string
	requests  []teardown.Request
	client    *etcdclient.Client
	count     int
}

func NewEtcdTests(addresses []string) *EtcdTests {
	t := EtcdTests{}

	t.addresses = addresses
	t.client = etcdclient.NewClient(addresses)
	t.requests = []teardown.Request{}

	return &t
}

func (t *EtcdTests) Step() error {
	request := teardown.Request{
		Key:   fmt.Sprintf("/%d", t.count),
		Value: "hi!",
	}

	setResp, err := t.client.Set(request.Key, request.Value, 0)

	if err != nil {
		request.Status = teardown.Fail
	} else {
		request.Status = teardown.Ack
	}

	getResp, err := t.client.Get(request.Key, false, false)

	if err != nil {
		request.Response = teardown.No
	} else {
		request.Response = teardown.Yes
	}

	log.Println("Response", setResp)
	log.Println("Error", err)
	log.Println("Get response", getResp)

	t.requests = append(t.requests, request)
	t.count++

	return nil
}

func (t *EtcdTests) Finalize() {
	successfulWrites := 0
	correctFailures := 0
	missingWrites := 0
	extraWrites := 0

	for _, request := range t.requests {
		if request.Status == teardown.Ack && request.Response == teardown.Yes {
			successfulWrites++
		}
		if request.Status == teardown.Fail && request.Response == teardown.Yes {
			extraWrites++
		}
		if request.Status == teardown.Fail && request.Response == teardown.No {
			correctFailures++
		}
		if request.Status == teardown.Ack && request.Response == teardown.No {
			missingWrites++
		}
	}

	log.Println("Successful writes:", successfulWrites)
	log.Println("Correct failures:", correctFailures)
	log.Println("Missing writes:", missingWrites)
	log.Println("Extra writes", extraWrites)
}

func main() {
	cluster := etcd.Cluster()
	tests := NewEtcdTests(cluster.Addresses())
	teardown.RunTests(cluster, tests)
}
