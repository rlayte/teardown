package teardown

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rlayte/teardown/mocks"
	"github.com/rlayte/teardown/mocks/iptables"
)

func setup(t *testing.T, addresses []string) (*gomock.Controller, *mock_teardown.MockCluster,
	*mock_iptables.MockIpTables) {
	ctrl := gomock.NewController(t)

	cluster := mock_teardown.NewMockCluster(ctrl)
	iptables := mock_iptables.NewMockIpTables(ctrl)

	return ctrl, cluster, iptables
}

func TestPartitionHalf(t *testing.T) {
	addresses := []string{"a", "b", "c"}
	ctrl, cluster, iptables := setup(t, addresses)
	defer ctrl.Finish()

	cluster.EXPECT().Addresses().Return(addresses)
	iptables.EXPECT().PartitionLevel(addresses, 1)

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.PartitionHalf()
}

func TestPartitionRandom(t *testing.T) {
	addresses := []string{"a", "b", "c"}
	ctrl, cluster, iptables := setup(t, addresses)
	defer ctrl.Finish()

	cluster.EXPECT().Addresses().Return(addresses)
	iptables.EXPECT().Deny("c", "a")
	iptables.EXPECT().Deny("c", "b")

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.PartitionRandom()
}

func TestPartitionSingle(t *testing.T) {
	addresses := []string{"a", "b", "c"}
	ctrl, cluster, iptables := setup(t, addresses)
	defer ctrl.Finish()

	cluster.EXPECT().Addresses().Return(addresses)
	iptables.EXPECT().Deny("a", "b")
	iptables.EXPECT().Deny("a", "c")

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.PartitionSingle(0)
}

func TestBridge(t *testing.T) {
	addresses := []string{"a", "b", "c", "d"}
	ctrl, cluster, iptables := setup(t, addresses)
	defer ctrl.Finish()

	cluster.EXPECT().Addresses().Return(addresses)
	iptables.EXPECT().Deny("a", "d")
	iptables.EXPECT().Deny("b", "d")

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.Bridge()
}

func TestHeal(t *testing.T) {
	addresses := []string{"a", "b", "c"}
	ctrl, cluster, iptables := setup(t, addresses)
	defer ctrl.Finish()

	cluster.EXPECT().Addresses().Return(addresses)
	iptables.EXPECT().Heal()

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.Heal()
}
