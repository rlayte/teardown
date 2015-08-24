// etcd.go

package main

import (
	"github.com/rlayte/teardown"
	"github.com/rlayte/teardown/examples/client"
	"github.com/rlayte/teardown/examples/server"
)

func main() {
	cluster := server.New()
	client := client.New(cluster)
	tests := teardown.NewTestRunner(cluster, client)
	tests.Run()
}
