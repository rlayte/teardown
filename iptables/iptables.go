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
		log.Fatal(err)
	}

	return out
}

func Partition(nodes []string, position int) {
	for i := 0; i < position; i++ {
		for j := position; j < len(nodes); j++ {
			n1 := nodes[i]
			n2 := nodes[j]

			iptables(fmt.Sprintf("-A INPUT -j DROP -s %s -d %s", n1, n2))
		}
	}
}

func Heal() {
	iptables("-F")
	iptables("-X")
}
