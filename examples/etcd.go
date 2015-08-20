package main

import (
	"github.com/rlayte/teardown"
	"github.com/rlayte/teardown/adapters/etcd"
)

func Test(cluster teardown.Cluster) {
	cluster.Setup()

	// Tests
	tests := teardown.NewEtcdTests(cluster.Addresses())
	tests.Step()

	cluster.Teardown()
}

func main() {
	cluster := etcd.Cluster()
	Test(cluster)
}
