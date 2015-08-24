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
	Bridge()
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
	n.PartitionSingle(position)
}

// What does this do?
func (n *LocalNemesis) PartitionSingle(position int) {
	n1 := n.positionToAddress(position)
	for _, n2 := range n.nodes {
		if n1 != n2 {
			iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) Bridge() {
	bridge_index := len(n.nodes) / 2 // Fails with empty list.
	for _, n1 := range n.nodes[:bridge_index] {
		for _, n2 := range n.nodes[bridge_index+1:] {
			iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) Heal() {
	iptables.Heal()
}

func NewNemesis(cluster Cluster) Nemesis {
	return &LocalNemesis{cluster.Addresses()}
}
