// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository.go

// Package mock_session is a generated GoMock package.
package mock_session

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
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

// CheckSession mocks base method.
func (m *MockRepository) CheckSession(value string) (bool, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", value)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockRepositoryMockRecorder) CheckSession(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockRepository)(nil).CheckSession), value)
}

// DeleteSession mocks base method.
func (m *MockRepository) DeleteSession(value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", value)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockRepositoryMockRecorder) DeleteSession(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockRepository)(nil).DeleteSession), value)
}

// InsertSession mocks base method.
func (m *MockRepository) InsertSession(userId uint64, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertSession", userId, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertSession indicates an expected call of InsertSession.
func (mr *MockRepositoryMockRecorder) InsertSession(userId, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertSession", reflect.TypeOf((*MockRepository)(nil).InsertSession), userId, value)
}