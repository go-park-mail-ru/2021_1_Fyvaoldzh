// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository.go

// Package mock_subscription is a generated GoMock package.
package mock_subscription

import (
	reflect "reflect"
	time "time"

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

// AddCountNotification mocks base method.
func (m *MockRepository) AddCountNotification(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCountNotification", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCountNotification indicates an expected call of AddCountNotification.
func (mr *MockRepositoryMockRecorder) AddCountNotification(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCountNotification", reflect.TypeOf((*MockRepository)(nil).AddCountNotification), id)
}

// AddPlanning mocks base method.
func (m *MockRepository) AddPlanning(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlanning", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPlanning indicates an expected call of AddPlanning.
func (mr *MockRepositoryMockRecorder) AddPlanning(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlanning", reflect.TypeOf((*MockRepository)(nil).AddPlanning), userId, eventId)
}

// AddPlanningNotification mocks base method.
func (m *MockRepository) AddPlanningNotification(eventId, userId uint64, eventDate, now time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlanningNotification", eventId, userId, eventDate, now)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPlanningNotification indicates an expected call of AddPlanningNotification.
func (mr *MockRepositoryMockRecorder) AddPlanningNotification(eventId, userId, eventDate, now interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlanningNotification", reflect.TypeOf((*MockRepository)(nil).AddPlanningNotification), eventId, userId, eventDate, now)
}

// AddSubscriptionAction mocks base method.
func (m *MockRepository) AddSubscriptionAction(subscriberId, subscribedToId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubscriptionAction", subscriberId, subscribedToId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSubscriptionAction indicates an expected call of AddSubscriptionAction.
func (mr *MockRepositoryMockRecorder) AddSubscriptionAction(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubscriptionAction", reflect.TypeOf((*MockRepository)(nil).AddSubscriptionAction), subscriberId, subscribedToId)
}

// AddUserEventAction mocks base method.
func (m *MockRepository) AddUserEventAction(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserEventAction", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserEventAction indicates an expected call of AddUserEventAction.
func (mr *MockRepositoryMockRecorder) AddUserEventAction(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserEventAction", reflect.TypeOf((*MockRepository)(nil).AddUserEventAction), userId, eventId)
}

// AddVisited mocks base method.
func (m *MockRepository) AddVisited(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVisited", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddVisited indicates an expected call of AddVisited.
func (mr *MockRepositoryMockRecorder) AddVisited(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVisited", reflect.TypeOf((*MockRepository)(nil).AddVisited), userId, eventId)
}

// CheckEventAdded mocks base method.
func (m *MockRepository) CheckEventAdded(userId, eventId uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEventAdded", userId, eventId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEventAdded indicates an expected call of CheckEventAdded.
func (mr *MockRepositoryMockRecorder) CheckEventAdded(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEventAdded", reflect.TypeOf((*MockRepository)(nil).CheckEventAdded), userId, eventId)
}

// CheckEventInList mocks base method.
func (m *MockRepository) CheckEventInList(eventId uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEventInList", eventId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEventInList indicates an expected call of CheckEventInList.
func (mr *MockRepositoryMockRecorder) CheckEventInList(eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEventInList", reflect.TypeOf((*MockRepository)(nil).CheckEventInList), eventId)
}

// CheckSubscription mocks base method.
func (m *MockRepository) CheckSubscription(subscriberId, subscribedToId uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSubscription", subscriberId, subscribedToId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSubscription indicates an expected call of CheckSubscription.
func (mr *MockRepositoryMockRecorder) CheckSubscription(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSubscription", reflect.TypeOf((*MockRepository)(nil).CheckSubscription), subscriberId, subscribedToId)
}

// GetTimeEvent mocks base method.
func (m *MockRepository) GetTimeEvent(eventId uint64) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimeEvent", eventId)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTimeEvent indicates an expected call of GetTimeEvent.
func (mr *MockRepositoryMockRecorder) GetTimeEvent(eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimeEvent", reflect.TypeOf((*MockRepository)(nil).GetTimeEvent), eventId)
}

// RemoveEvent mocks base method.
func (m *MockRepository) RemoveEvent(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveEvent", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveEvent indicates an expected call of RemoveEvent.
func (mr *MockRepositoryMockRecorder) RemoveEvent(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEvent", reflect.TypeOf((*MockRepository)(nil).RemoveEvent), userId, eventId)
}

// RemovePlanningNotification mocks base method.
func (m *MockRepository) RemovePlanningNotification(eventId, userId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePlanningNotification", eventId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePlanningNotification indicates an expected call of RemovePlanningNotification.
func (mr *MockRepositoryMockRecorder) RemovePlanningNotification(eventId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePlanningNotification", reflect.TypeOf((*MockRepository)(nil).RemovePlanningNotification), eventId, userId)
}

// RemoveSubscriptionAction mocks base method.
func (m *MockRepository) RemoveSubscriptionAction(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSubscriptionAction", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSubscriptionAction indicates an expected call of RemoveSubscriptionAction.
func (mr *MockRepositoryMockRecorder) RemoveSubscriptionAction(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSubscriptionAction", reflect.TypeOf((*MockRepository)(nil).RemoveSubscriptionAction), userId, eventId)
}

// RemoveUserEventAction mocks base method.
func (m *MockRepository) RemoveUserEventAction(userId, eventId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserEventAction", userId, eventId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserEventAction indicates an expected call of RemoveUserEventAction.
func (mr *MockRepositoryMockRecorder) RemoveUserEventAction(userId, eventId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserEventAction", reflect.TypeOf((*MockRepository)(nil).RemoveUserEventAction), userId, eventId)
}

// SubscribeUser mocks base method.
func (m *MockRepository) SubscribeUser(subscriberId, subscribedToId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeUser", subscriberId, subscribedToId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeUser indicates an expected call of SubscribeUser.
func (mr *MockRepositoryMockRecorder) SubscribeUser(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeUser", reflect.TypeOf((*MockRepository)(nil).SubscribeUser), subscriberId, subscribedToId)
}

// UnsubscribeUser mocks base method.
func (m *MockRepository) UnsubscribeUser(subscriberId, subscribedToId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsubscribeUser", subscriberId, subscribedToId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsubscribeUser indicates an expected call of UnsubscribeUser.
func (mr *MockRepositoryMockRecorder) UnsubscribeUser(subscriberId, subscribedToId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeUser", reflect.TypeOf((*MockRepository)(nil).UnsubscribeUser), subscriberId, subscribedToId)
}
