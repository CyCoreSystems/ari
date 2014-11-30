package ari

// DeviceState describes the state of a device
type DeviceState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

//ListDeviceStates returns the list of all ARI controlled device states
//Equivalent to GET /deviceStates
func (c *Client) ListDeviceStates() ([]DeviceState, error) {
	var m []DeviceState
	err := c.AriGet("/deviceStates", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Retrieve the current state of a specified device
//Equivalent to GET /deviceStates/{deviceName}
func (c *Client) GetDeviceState(deviceName string) (DeviceState, error) {
	var m DeviceState
	err := c.AriGet("/deviceStates/"+deviceName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Change the state of a device controlled by ARI. (Note - implicitly creates the device state).
//Equivalent to PUT /deviceStates/{deviceName}
func (c *Client) ChangeDeviceState(deviceName string, deviceState string) error {
	type request struct {
		DeviceState string `json:"deviceState"`
	}

	req := request{deviceState}

	//send request
	err := c.AriPut("/deviceStates/"+deviceName, nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Destroy a device-state controlled by ARI.
//Equivalent to DELETE /deviceStates/{deviceName}
func (c *Client) DeleteDeviceState(deviceName string) error {
	err := c.AriDelete("/deviceStates/"+deviceName, nil, nil)
	return err
}
