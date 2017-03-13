package native

import "github.com/CyCoreSystems/ari"

// DeviceState provides the ARI DeviceState accessors for the native client
type DeviceState struct {
	client *Client
}

// Get returns the lazy handle for the given device name
func (ds *DeviceState) Get(name string) ari.DeviceStateHandle {
	return NewDeviceStateHandle(name, ds)
}

// List lists the current devices and returns a list of handles
func (ds *DeviceState) List() (dx []ari.DeviceStateHandle, err error) {

	type device struct {
		Name string `json:"name"`
	}

	var devices []device
	err = ds.client.get("/deviceStates", &devices)
	for _, i := range devices {
		dx = append(dx, ds.Get(i.Name))
	}

	return
}

// Data retrieves the current state of the device
func (ds *DeviceState) Data(name string) (d *ari.DeviceStateData, err error) {
	device := struct {
		State string `json:"state"`
	}{}
	err = ds.client.get("/deviceStates/"+name, &device)
	if err != nil {
		d = nil
		err = dataGetError(err, "deviceState", "%v", name)
		return
	}
	x := ari.DeviceStateData(device.State) //TODO: we can make DeviceStateData implement MarshalJSON/UnmarshalJSON
	d = &x
	return
}

// Update updates the state of the device
func (ds *DeviceState) Update(name string, state string) (err error) {
	req := map[string]string{
		"deviceState": state,
	}
	err = ds.client.put("/deviceStates/"+name, nil, &req)
	return
}

// Delete deletes the device
func (ds *DeviceState) Delete(name string) (err error) {
	err = ds.client.del("/deviceStates/"+name, nil, "")
	return
}

// DeviceStateHandle is a representation of a device state
// that can be interacted with
type DeviceStateHandle struct {
	name string
	d    *DeviceState
}

// NewDeviceStateHandle creates a new deviceState handle
func NewDeviceStateHandle(name string, d *DeviceState) ari.DeviceStateHandle {
	return &DeviceStateHandle{
		name: name,
		d:    d,
	}
}

// ID returns the identifier for the device
func (dsh *DeviceStateHandle) ID() string {
	return dsh.name
}

// Data gets the device state
func (dsh *DeviceStateHandle) Data() (d *ari.DeviceStateData, err error) {
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
