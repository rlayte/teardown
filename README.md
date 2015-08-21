# Teardown

Teardown is a library for testing distributed systems --- specifically distributed key/value stores. We were heavily inspired by [@aphyr's]() [Jepsen](), but teardown differs in a number of ways:

- Written in Go.
- Doesn't bundle installation into tests.
- Only tests simple consistency models (i.e. doesn't do the difficult linearizability stuff).

## Getting started

To use teardown we make a few assumptions about your system:

- Hosts 127.0.0.1/24 are available to bind to.
- You have [iptables]() installed.
- You have [tc]() installed.
- You have whatever you're testing installed.

Because teardown messes with some core networking options it's probably safest to run in a sandbox environment like Docker.

### Install

    $ go get github.com/rlayte/teardown

### Setup

teardown exposes two interfaces, `Cluster` and `Client`, that you must implement to with the specific details of your system.

`Cluster` exposes three methods --- `Setup`, `Teardown` and `Addresses`.

```go
// mycluster.go

type MyCluster struct {
  addresses []string
}

func (c *MyCluster) Setup() {
  // Calls start on an imaginary service
  c.addresses = []string{"127.0.0.2", "127.0.0.3", "127.0.0.4"}
  exec.Command("mycluster", "start", c.addresses)
}

func (c *MyCluster) Teardown() {
  // Calls stop on an imaginary service
  exec.Command("mything", "stop", c.addresses)
}

func (c *MyCluster) Addresses() []string {
  // Returns a list of all node addresses
  return c.addresses
}
```

`Client` handles getting and setting values in your system. Here's a very simple example that uses http:

```go
// mycluster.go

type MyClient struct {
  cluster Cluster
}

func (c *MyClient) currentLeader() string {
  // Returns a random node in the cluster
  addresses := c.cluster.Addresses()
  return addresses[rand.Intn(len(addresses))]
}

func (c *MyClient) Put(key string, value string) error {
  // Sends value=bar to http://127.0.0.z:4000/:key
  address := fmt.Sprintf("http://%s:4000/%s", c.currentLeader(), key)
  data := url.Values{}
  data.Set("value", "bar")
  resp, err := http.PostForm(address, data)
  return err
}

func (c *MyClient) Get(key string) ([]byte, error) {
  // Returns the body of http://127.0.0.x:4000/:key
  address := fmt.Sprintf("http://%s:4000/%s", c.currentLeader(), key)
  resp, err := http.Get(address) 
  defer resp.Body.Close()
  return ioutils.ReadAll(resp.Body), err
}
```

### Running the tests

Once you have a concrete implementation of `Cluster` and `Client` pass them to `NewTestRunner` to execute the tests. E.g.

```go
// mycluster.go

main () {
  cluster := NewMyCluster()
  client := NewMyClient(cluster)

  tests := teardown.NewTestRunner(cluster, client)
  tests.Run()
}
```

And run this with:

    $ go run mycluster.go

## Documentation

[Full API documentation]()
