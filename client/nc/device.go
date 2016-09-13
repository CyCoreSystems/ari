package nc

import (
	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsDeviceState struct {
	conn *nats.Conn
}

func (ds *natsDeviceState) Get(name string) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(name, ds)
}

func (ds *natsDeviceState) List() (dx []*ari.DeviceStateHandle, err error) {
	var devices []string
	err = request(ds.conn, "ari.devices.all", nil, &devices)
	for _, d := range devices {
		dx = append(dx, ds.Get(d))
	}
	return
}

func (ds *natsDeviceState) Data(name string) (d ari.DeviceStateData, err error) {
	err = request(ds.conn, "ari.devices.data."+name, nil, &d)
	return
}

func (ds *natsDeviceState) Update(name string, state string) (err error) {
	err = request(ds.conn, "ari.devices.update."+name, &state, nil)
	return
}

func (ds *natsDeviceState) Delete(name string) (err error) {
	err = request(ds.conn, "ari.devices.delete."+name, nil, nil)
	return
}
