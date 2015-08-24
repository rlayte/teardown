// Automatically generated by MockGen. DO NOT EDIT!
// Source: iptables/iptables.go

package mock_iptables

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of IpTables interface
type MockIpTables struct {
	ctrl     *gomock.Controller
	recorder *_MockIpTablesRecorder
}

// Recorder for MockIpTables (not exported)
type _MockIpTablesRecorder struct {
	mock *MockIpTables
}

func NewMockIpTables(ctrl *gomock.Controller) *MockIpTables {
	mock := &MockIpTables{ctrl: ctrl}
	mock.recorder = &_MockIpTablesRecorder{mock}
	return mock
}

func (_m *MockIpTables) EXPECT() *_MockIpTablesRecorder {
	return _m.recorder
}

func (_m *MockIpTables) PartitionLevel(nodes []string, position int) {
	_m.ctrl.Call(_m, "PartitionLevel", nodes, position)
}

func (_mr *_MockIpTablesRecorder) PartitionLevel(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PartitionLevel", arg0, arg1)
}

func (_m *MockIpTables) DenyDirection(in string, out string) {
	_m.ctrl.Call(_m, "DenyDirection", in, out)
}

func (_mr *_MockIpTablesRecorder) DenyDirection(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DenyDirection", arg0, arg1)
}

func (_m *MockIpTables) Deny(in string, out string) {
	_m.ctrl.Call(_m, "Deny", in, out)
}

func (_mr *_MockIpTablesRecorder) Deny(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Deny", arg0, arg1)
}

func (_m *MockIpTables) Heal() {
	_m.ctrl.Call(_m, "Heal")
}

func (_mr *_MockIpTablesRecorder) Heal() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Heal")
}
