package native

import "github.com/CyCoreSystems/ari"

// DeviceState provides the ARI DeviceState accessors for the native client
type DeviceState struct {
	client *Client
}

// Get returns the lazy handle for the given device name
func (ds *DeviceState) Get(key *ari.Key) ari.DeviceStateHandle {
	return NewDeviceStateHandle(key, ds)
}

// List lists the current devices and returns a list of handles
func (ds *DeviceState) List(filter *ari.Key) (dx []*ari.Key, err error) {

	type device struct {
		Name string `json:"name"`
	}

	if filter == nil {
		filter = ari.NodeKey(ds.client.ApplicationName(), ds.client.node)
	}

	var devices []device
	err = ds.client.get("/deviceStates", &devices)
	for _, i := range devices {
		k := ari.NewKey(ari.DeviceStateKey, i.Name, ari.WithApp(ds.client.ApplicationName()), ari.WithNode(ds.client.node))
		if filter.Match(k) {
			dx = append(dx, k)
		}
	}

	return
}

// Data retrieves the current state of the device
func (ds *DeviceState) Data(key *ari.Key) (d *ari.DeviceStateData, err error) {
	device := struct {
		State string `json:"state"`
	}{}
	name := key.ID
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
func (ds *DeviceState) Update(key *ari.Key, state string) (err error) {
	req := map[string]string{
		"deviceState": state,
	}
	name := key.ID
	err = ds.client.put("/deviceStates/"+name, nil, &req)
	return
}

// Delete deletes the device
func (ds *DeviceState) Delete(key *ari.Key) (err error) {
	name := key.ID
	err = ds.client.del("/deviceStates/"+name, nil, "")
	return
}

// DeviceStateHandle is a representation of a device state
// that can be interacted with
type DeviceStateHandle struct {
	key *ari.Key
	d   *DeviceState
}

// NewDeviceStateHandle creates a new deviceState handle
func NewDeviceStateHandle(key *ari.Key, d *DeviceState) ari.DeviceStateHandle {
	return &DeviceStateHandle{
		key: key,
		d:   d,
	}
}

// ID returns the identifier for the device
func (dsh *DeviceStateHandle) ID() string {
	return dsh.key.ID
}

// Data gets the device state
func (dsh *DeviceStateHandle) Data() (d *ari.DeviceStateData, err error) {
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
