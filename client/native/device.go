package native

import "github.com/CyCoreSystems/ari"

// DeviceState provides the ARI DeviceState accessors for the native client
type DeviceState struct {
	client *Client
}

// Get returns the lazy handle for the given device name
func (ds *DeviceState) Get(name string) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(name, ds)
}

// List lists the current devices and returns a list of handles
func (ds *DeviceState) List() (dx []*ari.DeviceStateHandle, err error) {

	type device struct {
		Name string `json:"name"`
	}

	var devices []device
	err = ds.client.conn.Get("/deviceStates", &devices)
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
	err = ds.client.conn.Get("/deviceStates/"+name, &device)
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
	err = ds.client.conn.Put("/deviceStates/"+name, nil, &req)
	return
}

// Delete deletes the device
func (ds *DeviceState) Delete(name string) (err error) {
	err = ds.client.conn.Delete("/deviceStates/"+name, nil, "")
	return
}
