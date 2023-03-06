// Code generated by MockGen. DO NOT EDIT.
// Source: node.go

// Package models is a generated GoMock package.
package models

import (
	reflect "reflect"

	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gomock "github.com/golang/mock/gomock"
	substrate "github.com/threefoldtech/substrate-client"
)

// MockSub is a mock of Sub interface.
type MockSub struct {
	ctrl     *gomock.Controller
	recorder *MockSubMockRecorder
}

// MockSubMockRecorder is the mock recorder for MockSub.
type MockSubMockRecorder struct {
	mock *MockSub
}

// NewMockSub creates a new mock instance.
func NewMockSub(ctrl *gomock.Controller) *MockSub {
	mock := &MockSub{ctrl: ctrl}
	mock.recorder = &MockSubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSub) EXPECT() *MockSubMockRecorder {
	return m.recorder
}

// SetNodePowerState mocks base method.
func (m *MockSub) SetNodePowerState(identity substrate.Identity, up bool) (types.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNodePowerState", identity, up)
	ret0, _ := ret[0].(types.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetNodePowerState indicates an expected call of SetNodePowerState.
func (mr *MockSubMockRecorder) SetNodePowerState(identity, up interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNodePowerState", reflect.TypeOf((*MockSub)(nil).SetNodePowerState), identity, up)
}

// GetNodeRentContract mocks base method.
func (m *MockSub) GetNodeRentContract(node uint32) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeRentContract", node)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeRentContract indicates an expected call of GetNodeRentContract.
func (mr *MockSubMockRecorder) GetNodeRentContract(node uint32) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeRentContract", reflect.TypeOf((*MockSub)(nil).GetNodeRentContract), node)
}
