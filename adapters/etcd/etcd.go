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
	peer_addresses   []string
	client_addresses []string // for clients
}

func (c *EtcdAdapter) Addresses() []string {
	return c.client_addresses
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

	out_file, err := os.Create(fmt.Sprintf("log/etcd_%d.out", i))
	check(err)
	err_file, err := os.Create(fmt.Sprintf("log/etcd_%d.err", i))
	check(err)

	out_writer := bufio.NewWriter(out_file)
	defer out_writer.Flush()
	err_writer := bufio.NewWriter(err_file)
	defer err_writer.Flush()

	err = cmd.Start()
	check(err)

	go io.Copy(out_writer, stdoutPipe)
	go io.Copy(err_writer, stderrPipe)

	cmd.Wait()
}

func (c *EtcdAdapter) Setup() {
	var cmd *exec.Cmd
	for i, address := range c.peer_addresses {
		cmd = exec.Command("etcd", "--initial-cluster", address)
		ExecWithLog(cmd, i) // will panic if something goes wrong.
	}
}

func (c *EtcdAdapter) Teardown() {
}

func Cluster() *EtcdAdapter {
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
		peer_addresses[i] = host + PeerPort
		client_addresses[i] = host + ClientPort
	}

	return &EtcdAdapter{peer_addresses, client_addresses}
}