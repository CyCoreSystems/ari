package native

import "github.com/CyCoreSystems/ari"

type nativeDeviceState struct {
	client *Client
}

func (ds *nativeDeviceState) Get(name string) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(name, ds)
}

func (ds *nativeDeviceState) List() (dx []*ari.DeviceStateHandle, err error) {

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

func (ds *nativeDeviceState) Data(name string) (d ari.DeviceStateData, err error) {
	device := struct {
		State string `json:"state"`
	}{}
	err = ds.client.conn.Get("/deviceStates/"+name, &device)
	d = ari.DeviceStateData(device.State) //TODO: we can make DeviceStateData implement MarshalJSON/UnmarshalJSON
	return
}

func (ds *nativeDeviceState) Update(name string, state string) (err error) {
	req := map[string]string{
		"deviceState": state,
	}
	err = ds.client.conn.Put("/deviceStates/"+name, nil, &req)
	return
}

func (ds *nativeDeviceState) Delete(name string) (err error) {
	err = ds.client.conn.Delete("/deviceStates/"+name, nil, "")
	return
}
