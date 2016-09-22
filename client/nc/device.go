package nc

import "github.com/CyCoreSystems/ari"

type natsDeviceState struct {
	conn *Conn
}

func (ds *natsDeviceState) Get(name string) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(name, ds)
}

func (ds *natsDeviceState) List() (dx []*ari.DeviceStateHandle, err error) {
	var devices []string
	err = ds.conn.readRequest("ari.devices.all", nil, &devices)
	for _, d := range devices {
		dx = append(dx, ds.Get(d))
	}
	return
}

func (ds *natsDeviceState) Data(name string) (d ari.DeviceStateData, err error) {
	err = ds.conn.readRequest("ari.devices.data."+name, nil, &d)
	return
}

func (ds *natsDeviceState) Update(name string, state string) (err error) {
	err = ds.conn.standardRequest("ari.devices.update."+name, &state, nil)
	return
}

func (ds *natsDeviceState) Delete(name string) (err error) {
	err = ds.conn.standardRequest("ari.devices.delete."+name, nil, nil)
	return
}
