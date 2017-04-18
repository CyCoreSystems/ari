package ari

// DeviceState represents a communication path interacting with an
// Asterisk server for device state resources
type DeviceState interface {
	Get(key *Key) DeviceStateHandle

	List(filter *Key) ([]*Key, error)

	Data(key *Key) (*DeviceStateData, error)

	Update(key *Key, state string) error

	Delete(key *Key) error
}

// DeviceStateData is the device state for the device
type DeviceStateData string

// DeviceStateHandle is a representation of a device state
// that can be interacted with
type DeviceStateHandle interface {
	// ID returns the identifier for the device
	ID() string

	// Data gets the device state
	Data() (d *DeviceStateData, err error)

	// Update updates the device state, implicitly creating it if not exists
	Update(state string) (err error)

	// Delete deletes the device state
	Delete() (err error)
}
