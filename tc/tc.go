package tc

import (
	"log"
	"os/exec"
	"strings"
)

func tc(args string) []byte {
	out, err := exec.Command("tc", strings.Split(args, " ")...).Output()

	if err != nil {
		log.Fatal(err)
	}

	return out
}

func Slow() {
	tc("qdisc add dev eth0 root netem delay 200ms")
}

func Flaky() {
	tc("qdisc change dev eth0 root netem loss 20% 75%")
}

func Fast() {
	tc("qdisc del dev eth0 root")
}
