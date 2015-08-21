package teardown

type Nemesis struct {
}

func (n *Nemesis) Partition(cluster Cluster) bool {
	return false
}

func (n *Nemesis) Heal(cluster Cluster) bool {
	return false
}
