package generic

import "github.com/CyCoreSystems/ari"

type DeviceState struct {
	Conn Conn
}

func (ds *DeviceState) Get(name string) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(name, ds)
}

func (ds *DeviceState) List() (dx []*ari.DeviceStateHandle, err error) {

	type device struct {
		Name string `json:"name"`
	}

	var devices []device
	err = ds.Conn.Get("/deviceStates", nil, &devices)
	for _, i := range devices {
		dx = append(dx, ds.Get(i.Name))
	}

	return
}

func (ds *DeviceState) Data(name string) (d ari.DeviceStateData, err error) {
	device := struct {
		State string `json:"state"`
	}{}
	err = ds.Conn.Get("/deviceStates/%s", []interface{}{name}, &device)
	d = ari.DeviceStateData(device.State) //TODO: we can make DeviceStateData implement MarshalJSON/UnmarshalJSON
	return
}

func (ds *DeviceState) Update(name string, state string) (err error) {
	req := map[string]string{
		"deviceState": state,
	}
	err = ds.Conn.Put("/deviceStates/%s", []interface{}{name}, nil, &req)
	return
}

func (ds *DeviceState) Delete(name string) (err error) {
	err = ds.Conn.Delete("/deviceStates/%s", []interface{}{name}, nil, "")
	return
}
