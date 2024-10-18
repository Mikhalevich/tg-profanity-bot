// Code generated by MockGen. DO NOT EDIT.
// Source: internal/processor/port/msgsender.go

// Package port is a generated GoMock package.
package port

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMsgSender is a mock of MsgSender interface.
type MockMsgSender struct {
	ctrl     *gomock.Controller
	recorder *MockMsgSenderMockRecorder
}

// MockMsgSenderMockRecorder is the mock recorder for MockMsgSender.
type MockMsgSenderMockRecorder struct {
	mock *MockMsgSender
}

// NewMockMsgSender creates a new mock instance.
func NewMockMsgSender(ctrl *gomock.Controller) *MockMsgSender {
	mock := &MockMsgSender{ctrl: ctrl}
	mock.recorder = &MockMsgSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMsgSender) EXPECT() *MockMsgSenderMockRecorder {
	return m.recorder
}

// Edit mocks base method.
func (m *MockMsgSender) Edit(ctx context.Context, originMsgInfo MessageInfo, msg string, options ...Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, originMsgInfo, msg}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Edit", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit.
func (mr *MockMsgSenderMockRecorder) Edit(ctx, originMsgInfo, msg interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, originMsgInfo, msg}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockMsgSender)(nil).Edit), varargs...)
}

// Reply mocks base method.
func (m *MockMsgSender) Reply(ctx context.Context, originMsgInfo MessageInfo, msg string, options ...Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, originMsgInfo, msg}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Reply", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reply indicates an expected call of Reply.
func (mr *MockMsgSenderMockRecorder) Reply(ctx, originMsgInfo, msg interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, originMsgInfo, msg}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reply", reflect.TypeOf((*MockMsgSender)(nil).Reply), varargs...)
}