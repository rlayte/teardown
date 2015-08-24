// server.go

package server

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
	peer_addresses   []string
	client_addresses []string // for clients
	processes        []*os.Process
	launchProcess    chan *os.Process
	killAll          chan bool
}

func (c *EtcdAdapter) Addresses() []string {
	return c.client_addresses
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

	out_file, err := os.Create(fmt.Sprintf("/tmp/etcd_%d.out", i))
	check(err)
	err_file, err := os.Create(fmt.Sprintf("/tmp/etcd_%d.err", i))
	check(err)

	out_writer := bufio.NewWriter(out_file)
	err_writer := bufio.NewWriter(err_file)

	defer out_writer.Flush()
	defer err_writer.Flush()

	// Start the command
	err = cmd.Start()
	check(err)

	go io.Copy(out_writer, stdoutPipe)
	go io.Copy(err_writer, stderrPipe)

	c.launchProcess <- cmd.Process

	cmd.Wait()
}

func (c *EtcdAdapter) nameFromPeer(peer string) string {
	for i, address := range c.peer_addresses {
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

	var all_peers string

	go c.serveProcesses()

	for i, peer_address := range c.peer_addresses {
		if i != 0 {
			all_peers += ","
		}
		all_peers += name(i) + "=" + peer_address
	}

	for i, peer_address := range c.peer_addresses {
		client_address := c.client_addresses[i]
		cmd = exec.Command(
			"etcd",
			"--name", c.nameFromPeer(peer_address),
			"--listen-peer-urls", peer_address,
			"--initial-advertise-peer-urls", peer_address,
			"--listen-client-urls", client_address,
			"--advertise-client-urls", client_address,
			"--initial-cluster", all_peers,
			"--initial-cluster-state", "new",
			"--initial-cluster-token", "etcd-teardown-cluster-1",
		)
		go c.ExecWithLog(cmd, i) // will panic if something goes wrong.
	}
}

func (c *EtcdAdapter) Teardown() {
	c.killAll <- true
	// Might be a race condition:

	for i := range c.peer_addresses {
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
		case random_bool := <-c.killAll:
			fmt.Println("2")
			random_bool = random_bool
			for _, p := range c.processes {
				p.Kill()
			}
		}
	}
}

func New() *EtcdAdapter {
	var peer_addresses, client_addresses []string
	hosts := []string{
		"127.0.0.12",
		"127.0.0.13",
		"127.0.0.14",
		"127.0.0.15",
		"127.0.0.16",
	}
	peer_addresses = make([]string, len(hosts))
	client_addresses = make([]string, len(hosts))

	for i, host := range hosts {
		peer_addresses[i] = "http://" + host + PeerPort
		client_addresses[i] = "http://" + host + ClientPort
	}

	return &EtcdAdapter{
		peer_addresses:   peer_addresses,
		client_addresses: client_addresses,
		processes:        []*os.Process{},
		launchProcess:    make(chan *os.Process),
		killAll:          make(chan bool),
	}
}
