// Code generated by MockGen. DO NOT EDIT.
// Source: ../usecase.go

// Package mock_subscription is a generated GoMock package.
package mock_subscription

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// AddPlanning mocks base method.
func (m *MockUseCase) AddPlanning(userId, eid uint64) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlanning", userId, eid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddPlanning indicates an expected call of AddPlanning.
func (mr *MockUseCaseMockRecorder) AddPlanning(userId, eid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlanning", reflect.TypeOf((*MockUseCase)(nil).AddPlanning), userId, eid)
}

// AddVisited mocks base method.
func (m *MockUseCase) AddVisited(userId, eid uint64) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVisited", userId, eid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddVisited indicates an expected call of AddVisited.
func (mr *MockUseCaseMockRecorder) AddVisited(userId, eid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVisited", reflect.TypeOf((*MockUseCase)(nil).AddVisited), userId, eid)
}

// RemoveEvent mocks base method.
func (m *MockUseCase) RemoveEvent(userId, eid uint64) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveEvent", userId, eid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RemoveEvent indicates an expected call of RemoveEvent.
func (mr *MockUseCaseMockRecorder) RemoveEvent(userId, eid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEvent", reflect.TypeOf((*MockUseCase)(nil).RemoveEvent), userId, eid)
}

// SubscribeUser mocks base method.
func (m *MockUseCase) SubscribeUser(subscriberId, subscribedToId uint64) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeUser", subscriberId, subscribedToId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SubscribeUser indicates an expected call of SubscribeUser.
func (mr *MockUseCaseMockRecorder) SubscribeUser(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeUser", reflect.TypeOf((*MockUseCase)(nil).SubscribeUser), subscriberId, subscribedToId)
}

// UnsubscribeUser mocks base method.
func (m *MockUseCase) UnsubscribeUser(subscriberId, subscribedToId uint64) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsubscribeUser", subscriberId, subscribedToId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UnsubscribeUser indicates an expected call of UnsubscribeUser.
func (mr *MockUseCaseMockRecorder) UnsubscribeUser(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeUser", reflect.TypeOf((*MockUseCase)(nil).UnsubscribeUser), subscriberId, subscribedToId)
}