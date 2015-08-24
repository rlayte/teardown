package teardown

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rlayte/teardown/mocks"
	"github.com/rlayte/teardown/mocks/iptables"
)

func TestPartitionHalf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addresses := []string{"a", "b", "c"}

	cluster := mock_teardown.NewMockCluster(ctrl)
	cluster.EXPECT().Addresses().Return(addresses)

	iptables := mock_iptables.NewMockIpTables(ctrl)
	iptables.EXPECT().PartitionLevel(addresses, 1)

	n := NewNemesis(cluster)
	n.iptables = iptables
	n.PartitionHalf()
}
