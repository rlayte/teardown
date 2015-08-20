package teardown

type Client interface {
	Get(key string) (value string, err error)
	Put(key string, value string) error
}
