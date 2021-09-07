// Code generated by mockery v1.0.0. DO NOT EDIT.

package arimocks

import (
	ari "github.com/CyCoreSystems/ari/v5"
	mock "github.com/stretchr/testify/mock"
)

// StoredRecording is an autogenerated mock type for the StoredRecording type
type StoredRecording struct {
	mock.Mock
}

// Copy provides a mock function with given fields: key, dest
func (_m *StoredRecording) Copy(key *ari.Key, dest string) (*ari.StoredRecordingHandle, error) {
	ret := _m.Called(key, dest)

	var r0 *ari.StoredRecordingHandle
	if rf, ok := ret.Get(0).(func(*ari.Key, string) *ari.StoredRecordingHandle); ok {
		r0 = rf(key, dest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.StoredRecordingHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key, string) error); ok {
		r1 = rf(key, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Data provides a mock function with given fields: key
func (_m *StoredRecording) Data(key *ari.Key) (*ari.StoredRecordingData, error) {
	ret := _m.Called(key)

	var r0 *ari.StoredRecordingData
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.StoredRecordingData); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.StoredRecordingData)
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

// Data provides a mock function with given fields: key
func (_m *StoredRecording) File(key *ari.Key) (*ari.StoredRecordingFile, error) {
	ret := _m.Called(key)

	var r0 *ari.StoredRecordingFile
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.StoredRecordingFile); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.StoredRecordingFile)
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
func (_m *StoredRecording) Delete(key *ari.Key) error {
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
func (_m *StoredRecording) Get(key *ari.Key) *ari.StoredRecordingHandle {
	ret := _m.Called(key)

	var r0 *ari.StoredRecordingHandle
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.StoredRecordingHandle); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.StoredRecordingHandle)
		}
	}

	return r0
}

// List provides a mock function with given fields: filter
func (_m *StoredRecording) List(filter *ari.Key) ([]*ari.Key, error) {
	ret := _m.Called(filter)

	var r0 []*ari.Key
	if rf, ok := ret.Get(0).(func(*ari.Key) []*ari.Key); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ari.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
