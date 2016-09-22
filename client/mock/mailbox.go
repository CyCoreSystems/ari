// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/CyCoreSystems/ari (interfaces: Mailbox)

package mock

import (
	ari "github.com/CyCoreSystems/ari"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Mailbox interface
type MockMailbox struct {
	ctrl     *gomock.Controller
	recorder *_MockMailboxRecorder
}

// Recorder for MockMailbox (not exported)
type _MockMailboxRecorder struct {
	mock *MockMailbox
}

func NewMockMailbox(ctrl *gomock.Controller) *MockMailbox {
	mock := &MockMailbox{ctrl: ctrl}
	mock.recorder = &_MockMailboxRecorder{mock}
	return mock
}

func (_m *MockMailbox) EXPECT() *_MockMailboxRecorder {
	return _m.recorder
}

func (_m *MockMailbox) Data(_param0 string) (ari.MailboxData, error) {
	ret := _m.ctrl.Call(_m, "Data", _param0)
	ret0, _ := ret[0].(ari.MailboxData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockMailboxRecorder) Data(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Data", arg0)
}

func (_m *MockMailbox) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMailboxRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockMailbox) Get(_param0 string) *ari.MailboxHandle {
	ret := _m.ctrl.Call(_m, "Get", _param0)
	ret0, _ := ret[0].(*ari.MailboxHandle)
	return ret0
}

func (_mr *_MockMailboxRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockMailbox) List() ([]*ari.MailboxHandle, error) {
	ret := _m.ctrl.Call(_m, "List")
	ret0, _ := ret[0].([]*ari.MailboxHandle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockMailboxRecorder) List() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "List")
}

func (_m *MockMailbox) Update(_param0 string, _param1 int, _param2 int) error {
	ret := _m.ctrl.Call(_m, "Update", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMailboxRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1, arg2)
}
