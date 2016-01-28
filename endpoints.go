package ari

// Endpoint describes an external device which may offer or accept calls
// to or from Asterisk.  Devices are defined by a technology/resource pair.
//
// Allowed states:  'unknown', 'offline', 'online'
type Endpoint struct {
	ChannelIds []string `json:"channel_ids"`     // List of channel Ids which are associated with this endpoint
	Resource   string   `json:"resource"`        // The endpoint's resource name
	State      string   `json:"state,omitempty"` // The state of the endpoint
	Technology string   `json:"technology"`      // The technology of the endpoint (e.g. SIP, PJSIP, DAHDI, etc)
}

// ListEndpoints lists all endpoints
// TODO: associated with the application, or on the entire system?
func (c *Client) ListEndpoints() ([]Endpoint, error) {
	var m []Endpoint
	err := c.Get("/endpoints", &m)
	return m, err
}

//List available endpoints for a given endpoint technology
//Equivalent to Get /endpoints/{tech}
func (c *Client) GetEndpointsByTech(tech string) ([]Endpoint, error) {
	var m []Endpoint
	err := c.Get("/endpoints/"+tech, &m)
	return m, err
}

// GetEndpoint requests the details of an endpoint from Asterisk
func (c *Client) GetEndpoint(tech string, resource string) (Endpoint, error) {
	var m Endpoint
	err := c.Get("/endpoints/"+tech+"/"+resource, &m)
	return m, err
}
