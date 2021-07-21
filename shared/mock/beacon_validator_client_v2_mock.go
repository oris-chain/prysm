// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/prysm/v2 (interfaces: BeaconNodeValidatorAltairClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockBeaconNodeValidatorAltairClient is a mock of BeaconNodeValidatorAltairClient interface
type MockBeaconNodeValidatorAltairClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconNodeValidatorAltairClientMockRecorder
}

// MockBeaconNodeValidatorAltairClientMockRecorder is the mock recorder for MockBeaconNodeValidatorAltairClient
type MockBeaconNodeValidatorAltairClientMockRecorder struct {
	mock *MockBeaconNodeValidatorAltairClient
}

// NewMockBeaconNodeValidatorAltairClient creates a new mock instance
func NewMockBeaconNodeValidatorAltairClient(ctrl *gomock.Controller) *MockBeaconNodeValidatorAltairClient {
	mock := &MockBeaconNodeValidatorAltairClient{ctrl: ctrl}
	mock.recorder = &MockBeaconNodeValidatorAltairClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconNodeValidatorAltairClient) EXPECT() *MockBeaconNodeValidatorAltairClientMockRecorder {
	return m.recorder
}

// GetBlock mocks base method
func (m *MockBeaconNodeValidatorAltairClient) GetBlock(arg0 context.Context, arg1 *v1alpha1.BlockRequest, arg2 ...grpc.CallOption) (*v2.BeaconBlockAltair, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBlock", varargs...)
	ret0, _ := ret[0].(*v2.BeaconBlockAltair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) GetBlock(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).GetBlock), varargs...)
}

// GetSyncCommitteeContribution mocks base method
func (m *MockBeaconNodeValidatorAltairClient) GetSyncCommitteeContribution(arg0 context.Context, arg1 *v2.SyncCommitteeContributionRequest, arg2 ...grpc.CallOption) (*v2.SyncCommitteeContribution, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSyncCommitteeContribution", varargs...)
	ret0, _ := ret[0].(*v2.SyncCommitteeContribution)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncCommitteeContribution indicates an expected call of GetSyncCommitteeContribution
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) GetSyncCommitteeContribution(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncCommitteeContribution", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).GetSyncCommitteeContribution), varargs...)
}

// GetSyncMessageBlockRoot mocks base method
func (m *MockBeaconNodeValidatorAltairClient) GetSyncMessageBlockRoot(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v2.SyncMessageBlockRootResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSyncMessageBlockRoot", varargs...)
	ret0, _ := ret[0].(*v2.SyncMessageBlockRootResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncMessageBlockRoot indicates an expected call of GetSyncMessageBlockRoot
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) GetSyncMessageBlockRoot(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncMessageBlockRoot", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).GetSyncMessageBlockRoot), varargs...)
}

// GetSyncSubcommitteeIndex mocks base method
func (m *MockBeaconNodeValidatorAltairClient) GetSyncSubcommitteeIndex(arg0 context.Context, arg1 *v2.SyncSubcommitteeIndexRequest, arg2 ...grpc.CallOption) (*v2.SyncSubcommitteeIndexResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSyncSubcommitteeIndex", varargs...)
	ret0, _ := ret[0].(*v2.SyncSubcommitteeIndexResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncSubcommitteeIndex indicates an expected call of GetSyncSubcommitteeIndex
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) GetSyncSubcommitteeIndex(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncSubcommitteeIndex", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).GetSyncSubcommitteeIndex), varargs...)
}

// ProposeBlock mocks base method
func (m *MockBeaconNodeValidatorAltairClient) ProposeBlock(arg0 context.Context, arg1 *v2.SignedBeaconBlockAltair, arg2 ...grpc.CallOption) (*v1alpha1.ProposeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ProposeBlock", varargs...)
	ret0, _ := ret[0].(*v1alpha1.ProposeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProposeBlock indicates an expected call of ProposeBlock
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) ProposeBlock(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProposeBlock", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).ProposeBlock), varargs...)
}

// StreamBlocks mocks base method
func (m *MockBeaconNodeValidatorAltairClient) StreamBlocks(arg0 context.Context, arg1 *v1alpha1.StreamBlocksRequest, arg2 ...grpc.CallOption) (v2.BeaconNodeValidatorAltair_StreamBlocksClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StreamBlocks", varargs...)
	ret0, _ := ret[0].(v2.BeaconNodeValidatorAltair_StreamBlocksClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamBlocks indicates an expected call of StreamBlocks
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) StreamBlocks(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamBlocks", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).StreamBlocks), varargs...)
}

// SubmitSignedContributionAndProof mocks base method
func (m *MockBeaconNodeValidatorAltairClient) SubmitSignedContributionAndProof(arg0 context.Context, arg1 *v2.SignedContributionAndProof, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubmitSignedContributionAndProof", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitSignedContributionAndProof indicates an expected call of SubmitSignedContributionAndProof
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) SubmitSignedContributionAndProof(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitSignedContributionAndProof", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).SubmitSignedContributionAndProof), varargs...)
}

// SubmitSyncMessage mocks base method
func (m *MockBeaconNodeValidatorAltairClient) SubmitSyncMessage(arg0 context.Context, arg1 *v2.SyncCommitteeMessage, arg2 ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubmitSyncMessage", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitSyncMessage indicates an expected call of SubmitSyncMessage
func (mr *MockBeaconNodeValidatorAltairClientMockRecorder) SubmitSyncMessage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitSyncMessage", reflect.TypeOf((*MockBeaconNodeValidatorAltairClient)(nil).SubmitSyncMessage), varargs...)
}
