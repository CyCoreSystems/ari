package ari

// DeviceState represents a communication path interacting with an
// Asterisk server for device state resources
type DeviceState interface {
	Get(key *Key) *DeviceStateHandle

	List(filter *Key) ([]*Key, error)

	Data(key *Key) (*DeviceStateData, error)

	Update(key *Key, state string) error

	Delete(key *Key) error
}

// DeviceStateData is the device state for the device
type DeviceStateData struct {
	// Key is the cluster-unique identifier for this device state
	Key *Key `json:"key"`

	State string `json:"state"`
}

// DeviceStateHandle is a representation of a device state
// that can be interacted with
type DeviceStateHandle struct {
	key *Key
	d   DeviceState
}

// NewDeviceStateHandle creates a new deviceState handle
func NewDeviceStateHandle(key *Key, d DeviceState) *DeviceStateHandle {
	return &DeviceStateHandle{
		key: key,
		d:   d,
	}
}

// ID returns the identifier for the device
func (dsh *DeviceStateHandle) ID() string {
	return dsh.key.ID
}

// Key returns the key for the device
func (dsh *DeviceStateHandle) Key() *Key {
	return dsh.key
}

// Data gets the device state
func (dsh *DeviceStateHandle) Data() (d *DeviceStateData, err error) {
	d, err = dsh.d.Data(dsh.key)
	return
}

// Update updates the device state, implicitly creating it if not exists
func (dsh *DeviceStateHandle) Update(state string) (err error) {
	err = dsh.d.Update(dsh.key, state)
	return
}

// Delete deletes the device state
func (dsh *DeviceStateHandle) Delete() (err error) {
	err = dsh.d.Delete(dsh.key)
	//NOTE: if err is not nil,
	// we could replace 'd' with a version of it
	// that always returns ErrNotFound. Not required, as the
	// handle could "come back" at any moment via an 'Update'
	return
}
