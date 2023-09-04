// Code generated by mockery v1.0.0. DO NOT EDIT.

package arimocks

import (
	ari "github.com/CyCoreSystems/ari/v6"
	mock "github.com/stretchr/testify/mock"
)

// Bridge is an autogenerated mock type for the Bridge type
type Bridge struct {
	mock.Mock
}

// AddChannel provides a mock function with given fields: key, channelID
func (_m *Bridge) AddChannel(key *ari.Key, channelID string) error {
	ret := _m.Called(key, channelID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, string) error); ok {
		r0 = rf(key, channelID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddChannelWithOptions provides a mock function with given fields: key, channelID, options
func (_m *Bridge) AddChannelWithOptions(key *ari.Key, channelID string, options *ari.BridgeAddChannelOptions) error {
	ret := _m.Called(key, channelID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, string, *ari.BridgeAddChannelOptions) error); ok {
		r0 = rf(key, channelID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: key, btype, name
func (_m *Bridge) Create(key *ari.Key, btype string, name string) (*ari.BridgeHandle, error) {
	ret := _m.Called(key, btype, name)

	var r0 *ari.BridgeHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, string) *ari.BridgeHandle); ok {
		r0 = rf(key, btype, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.BridgeHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, string) error); ok {
		r1 = rf(key, btype, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Data provides a mock function with given fields: key
func (_m *Bridge) Data(key *ari.Key) (*ari.BridgeData, error) {
	ret := _m.Called(key)

	var r0 *ari.BridgeData
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.BridgeData); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.BridgeData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: key
func (_m *Bridge) Delete(key *ari.Key) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *Bridge) Get(key *ari.Key) *ari.BridgeHandle {
	ret := _m.Called(key)

	var r0 *ari.BridgeHandle
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.BridgeHandle); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.BridgeHandle)
		}
	}

	return r0
}

// List provides a mock function with given fields: _a0
func (_m *Bridge) List(_a0 *ari.Key) ([]*ari.Key, error) {
	ret := _m.Called(_a0)

	var r0 []*ari.Key
	if rf, ok := ret.Get(0).(func(*ari.Key) []*ari.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ari.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MOH provides a mock function with given fields: key, moh
func (_m *Bridge) MOH(key *ari.Key, moh string) error {
	ret := _m.Called(key, moh)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, string) error); ok {
		r0 = rf(key, moh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Play provides a mock function with given fields: key, playbackID, mediaURI
func (_m *Bridge) Play(key *ari.Key, playbackID string, mediaURI ...string) (*ari.PlaybackHandle, error) {
	_va := make([]interface{}, len(mediaURI))
	for _i := range mediaURI {
		_va[_i] = mediaURI[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key, playbackID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ari.PlaybackHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, ...string) *ari.PlaybackHandle); ok {
		r0 = rf(key, playbackID, mediaURI...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.PlaybackHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, ...string) error); ok {
		r1 = rf(key, playbackID, mediaURI...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Record provides a mock function with given fields: key, name, opts
func (_m *Bridge) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	ret := _m.Called(key, name, opts)

	var r0 *ari.LiveRecordingHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, *ari.RecordingOptions) *ari.LiveRecordingHandle); ok {
		r0 = rf(key, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.LiveRecordingHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, *ari.RecordingOptions) error); ok {
		r1 = rf(key, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveChannel provides a mock function with given fields: key, channelID
func (_m *Bridge) RemoveChannel(key *ari.Key, channelID string) error {
	ret := _m.Called(key, channelID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, string) error); ok {
		r0 = rf(key, channelID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StageCreate provides a mock function with given fields: key, btype, name
func (_m *Bridge) StageCreate(key *ari.Key, btype string, name string) (*ari.BridgeHandle, error) {
	ret := _m.Called(key, btype, name)

	var r0 *ari.BridgeHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, string) *ari.BridgeHandle); ok {
		r0 = rf(key, btype, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.BridgeHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, string) error); ok {
		r1 = rf(key, btype, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StagePlay provides a mock function with given fields: key, playbackID, mediaURI
func (_m *Bridge) StagePlay(key *ari.Key, playbackID string, mediaURI ...string) (*ari.PlaybackHandle, error) {
	_va := make([]interface{}, len(mediaURI))
	for _i := range mediaURI {
		_va[_i] = mediaURI[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key, playbackID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ari.PlaybackHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, ...string) *ari.PlaybackHandle); ok {
		r0 = rf(key, playbackID, mediaURI...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.PlaybackHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, ...string) error); ok {
		r1 = rf(key, playbackID, mediaURI...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StageRecord provides a mock function with given fields: key, name, opts
func (_m *Bridge) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	ret := _m.Called(key, name, opts)

	var r0 *ari.LiveRecordingHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string, *ari.RecordingOptions) *ari.LiveRecordingHandle); ok {
		r0 = rf(key, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.LiveRecordingHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string, *ari.RecordingOptions) error); ok {
		r1 = rf(key, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopMOH provides a mock function with given fields: key
func (_m *Bridge) StopMOH(key *ari.Key) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: key, n
func (_m *Bridge) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	_va := make([]interface{}, len(n))
	for _i := range n {
		_va[_i] = n[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 ari.Subscription
	if rf, ok := ret.Get(0).(func(*ari.Key, ...string) ari.Subscription); ok {
		r0 = rf(key, n...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ari.Subscription)
		}
	}

	return r0
}

// VideoSource provides a mock function with given fields: key, channelID
func (_m *Bridge) VideoSource(key *ari.Key, channelID string) error {
	ret := _m.Called(key, channelID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, string) error); ok {
		r0 = rf(key, channelID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VideoSourceDelete provides a mock function with given fields: key
func (_m *Bridge) VideoSourceDelete(key *ari.Key) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
