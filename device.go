package ari

// DeviceState represents a communication path interacting with an
// Asterisk server for device state resources
type DeviceState interface {
	Get(name string) *DeviceStateHandle

	List() ([]*DeviceStateHandle, error)

	Data(name string) (DeviceStateData, error)

	Update(name string, state string) error

	Delete(name string) error
}

// DeviceStateData is the device state for the device
type DeviceStateData string

// NewDeviceStateHandle creates a new deviceState handle
func NewDeviceStateHandle(name string, d DeviceState) *DeviceStateHandle {
	return &DeviceStateHandle{
		name: name,
		d:    d,
	}
}

// DeviceStateHandle is a representation of a device state
// that can be interacted with
type DeviceStateHandle struct {
	name string
	d    DeviceState
}

// ID returns the identifier for the device
func (dsh *DeviceStateHandle) ID() string {
	return dsh.name
}

// Data gets the device state
func (dsh *DeviceStateHandle) Data() (d DeviceStateData, err error) {
	d, err = dsh.d.Data(dsh.name)
	return
}

// Update updates the device state, implicitly creating it if not exists
func (dsh *DeviceStateHandle) Update(state string) (err error) {
	err = dsh.d.Update(dsh.name, state)
	return
}

// Delete deletes the device state
func (dsh *DeviceStateHandle) Delete() (err error) {
	err = dsh.d.Delete(dsh.name)
	//NOTE: if err is not nil,
	// we could replace 'd' with a version of it
	// that always returns ErrNotFound. Not required, as the
	// handle could "come back" at any moment via an 'Update'
	return
}
