package main

import "github.com/rlayte/jepsen-go"

type EtcdCluster struct {
	addresses []string
}

func (c *EtcdCluster) Setup() error {
	return nil
}

func (c *EtcdCluster) Teardown() error {
	return nil
}

func main() {
	cluster := EtcdCluster{}
	cluster.Setup()

	// Tests
	tests := teardown.NewEtcdTests(cluster.addresses)
	tests.Step()

	cluster.Teardown()
}
