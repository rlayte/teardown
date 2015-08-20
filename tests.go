package teardown

type Tests interface {
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

func RunTests(cluster Cluster, tests Tests) {
	cluster.Setup()

	tests.Step()

	cluster.Teardown()
}
