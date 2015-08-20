package teardown

import "log"

type Tests interface {
	Step() error
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
	Status   RequestStatus
	Response ResponseStatus
	Key      string
	Value    string
}

func RunTests(cluster Cluster, tests Tests) {
	log.Println("Running tests")
	cluster.Setup()

	tests.Step()
	tests.Finalize()

	cluster.Teardown()
}
