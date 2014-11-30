package ari

// Application describes a Stasis (Ari) application
type Application struct {
	Bridge_ids   []string `json:"bridge_ids"`   // Subscribed BridgeIds
	Channel_ids  []string `json:"channel_ids"`  // Subscribed ChannelIds
	Device_names []string `json:"device_names"` // Subscribed Device names
	Endpoint_ids []string `json:"endpoint_ids"` // Subscribed Endpoints (tech/resource format)
	Name         string   `json:"name"`         // Name of the application
}

// ListApplications returns the list of ARI applications on
// the Asterisk server
// Equivalent to GET /applications
func (c *Client) ListApplications() ([]Application, error) {
	var m []Application
	err := c.AriGet("/applications", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// GetApplication returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (c *Client) GetApplication(applicationName string) (Application, error) {
	var m Application
	err := c.AriGet("/applications/"+applicationName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// SubscribeApplication subscribes the given application to an event source
// event source may be one of:
//  - channel:<channelId>
//  - bridge:<bridgeId>
//  - endpoint:<tech>/<resource> (e.g. SIP/102)
//  - deviceState:<deviceName>
// Equivalent to POST /applications/{applicationName}/subscription
func (c *Client) SubscribeApplication(applicationName string, eventSource string) (Application, error) {
	var m Application

	type request struct {
		EventSource string `json:"eventSource"`
	}

	req := request{EventSource: eventSource}

	// Make the request
	err := c.AriPost("/applications/"+applicationName+"/subscription", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

// UnsubscribeApplication unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (c *Client) UnsubscribeApplication(applicationName string, eventSource string) (Application, error) {
	var m Application

	type request struct {
		EventSource string `json:"eventSource"`
	}

	req := request{eventSource}

	// TODO: handle Error Responses individually

	// Make the request
	err := c.AriDelete("/applications/"+applicationName+"/subscription", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Request structure for subscribing or unsubscribing to/from an application. EventSource is required.
