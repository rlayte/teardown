// cluster.go

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
	processes       []*os.Process
	launchProcess   chan *os.Process
	killAll         chan bool
}

func (c *EtcdAdapter) Addresses() []string {
	return c.clientAddresses
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (c *EtcdAdapter) ExecWithLog(cmd *exec.Cmd, i int) {
	stdoutPipe, err := cmd.StdoutPipe()
	check(err)
	stderrPipe, err := cmd.StderrPipe()
	check(err)

	outFile, err := os.Create(fmt.Sprintf("/tmp/etcd_%d.out", i))
	check(err)
	errFile, err := os.Create(fmt.Sprintf("/tmp/etcd_%d.err", i))
	check(err)

	outWriter := bufio.NewWriter(outFile)
	errWriter := bufio.NewWriter(errFile)

	defer outWriter.Flush()
	defer errWriter.Flush()

	// Start the command
	err = cmd.Start()
	check(err)

	go io.Copy(outWriter, stdoutPipe)
	go io.Copy(errWriter, stderrPipe)

	c.launchProcess <- cmd.Process

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

	go c.serveProcesses()

	for i, peerAddress := range c.peerAddresses {
		if i != 0 {
			allPeers += ","
		}
		allPeers += name(i) + "=" + peerAddress
	}

	for i, peerAddress := range c.peerAddresses {
		clientAddress := c.clientAddresses[i]
		cmd = exec.Command(
			"etcd",
			"--name", c.nameFromPeer(peerAddress),
			"--listen-peer-urls", peerAddress,
			"--initial-advertise-peer-urls", peerAddress,
			"--listen-client-urls", clientAddress,
			"--advertise-client-urls", clientAddress,
			"--initial-cluster", allPeers,
			"--initial-cluster-state", "new",
			"--initial-cluster-token", "etcd-teardown-cluster-1",
		)
		go c.ExecWithLog(cmd, i) // will panic if something goes wrong.
	}
}

func (c *EtcdAdapter) Teardown() {
	c.killAll <- true
	// Might be a race condition:

	for i := range c.peerAddresses {
		err := os.RemoveAll(name(i) + ".etcd")
		if err != nil {
			panic(err)
		}
	}

}

func (c *EtcdAdapter) serveProcesses() {
	var process *os.Process
	for {
		select {
		case process = <-c.launchProcess:
			c.processes = append(c.processes, process)
		case randomBool := <-c.killAll:
			fmt.Println("2")
			randomBool = randomBool
			for _, p := range c.processes {
				p.Kill()
			}
		}
	}
}

func NewEtcdCluster() *EtcdAdapter {
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

	cluster := &EtcdAdapter{
		peerAddresses:   peerAddresses,
		clientAddresses: clientAddresses,
		processes:       []*os.Process{},
		launchProcess:   make(chan *os.Process),
		killAll:         make(chan bool),
	}

	cluster.Setup()

	return cluster
}
