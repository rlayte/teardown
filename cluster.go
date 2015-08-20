package teardown

type Cluster interface {
	Setup()
	Teardown()
	Addresses() []string
}
