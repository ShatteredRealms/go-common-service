// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/go-common-service/pkg/pb/bus_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/dev/sro/go-common-service/pkg/pb/bus_grpc.pb.go -destination=/home/wil/dev/sro/go-common-service/pkg/pb/mocks/bus_grpc.pb.go
//

// Package mock_pb is a generated GoMock package.
package mock_pb

import (
	context "context"
	reflect "reflect"

	pb "github.com/ShatteredRealms/go-common-service/pkg/pb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockBusServiceClient is a mock of BusServiceClient interface.
type MockBusServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockBusServiceClientMockRecorder
	isgomock struct{}
}

// MockBusServiceClientMockRecorder is the mock recorder for MockBusServiceClient.
type MockBusServiceClientMockRecorder struct {
	mock *MockBusServiceClient
}

// NewMockBusServiceClient creates a new mock instance.
func NewMockBusServiceClient(ctrl *gomock.Controller) *MockBusServiceClient {
	mock := &MockBusServiceClient{ctrl: ctrl}
	mock.recorder = &MockBusServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBusServiceClient) EXPECT() *MockBusServiceClientMockRecorder {
	return m.recorder
}

// ResetReaderBus mocks base method.
func (m *MockBusServiceClient) ResetReaderBus(ctx context.Context, in *pb.BusTarget, opts ...grpc.CallOption) (*pb.ResetBusResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ResetReaderBus", varargs...)
	ret0, _ := ret[0].(*pb.ResetBusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetReaderBus indicates an expected call of ResetReaderBus.
func (mr *MockBusServiceClientMockRecorder) ResetReaderBus(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetReaderBus", reflect.TypeOf((*MockBusServiceClient)(nil).ResetReaderBus), varargs...)
}

// ResetWriterBus mocks base method.
func (m *MockBusServiceClient) ResetWriterBus(ctx context.Context, in *pb.BusTarget, opts ...grpc.CallOption) (*pb.ResetBusResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ResetWriterBus", varargs...)
	ret0, _ := ret[0].(*pb.ResetBusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetWriterBus indicates an expected call of ResetWriterBus.
func (mr *MockBusServiceClientMockRecorder) ResetWriterBus(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetWriterBus", reflect.TypeOf((*MockBusServiceClient)(nil).ResetWriterBus), varargs...)
}

// MockBusServiceServer is a mock of BusServiceServer interface.
type MockBusServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockBusServiceServerMockRecorder
	isgomock struct{}
}

// MockBusServiceServerMockRecorder is the mock recorder for MockBusServiceServer.
type MockBusServiceServerMockRecorder struct {
	mock *MockBusServiceServer
}

// NewMockBusServiceServer creates a new mock instance.
func NewMockBusServiceServer(ctrl *gomock.Controller) *MockBusServiceServer {
	mock := &MockBusServiceServer{ctrl: ctrl}
	mock.recorder = &MockBusServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBusServiceServer) EXPECT() *MockBusServiceServerMockRecorder {
	return m.recorder
}

// ResetReaderBus mocks base method.
func (m *MockBusServiceServer) ResetReaderBus(arg0 context.Context, arg1 *pb.BusTarget) (*pb.ResetBusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetReaderBus", arg0, arg1)
	ret0, _ := ret[0].(*pb.ResetBusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetReaderBus indicates an expected call of ResetReaderBus.
func (mr *MockBusServiceServerMockRecorder) ResetReaderBus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetReaderBus", reflect.TypeOf((*MockBusServiceServer)(nil).ResetReaderBus), arg0, arg1)
}

// ResetWriterBus mocks base method.
func (m *MockBusServiceServer) ResetWriterBus(arg0 context.Context, arg1 *pb.BusTarget) (*pb.ResetBusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetWriterBus", arg0, arg1)
	ret0, _ := ret[0].(*pb.ResetBusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetWriterBus indicates an expected call of ResetWriterBus.
func (mr *MockBusServiceServerMockRecorder) ResetWriterBus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetWriterBus", reflect.TypeOf((*MockBusServiceServer)(nil).ResetWriterBus), arg0, arg1)
}

// mustEmbedUnimplementedBusServiceServer mocks base method.
func (m *MockBusServiceServer) mustEmbedUnimplementedBusServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedBusServiceServer")
}

// mustEmbedUnimplementedBusServiceServer indicates an expected call of mustEmbedUnimplementedBusServiceServer.
func (mr *MockBusServiceServerMockRecorder) mustEmbedUnimplementedBusServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedBusServiceServer", reflect.TypeOf((*MockBusServiceServer)(nil).mustEmbedUnimplementedBusServiceServer))
}

// MockUnsafeBusServiceServer is a mock of UnsafeBusServiceServer interface.
type MockUnsafeBusServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeBusServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeBusServiceServerMockRecorder is the mock recorder for MockUnsafeBusServiceServer.
type MockUnsafeBusServiceServerMockRecorder struct {
	mock *MockUnsafeBusServiceServer
}

// NewMockUnsafeBusServiceServer creates a new mock instance.
func NewMockUnsafeBusServiceServer(ctrl *gomock.Controller) *MockUnsafeBusServiceServer {
	mock := &MockUnsafeBusServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeBusServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeBusServiceServer) EXPECT() *MockUnsafeBusServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedBusServiceServer mocks base method.
func (m *MockUnsafeBusServiceServer) mustEmbedUnimplementedBusServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedBusServiceServer")
}

// mustEmbedUnimplementedBusServiceServer indicates an expected call of mustEmbedUnimplementedBusServiceServer.
func (mr *MockUnsafeBusServiceServerMockRecorder) mustEmbedUnimplementedBusServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedBusServiceServer", reflect.TypeOf((*MockUnsafeBusServiceServer)(nil).mustEmbedUnimplementedBusServiceServer))
}
