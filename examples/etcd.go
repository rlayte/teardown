package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/rlayte/teardown"
)

func Test(cluster Cluster) {
	cluster.Setup()

	// Tests
	tests := teardown.NewEtcdTests(cluster.addresses)
	tests.Step()

	cluster.Teardown()
}

func main() {
	cluster := etcd.Cluster()
	Test(cluster)
}
