// etcd.go

package main

import (
	"github.com/rlayte/teardown"
	"github.com/rlayte/teardown/examples/etcd/client"
	"github.com/rlayte/teardown/examples/etcd/server"
)

func main() {
	cluster := server.New()
	client := client.New(cluster)
	tests := teardown.NewTestRunner(cluster, client)
	tests.Run()
}
