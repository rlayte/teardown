package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/coreos/go-etcd/etcd"
	"github.com/rlayte/teardown"
)

const (
	ClientPort string = ":4000"
	PeerPort   string = ":4004"
)

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

type EtcdAdapter struct {
	peerAddresses   []string
	clientAddresses []string // for clients
}

func (c *EtcdAdapter) Addresses() []string {
	return c.clientAddresses
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

type EtcdClient struct {
	client *etcd.Client
}

func (c EtcdClient) Get(key string) (string, error) {
	_, err := c.client.Get(key, false, false)
	return "", err
}

func (c EtcdClient) Put(key string, value string) error {
	_, err := c.client.Set(key, value, 0)
	return err
}

func main() {
	cluster := Cluster()
	client := EtcdClient{etcd.NewClient(cluster.Addresses())}
	tests := teardown.NewTestRunner(cluster, client)
	tests.Run()
}
