package teardown

import (
	"fmt"
	"log"

	"github.com/coreos/go-etcd/etcd"
)

type Tests interface {
	Setup(addresses []string)
	Step()
	Finalize()
}

type RequestStatus int
type ResponseStatus int

const (
	Ack RequestStatus = iota
	Unknown
	Fail

	Yes ResponseStatus = iota
	No
	Maybe
)

type Request struct {
	status   RequestStatus
	response ResponseStatus
	key      string
	value    string
}

type EtcdTests struct {
	addresses []string
	requests  []Request
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

func (t *EtcdTests) Step() {
	request := Request{
		key:   fmt.Sprintf("%d", t.count),
		value: "hi!",
	}

	response, err := t.client.Set(request.key, request.value, 0)

	t.requests = append(t.requests, request)
	t.count++

	log.Println("Response", response, err)
}
