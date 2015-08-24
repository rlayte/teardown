# Teardown

Teardown is a library for testing distributed systems --- specifically distributed key/value stores. We were heavily inspired by [@aphyr's](http://aphyr.com) [Jepsen](https://github.com/aphyr/jepsen), but teardown differs in a number of ways:

- Written in Go.
- Doesn't bundle installation into tests.
- Only tests simple consistency models (i.e. doesn't do the difficult linearizability stuff).

## Getting started

To use teardown we make a few assumptions about your system:

- Hosts 127.0.0.1/24 are available to bind to. (Default on Linux)
- You have [iptables](https://en.wikipedia.com/wiki/iptables) installed.
- You have [tc](http://lartc.org/manpages/tc.txt) installed. (`sudo apt-get install iproute`)
- You have whatever you're testing installed.

Because teardown messes with some core networking options it's probably safest to run in a sandbox environment like Docker.

### Install

    $ go get github.com/rlayte/teardown

### Setup

teardown exposes a `Cluster` interfaces that you must implement to with the specific details of your system.

`Cluster` exposes three methods --- `Setup`, `Teardown` and `Addresses`.

```go
// mycluster.go

package mycluster

type MyCluster struct {
  addresses []string
}

func (c *MyCluster) Setup() {
  // Calls start on an imaginary service
  c.addresses = []string{"127.0.0.2", "127.0.0.3", "127.0.0.4"}
  err := exec.Command("mycluster", "start", c.addresses).Start()
  if err != nil {
    // handle error
  }
}

func (c *MyCluster) Teardown() {
  // Calls stop on an imaginary service
  err := exec.Command("mything", "stop", c.addresses).Run()
  if err != nil {
    // handle error
  }
}

func (c *MyCluster) Addresses() []string {
  // Returns a list of all node addresses
  return c.addresses
}
```

### Running the tests

Once you have a concrete implementation of `Cluster` you can pass it to a `Nemesis` instance, manipulate the current state of the network and write tests as normal. E.g.

```go
// mycluster_test.go

package mycluster

var cluster teardown.Cluster
var nemesis teardown.Nemesis

init () {
  cluster = NewMyCluster()
  nemesis = teardown.NewNemesis(cluster)
}

func TestPartition(t *testing.T) {
  nemesis.PartitionHalf()

  // Run tests here
}
```

And run this with:

    $ go test

## Documentation

[Full API documentation]()
