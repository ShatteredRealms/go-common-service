// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/repository.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/repository.go -destination=/home/wil/dev/sro/go-common-service/pkg/bus/chat/channelbus/mocks/repository.go
//

// Package mock_channelbus is a generated GoMock package.
package mock_channelbus

import (
	context "context"
	reflect "reflect"

	channelbus "github.com/ShatteredRealms/go-common-service/pkg/bus/chat/channelbus"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockRepository) GetAll(ctx context.Context) (*channelbus.Channels, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].(*channelbus.Channels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRepositoryMockRecorder) GetAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepository)(nil).GetAll), ctx)
}

// GetByDimensionId mocks base method.
func (m *MockRepository) GetByDimensionId(ctx context.Context, ownerId *uuid.UUID) (*channelbus.Channels, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByDimensionId", ctx, ownerId)
	ret0, _ := ret[0].(*channelbus.Channels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByDimensionId indicates an expected call of GetByDimensionId.
func (mr *MockRepositoryMockRecorder) GetByDimensionId(ctx, ownerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByDimensionId", reflect.TypeOf((*MockRepository)(nil).GetByDimensionId), ctx, ownerId)
}

// GetById mocks base method.
func (m *MockRepository) GetById(ctx context.Context, channelId *uuid.UUID) (*channelbus.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", ctx, channelId)
	ret0, _ := ret[0].(*channelbus.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockRepositoryMockRecorder) GetById(ctx, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockRepository)(nil).GetById), ctx, channelId)
}

// Save mocks base method.
func (m *MockRepository) Save(ctx context.Context, data channelbus.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryMockRecorder) Save(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), ctx, data)
}
