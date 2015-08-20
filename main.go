package main

type EtcdCluster struct {
	addresses []string
}

func (c *EtcdCluster) Setup() error {
}

func (c *EtcdCluster) Teardown() error {
}

func main() {
	cluster := EtcdCluster{}
	cluster.Setup()

	// Tests

	cluster.Teardown()
}
