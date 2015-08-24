// etcd_test.go

package etcd

import (
	"log"
	"testing"

	"github.com/coreos/go-etcd/etcd"
	"github.com/rlayte/teardown"
)

var cluster teardown.Cluster
var client *etcd.Client
var nemesis teardown.Nemesis

func init() {
	log.Println("Setting up etcd cluster")
	cluster = NewEtcdCluster()
	client = etcd.NewClient(cluster.Addresses())
	nemesis = teardown.NewNemesis(cluster)
}

func TestPartition(t *testing.T) {
	client.Set("foo", "bar", 0)
	value, err := client.Get("foo", false, false)

	if err != nil {
		t.Error(err)
	}

	if value.Action != "bar" {
		t.Error("Foo should equal bar")
	}
}
