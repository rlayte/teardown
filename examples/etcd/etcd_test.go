// etcd_test.go

package etcd

import (
	"fmt"
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

type KeyValue struct {
	key   string
	value string
}

func TestPartition(t *testing.T) {
	requests := []KeyValue{}

	nemesis.PartitionHalf()

	for i := 0; i < 2000; i++ {
		key := fmt.Sprintf("etcd-key-%d", i)
		value := fmt.Sprintf("etcd-value-%d", i)
		requests = append(requests, KeyValue{key, value})
		log.Println("Putting", key, value)
		client.Set(key, value, 0)
	}

	for _, request := range requests {
		resp, err := client.Get(request.key, false, false)

		log.Println("Getting", request.key)

		if err != nil {
			t.Fatal(err)
		}

		if resp.Node.Value != request.value {
			t.Errorf("%s should equal %s", resp.Node.Value, request.value)
		}
	}

	nemesis.Heal()
}
