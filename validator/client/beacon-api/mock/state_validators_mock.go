// Code generated by MockGen. DO NOT EDIT.
// Source: validator/client/beacon-api/state_validators.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/api/eth/beacon"
	primitives "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

// MockstateValidatorsProvider is a mock of stateValidatorsProvider interface.
type MockstateValidatorsProvider struct {
	ctrl     *gomock.Controller
	recorder *MockstateValidatorsProviderMockRecorder
}

// MockstateValidatorsProviderMockRecorder is the mock recorder for MockstateValidatorsProvider.
type MockstateValidatorsProviderMockRecorder struct {
	mock *MockstateValidatorsProvider
}

// NewMockstateValidatorsProvider creates a new mock instance.
func NewMockstateValidatorsProvider(ctrl *gomock.Controller) *MockstateValidatorsProvider {
	mock := &MockstateValidatorsProvider{ctrl: ctrl}
	mock.recorder = &MockstateValidatorsProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstateValidatorsProvider) EXPECT() *MockstateValidatorsProviderMockRecorder {
	return m.recorder
}

// GetStateValidators mocks base method.
func (m *MockstateValidatorsProvider) GetStateValidators(arg0 context.Context, arg1 []string, arg2 []int64, arg3 []string) (*beacon.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidators", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*beacon.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidators indicates an expected call of GetStateValidators.
func (mr *MockstateValidatorsProviderMockRecorder) GetStateValidators(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidators", reflect.TypeOf((*MockstateValidatorsProvider)(nil).GetStateValidators), arg0, arg1, arg2, arg3)
}

// GetStateValidatorsForHead mocks base method.
func (m *MockstateValidatorsProvider) GetStateValidatorsForHead(arg0 context.Context, arg1 []string, arg2 []primitives.ValidatorIndex, arg3 []string) (*beacon.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidatorsForHead", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*beacon.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidatorsForHead indicates an expected call of GetStateValidatorsForHead.
func (mr *MockstateValidatorsProviderMockRecorder) GetStateValidatorsForHead(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidatorsForHead", reflect.TypeOf((*MockstateValidatorsProvider)(nil).GetStateValidatorsForHead), arg0, arg1, arg2, arg3)
}

// GetStateValidatorsForSlot mocks base method.
func (m *MockstateValidatorsProvider) GetStateValidatorsForSlot(arg0 context.Context, arg1 primitives.Slot, arg2 []string, arg3 []primitives.ValidatorIndex, arg4 []string) (*beacon.GetValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateValidatorsForSlot", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*beacon.GetValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateValidatorsForSlot indicates an expected call of GetStateValidatorsForSlot.
func (mr *MockstateValidatorsProviderMockRecorder) GetStateValidatorsForSlot(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateValidatorsForSlot", reflect.TypeOf((*MockstateValidatorsProvider)(nil).GetStateValidatorsForSlot), arg0, arg1, arg2, arg3, arg4)
}
