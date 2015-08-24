package teardown

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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
	Status   RequestStatus
	Response ResponseStatus
	Key      string
	Value    string
}

type TestRunner struct {
	requests []Request
	client   Client
	cluster  Cluster
	count    int
}

func (t *TestRunner) Step() error {
	request := Request{
		Key:   fmt.Sprintf("/%d", t.count),
		Value: "hi!",
	}

	err := t.client.Put(request.Key, request.Value)

	if err != nil {
		request.Status = Fail
	} else {
		request.Status = Ack
	}

	value, err := t.client.Get(request.Key)

	if err != nil {
		request.Response = No
	} else {
		request.Response = Yes
	}

	log.Println("Error", err)
	log.Println("Get response", value)

	t.requests = append(t.requests, request)
	t.count++

	return nil
}

func (t *TestRunner) Report() {
	successfulWrites := 0
	correctFailures := 0
	missingWrites := 0
	extraWrites := 0

	for _, request := range t.requests {
		if request.Status == Ack && request.Response == Yes {
			successfulWrites++
		}
		if request.Status == Fail && request.Response == Yes {
			extraWrites++
		}
		if request.Status == Fail && request.Response == No {
			correctFailures++
		}
		if request.Status == Ack && request.Response == No {
			missingWrites++
		}
	}

	log.Println("Successful writes:", successfulWrites)
	log.Println("Correct failures:", correctFailures)
	log.Println("Missing writes:", missingWrites)
	log.Println("Extra writes", extraWrites)
}

func (t *TestRunner) Run() {
	t.cluster.Setup()
	bio := bufio.NewReader(os.Stdin)
	for i := 0; i < 100; i++ {
		bio.ReadLine()
		log.Println("Next")
		t.Step()
	}

	t.Report()
	t.cluster.Teardown()
}

func NewTestRunner(cluster Cluster, client Client) *TestRunner {
	t := TestRunner{
		requests: []Request{},
		client:   client,
		cluster:  cluster,
	}

	return &t
}
