package native

import (
	"errors"

	"github.com/Amtelco-Software/ari/v6"
)

// DeviceState provides the ARI DeviceState accessors for the native client
type DeviceState struct {
	client *Client
}

// Get returns the lazy handle for the given device name
func (ds *DeviceState) Get(key *ari.Key) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(ds.client.stamp(key), ds)
}

// List lists the current devices and returns a list of handles
func (ds *DeviceState) List(filter *ari.Key) (dx []*ari.Key, err error) {
	type device struct {
		Name string `json:"name"`
	}

	if filter == nil {
		filter = ds.client.stamp(ari.NewKey(ari.DeviceStateKey, ""))
	}

	var devices []device

	if err = ds.client.get("/deviceStates", &devices); err != nil {
		return nil, err
	}

	for _, i := range devices {
		k := ds.client.stamp(ari.NewKey(ari.DeviceStateKey, i.Name))
		if filter.Match(k) {
			dx = append(dx, k)
		}
	}

	return
}

// Data retrieves the current state of the device
func (ds *DeviceState) Data(key *ari.Key) (*ari.DeviceStateData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("device key not supplied")
	}

	data := new(ari.DeviceStateData)
	if err := ds.client.get("/deviceStates/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "deviceState", "%v", key.ID)
	}

	data.Key = ds.client.stamp(key)

	return data, nil
}

// Update updates the state of the device
func (ds *DeviceState) Update(key *ari.Key, state string) error {
	req := map[string]string{
		"deviceState": state,
	}

	return ds.client.put("/deviceStates/"+key.ID, nil, &req)
}

// Delete deletes the device
func (ds *DeviceState) Delete(key *ari.Key) error {
	return ds.client.del("/deviceStates/"+key.ID, nil, "")
}
