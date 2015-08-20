package teardown

type Cluster interface {
	Setup() error
	Teardown() error
	Addresses() []string
}
