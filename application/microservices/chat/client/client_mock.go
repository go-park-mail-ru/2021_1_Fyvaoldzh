// Code generated by MockGen. DO NOT EDIT.
// Source: client_interface.go

// Package mock_client is a generated GoMock package.
package client

import (
	models "kudago/application/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIChatClient is a mock of IChatClient interface.
type MockIChatClient struct {
	ctrl     *gomock.Controller
	recorder *MockIChatClientMockRecorder
}

// MockIChatClientMockRecorder is the mock recorder for MockIChatClient.
type MockIChatClientMockRecorder struct {
	mock *MockIChatClient
}

// NewMockIChatClient creates a new mock instance.
func NewMockIChatClient(ctrl *gomock.Controller) *MockIChatClient {
	mock := &MockIChatClient{ctrl: ctrl}
	mock.recorder = &MockIChatClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIChatClient) EXPECT() *MockIChatClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockIChatClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockIChatClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIChatClient)(nil).Close))
}

// DeleteDialogue mocks base method.
func (m *MockIChatClient) DeleteDialogue(uid, id uint64) (error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDialogue", uid, id)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// DeleteDialogue indicates an expected call of DeleteDialogue.
func (mr *MockIChatClientMockRecorder) DeleteDialogue(uid, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDialogue", reflect.TypeOf((*MockIChatClient)(nil).DeleteDialogue), uid, id)
}

// DeleteMessage mocks base method.
func (m *MockIChatClient) DeleteMessage(uid, id uint64) (error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", uid, id)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockIChatClientMockRecorder) DeleteMessage(uid, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockIChatClient)(nil).DeleteMessage), uid, id)
}

// EditMessage mocks base method.
func (m *MockIChatClient) EditMessage(uid uint64, newMessage *models.RedactMessage) (error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditMessage", uid, newMessage)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// EditMessage indicates an expected call of EditMessage.
func (mr *MockIChatClientMockRecorder) EditMessage(uid, newMessage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditMessage", reflect.TypeOf((*MockIChatClient)(nil).EditMessage), uid, newMessage)
}

// GetAllCounts mocks base method.
func (m *MockIChatClient) GetAllCounts(uid uint64) (models.Counts, error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCounts", uid)
	ret0, _ := ret[0].(models.Counts)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// GetAllCounts indicates an expected call of GetAllCounts.
func (mr *MockIChatClientMockRecorder) GetAllCounts(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCounts", reflect.TypeOf((*MockIChatClient)(nil).GetAllCounts), uid)
}

// GetAllDialogues mocks base method.
func (m *MockIChatClient) GetAllDialogues(uid uint64, page int) (models.DialogueCards, error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDialogues", uid, page)
	ret0, _ := ret[0].(models.DialogueCards)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// GetAllDialogues indicates an expected call of GetAllDialogues.
func (mr *MockIChatClientMockRecorder) GetAllDialogues(uid, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDialogues", reflect.TypeOf((*MockIChatClient)(nil).GetAllDialogues), uid, page)
}

// GetAllNotifications mocks base method.
func (m *MockIChatClient) GetAllNotifications(uid uint64, page int) (models.Notifications, error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotifications", uid, page)
	ret0, _ := ret[0].(models.Notifications)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// GetAllNotifications indicates an expected call of GetAllNotifications.
func (mr *MockIChatClientMockRecorder) GetAllNotifications(uid, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotifications", reflect.TypeOf((*MockIChatClient)(nil).GetAllNotifications), uid, page)
}

// GetOneDialogue mocks base method.
func (m *MockIChatClient) GetOneDialogue(uid1, uid2 uint64, page int) (models.Dialogue, error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneDialogue", uid1, uid2, page)
	ret0, _ := ret[0].(models.Dialogue)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// GetOneDialogue indicates an expected call of GetOneDialogue.
func (mr *MockIChatClientMockRecorder) GetOneDialogue(uid1, uid2, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneDialogue", reflect.TypeOf((*MockIChatClient)(nil).GetOneDialogue), uid1, uid2, page)
}

// Mailing mocks base method.
func (m *MockIChatClient) Mailing(uid uint64, mailing *models.Mailing) (error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mailing", uid, mailing)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// Mailing indicates an expected call of Mailing.
func (mr *MockIChatClientMockRecorder) Mailing(uid, mailing interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mailing", reflect.TypeOf((*MockIChatClient)(nil).Mailing), uid, mailing)
}

// Search mocks base method.
func (m *MockIChatClient) Search(uid uint64, id int, str string, page int) (models.DialogueCards, error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", uid, id, str, page)
	ret0, _ := ret[0].(models.DialogueCards)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(int)
	return ret0, ret1, ret2
}

// Search indicates an expected call of Search.
func (mr *MockIChatClientMockRecorder) Search(uid, id, str, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIChatClient)(nil).Search), uid, id, str, page)
}

// SendMessage mocks base method.
func (m *MockIChatClient) SendMessage(newMessage *models.NewMessage, uid uint64) (error, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", newMessage, uid)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockIChatClientMockRecorder) SendMessage(newMessage, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockIChatClient)(nil).SendMessage), newMessage, uid)
}
