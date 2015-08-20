// adapters/etcd.go

package etcd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const (
	ClientPort string = ":4000"
	PeerPort   string = ":4004"
)

type EtcdAdapter struct {
	peerAddresses   []string
	clientAddresses []string // for clients
}

func (c *EtcdAdapter) Addresses() []string {
	return c.clientAddresses
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ExecWithLog(cmd *exec.Cmd, i int) {
	stdoutPipe, err := cmd.StdoutPipe()
	check(err)
	stderrPipe, err := cmd.StderrPipe()
	check(err)

	outFile, err := os.Create(fmt.Sprintf("log/etcd_%d.out", i))
	check(err)
	errFile, err := os.Create(fmt.Sprintf("log/etcd_%d.err", i))
	check(err)

	outWriter := bufio.NewWriter(outFile)
	defer outWriter.Flush()
	errWriter := bufio.NewWriter(errFile)
	defer errWriter.Flush()

	err = cmd.Start()
	check(err)

	go io.Copy(outWriter, stdoutPipe)
	go io.Copy(errWriter, stderrPipe)

	cmd.Wait()
}

func (c *EtcdAdapter) nameFromPeer(peer string) string {
	for i, address := range c.peerAddresses {
		if peer == address {
			return name(i)
		}
	}
	panic(fmt.Sprintf("No such peer: %s", peer))
}

func name(i int) string {
	return fmt.Sprintf("peer%d", i)
}

func (c *EtcdAdapter) Setup() {
	var cmd *exec.Cmd

	var allPeers string

	for i, peerAddress := range c.peerAddresses {
		if i != 0 {
			allPeers += ","
		}
		allPeers += name(i) + "=" + peerAddress
	}

	for i, peer_address := range c.peerAddresses {
		client_address := c.clientAddresses[i]
		cmd = exec.Command(
			"etcd",
			"--name", c.nameFromPeer(peer_address),
			"--listen-peer-urls", peer_address,
			"--initial-advertise-peer-urls", peer_address,
			"--listen-client-urls", client_address,
			"--advertise-client-urls", client_address,
			"--initial-cluster", allPeers,
			"--initial-cluster-state", "new",
			"--initial-cluster-token", "etcd-teardown-cluster-1",
		)
		go ExecWithLog(cmd, i) // will panic if something goes wrong.
	}
}

func (c *EtcdAdapter) Teardown() {
}

func Cluster() *EtcdAdapter {
	var peerAddresses, clientAddresses []string
	hosts := []string{
		"127.0.0.12",
		"127.0.0.13",
		"127.0.0.14",
		"127.0.0.15",
		"127.0.0.16",
	}
	peerAddresses = make([]string, len(hosts))
	clientAddresses = make([]string, len(hosts))

	for i, host := range hosts {
		peerAddresses[i] = "http://" + host + PeerPort
		clientAddresses[i] = "http://" + host + ClientPort
	}

	return &EtcdAdapter{peerAddresses, clientAddresses}
}
