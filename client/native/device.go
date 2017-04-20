package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

// DeviceState provides the ARI DeviceState accessors for the native client
type DeviceState struct {
	client *Client
}

// Get returns the lazy handle for the given device name
func (ds *DeviceState) Get(key *ari.Key) *ari.DeviceStateHandle {
	return ari.NewDeviceStateHandle(key, ds)
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
func (ds *DeviceState) Data(key *ari.Key) (*ari.DeviceStateData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("device key not supplied")
	}

	var data = new(ari.DeviceStateData)
	if err := ds.client.get("/deviceStates/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "deviceState", "%v", key.ID)
	}

	data.Key = ds.client.stamp(key)
	return data, nil
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
