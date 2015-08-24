package teardown

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"

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
	log.Println("Partitioning network:", half)
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
	nodes := []string{}

	for _, address := range cluster.Addresses() {
		u, err := url.Parse(address)

		if err != nil {
			panic(err)
		}

		host, _, _ := net.SplitHostPort(u.Host)
		nodes = append(nodes, host)
	}

	return &LocalNemesis{nodes, iptables.UnixIpTables{}}
}
