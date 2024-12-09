// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/dev/sro/go-common-service/pkg/bus/gameserver/mapbus/repository.go
//
// Generated by this command:
//
//	mockgen -source=/home/wil/dev/sro/go-common-service/pkg/bus/gameserver/mapbus/repository.go -destination=/home/wil/dev/sro/go-common-service/pkg/bus/gameserver/mapbus/mocks/repository.go
//

// Package mock_mapbus is a generated GoMock package.
package mock_mapbus

import (
	context "context"
	reflect "reflect"

	mapbus "github.com/ShatteredRealms/go-common-service/pkg/bus/gameserver/mapbus"
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

// CreateMap mocks base method.
func (m *MockRepository) CreateMap(ctx context.Context, dimensionId string) (*mapbus.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMap", ctx, dimensionId)
	ret0, _ := ret[0].(*mapbus.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMap indicates an expected call of CreateMap.
func (mr *MockRepositoryMockRecorder) CreateMap(ctx, dimensionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMap", reflect.TypeOf((*MockRepository)(nil).CreateMap), ctx, dimensionId)
}

// DeleteMap mocks base method.
func (m *MockRepository) DeleteMap(ctx context.Context, dimensionId string) (*mapbus.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMap", ctx, dimensionId)
	ret0, _ := ret[0].(*mapbus.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMap indicates an expected call of DeleteMap.
func (mr *MockRepositoryMockRecorder) DeleteMap(ctx, dimensionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMap", reflect.TypeOf((*MockRepository)(nil).DeleteMap), ctx, dimensionId)
}

// GetMapById mocks base method.
func (m *MockRepository) GetMapById(ctx context.Context, dimensionId string) (*mapbus.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMapById", ctx, dimensionId)
	ret0, _ := ret[0].(*mapbus.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMapById indicates an expected call of GetMapById.
func (mr *MockRepositoryMockRecorder) GetMapById(ctx, dimensionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMapById", reflect.TypeOf((*MockRepository)(nil).GetMapById), ctx, dimensionId)
}

// GetMaps mocks base method.
func (m *MockRepository) GetMaps(ctx context.Context) (*mapbus.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaps", ctx)
	ret0, _ := ret[0].(*mapbus.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMaps indicates an expected call of GetMaps.
func (mr *MockRepositoryMockRecorder) GetMaps(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaps", reflect.TypeOf((*MockRepository)(nil).GetMaps), ctx)
}
