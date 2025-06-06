// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/service.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/service.go -destination=/home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/mocks/service.go
//

// Package mock_channelbus is a generated GoMock package.
package mock_channelbus

import (
	context "context"
	reflect "reflect"

	bus "github.com/ShatteredRealms/go-common-service/pkg/bus"
	channelbus "github.com/ShatteredRealms/go-common-service/pkg/bus/chat/channelbus"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetChannelById mocks base method.
func (m *MockService) GetChannelById(ctx context.Context, channelId string) (*channelbus.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannelById", ctx, channelId)
	ret0, _ := ret[0].(*channelbus.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannelById indicates an expected call of GetChannelById.
func (mr *MockServiceMockRecorder) GetChannelById(ctx, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelById", reflect.TypeOf((*MockService)(nil).GetChannelById), ctx, channelId)
}

// GetChannels mocks base method.
func (m *MockService) GetChannels(ctx context.Context) (*channelbus.Channels, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannels", ctx)
	ret0, _ := ret[0].(*channelbus.Channels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannels indicates an expected call of GetChannels.
func (mr *MockServiceMockRecorder) GetChannels(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannels", reflect.TypeOf((*MockService)(nil).GetChannels), ctx)
}

// GetResetter mocks base method.
func (m *MockService) GetResetter() bus.Resettable {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResetter")
	ret0, _ := ret[0].(bus.Resettable)
	return ret0
}

// GetResetter indicates an expected call of GetResetter.
func (mr *MockServiceMockRecorder) GetResetter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResetter", reflect.TypeOf((*MockService)(nil).GetResetter))
}

// IsProcessing mocks base method.
func (m *MockService) IsProcessing() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsProcessing")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsProcessing indicates an expected call of IsProcessing.
func (mr *MockServiceMockRecorder) IsProcessing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsProcessing", reflect.TypeOf((*MockService)(nil).IsProcessing))
}

// RegisterListener mocks base method.
func (m *MockService) RegisterListener(listener bus.BusListener[channelbus.Message]) bus.BusListenerHandler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterListener", listener)
	ret0, _ := ret[0].(bus.BusListenerHandler)
	return ret0
}

// RegisterListener indicates an expected call of RegisterListener.
func (mr *MockServiceMockRecorder) RegisterListener(listener any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterListener", reflect.TypeOf((*MockService)(nil).RegisterListener), listener)
}

// RemoveListener mocks base method.
func (m *MockService) RemoveListener(listenerHandle bus.BusListenerHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveListener", listenerHandle)
}

// RemoveListener indicates an expected call of RemoveListener.
func (mr *MockServiceMockRecorder) RemoveListener(listenerHandle any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveListener", reflect.TypeOf((*MockService)(nil).RemoveListener), listenerHandle)
}

// StartProcessing mocks base method.
func (m *MockService) StartProcessing(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartProcessing", ctx)
}

// StartProcessing indicates an expected call of StartProcessing.
func (mr *MockServiceMockRecorder) StartProcessing(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartProcessing", reflect.TypeOf((*MockService)(nil).StartProcessing), ctx)
}

// StopProcessing mocks base method.
func (m *MockService) StopProcessing() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopProcessing")
}

// StopProcessing indicates an expected call of StopProcessing.
func (mr *MockServiceMockRecorder) StopProcessing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopProcessing", reflect.TypeOf((*MockService)(nil).StopProcessing))
}
