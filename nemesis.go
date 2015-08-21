package teardown

import (
	"math/rand"

	"github.com/rlayte/teardown/iptables"
)

type Nemesis interface {
	PartitionHalf()
	PartitionRandom()
	PartitionSingle(node int)
	PartitionLeader()
	Bridge()
	Fail(node int)
	FailRandom()
	FailLeader()
	Heal()
}

type LocalNemesis struct {
	nodes []string
}

func (n *LocalNemesis) PartitionHalf() {
	half := len(n.nodes) / 2
	iptables.Partition(n.nodes, half)
}

func (n *LocalNemesis) PartitionRandom() {
	position := rand.Intn(len(n.nodes))
	iptables.Partition(n.nodes, position)
}
