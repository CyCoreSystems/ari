// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/CyCoreSystems/ari (interfaces: Bus)

package mock

import (
	ari "github.com/CyCoreSystems/ari"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Bus interface
type MockBus struct {
	ctrl     *gomock.Controller
	recorder *_MockBusRecorder
}

// Recorder for MockBus (not exported)
type _MockBusRecorder struct {
	mock *MockBus
}

func NewMockBus(ctrl *gomock.Controller) *MockBus {
	mock := &MockBus{ctrl: ctrl}
	mock.recorder = &_MockBusRecorder{mock}
	return mock
}

func (_m *MockBus) EXPECT() *_MockBusRecorder {
	return _m.recorder
}

func (_m *MockBus) Send(_param0 *ari.Message) {
	_m.ctrl.Call(_m, "Send", _param0)
}

func (_mr *_MockBusRecorder) Send(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Send", arg0)
}

func (_m *MockBus) Subscribe(_param0 ...string) ari.Subscription {
	_s := []interface{}{}
	for _, _x := range _param0 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Subscribe", _s...)
	ret0, _ := ret[0].(ari.Subscription)
	return ret0
}

func (_mr *_MockBusRecorder) Subscribe(arg0 ...interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Subscribe", arg0...)
}
