package iptables

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func iptables(args string) []byte {
	out, err := exec.Command("iptables", strings.Split(args, " ")...).Output()

	if err != nil {
		log.Fatal("Iptables error: ", err, args)
	}

	return out
}

type IpTables interface {
	PartitionLevel(nodes []string, position int)
	DenyDirection(in string, out string)
	Deny(in string, out string)
	Heal()
}

type UnixIpTables struct {
}

func (ip UnixIpTables) PartitionLevel(nodes []string, position int) {
	for i := 0; i < position; i++ {
		for j := position; j < len(nodes); j++ {
			n1 := nodes[i]
			n2 := nodes[j]
			ip.Deny(n1, n2)
		}
	}
}

func (ip UnixIpTables) DenyDirection(incoming string, outgoing string) {
	// Should we be doing this for OUTPUT too?
	iptables(fmt.Sprintf("-A INPUT -j DROP -s %s -d %s", incoming, outgoing))
}

func (ip UnixIpTables) Deny(incoming string, outgoing string) {
	ip.DenyDirection(incoming, outgoing)
	ip.DenyDirection(outgoing, incoming)
}

func (ip UnixIpTables) Heal() {
	iptables("-F")
	iptables("-X")
}
