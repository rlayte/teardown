package teardown

import (
	"fmt"
	"math/rand"

	"github.com/rlayte/teardown/iptables"
)

type Nemesis interface {
	PartitionHalf()
	PartitionRandom()
	PartitionSingle(position int)
	PartitionLeader()
	Bridge()
	Fail(position int)
	FailRandom()
	FailLeader()
	Heal()
}

type LocalNemesis struct {
	nodes []string
}

func (n *LocalNemesis) positionToAddress(position int) string {
	if position >= len(n.nodes) {
		panic(fmt.Sprintf(
			"Position %d too large for set of %d",
			position,
			len(n.nodes),
		))
	}
	return n.nodes[position]
}

func (n *LocalNemesis) PartitionHalf() {
	half := len(n.nodes) / 2
	iptables.PartitionLevel(n.nodes, half)
}

func (n *LocalNemesis) PartitionRandom() {
	position := rand.Intn(len(n.nodes))
	iptables.PartitionLevel(n.nodes, position)
}

// What does this do?
func (n *LocalNemesis) PartitionSingle(n1 int) {
}

func (n *LocalNemesis) PartitionLeader() {
}

func (n *LocalNemesis) Bridge() {
	bridge_index := len(n.nodes) / 2 // Fails with empty list.
	for _, n1 := range n.nodes[:bridge_index] {
		for _, n2 := range n.nodes[bridge_index+1:] {
			iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) Fail(position int) {
	n1 := n.positionToAddress(position)
	for _, n2 := range n.nodes {
		if n1 != n2 {
			iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) FailRandom() {
	position := rand.Intn(len(n.nodes))
	n.Fail(position)
}

func (n *LocalNemesis) FailLeader() {
}

func (n *LocalNemesis) Heal() {
	iptables.Heal()
}
