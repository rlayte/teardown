package teardown

import (
	"math/rand"

	"github.com/rlayte/teardown/iptables"
)

type Nemesis struct {
}

func PartitionHalf(nodes []string) {
	half := len(nodes) / 2
	iptables.Partition(nodes, half)
}

func PartitionRandom(nodes []string) {
	position := rand.Intn(len(nodes))
	iptables.Partition(nodes, position)
}
