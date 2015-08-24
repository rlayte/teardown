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
	nodes    []string
	iptables iptables.IpTables
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
	n.iptables.PartitionLevel(n.nodes, half)
}

func (n *LocalNemesis) PartitionRandom() {
	position := rand.Intn(len(n.nodes))
	n.PartitionSingle(position)
}

func (n *LocalNemesis) PartitionSingle(position int) {
	n1 := n.positionToAddress(position)
	for _, n2 := range n.nodes {
		if n1 != n2 {
			n.iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) Bridge() {
	position := len(n.nodes) / 2
	for _, n1 := range n.nodes[:position] {
		for _, n2 := range n.nodes[position+1:] {
			n.iptables.Deny(n1, n2)
		}
	}
}

func (n *LocalNemesis) Heal() {
	n.iptables.Heal()
}

func NewNemesis(cluster Cluster) *LocalNemesis {
	return &LocalNemesis{cluster.Addresses(), iptables.UnixIpTables{}}
}
