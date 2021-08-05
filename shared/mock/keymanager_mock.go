// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1 (interfaces: RemoteSignerClient)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockRemoteSignerClient is a mock of RemoteSignerClient interface.
type MockRemoteSignerClient struct {
	ctrl     *gomock.Controller
	recorder *MockRemoteSignerClientMockRecorder
}

// MockRemoteSignerClientMockRecorder is the mock recorder for MockRemoteSignerClient.
type MockRemoteSignerClientMockRecorder struct {
	mock *MockRemoteSignerClient
}

// NewMockRemoteSignerClient creates a new mock instance.
func NewMockRemoteSignerClient(ctrl *gomock.Controller) *MockRemoteSignerClient {
	mock := &MockRemoteSignerClient{ctrl: ctrl}
	mock.recorder = &MockRemoteSignerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRemoteSignerClient) EXPECT() *MockRemoteSignerClientMockRecorder {
	return m.recorder
}

// ListValidatingPublicKeys mocks base method
func (m *MockRemoteSignerClient) ListValidatingPublicKeys(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v2.ListPublicKeysResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListValidatingPublicKeys", varargs...)
	ret0, _ := ret[0].(*v2.ListPublicKeysResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListValidatingPublicKeys indicates an expected call of ListValidatingPublicKeys.
func (mr *MockRemoteSignerClientMockRecorder) ListValidatingPublicKeys(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListValidatingPublicKeys", reflect.TypeOf((*MockRemoteSignerClient)(nil).ListValidatingPublicKeys), varargs...)
}

// Sign mocks base method
func (m *MockRemoteSignerClient) Sign(arg0 context.Context, arg1 *v2.SignRequest, arg2 ...grpc.CallOption) (*v2.SignResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Sign", varargs...)
	ret0, _ := ret[0].(*v2.SignResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign.
func (mr *MockRemoteSignerClientMockRecorder) Sign(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockRemoteSignerClient)(nil).Sign), varargs...)
}
