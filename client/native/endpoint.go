package native

import "github.com/CyCoreSystems/ari"

// Endpoint provides the ARI Endpoint accessors for the native client
type Endpoint struct {
	client *Client
}

// Get gets a lazy handle for the endpoint entity
func (e *Endpoint) Get(tech string, resource string) ari.EndpointHandle {
	return NewEndpointHandle(tech, resource, e)
}

// List lists the current endpoints and returns a list of handles
func (e *Endpoint) List() (ex []ari.EndpointHandle, err error) {
	endpoints := []struct {
		Tech     string `json:"technology"`
		Resource string `json:"resource"`
	}{}
	err = e.client.get("/endpoints", &endpoints)
	for _, i := range endpoints {
		ex = append(ex, e.Get(i.Tech, i.Resource))
	}

	return
}

// ListByTech lists the current endpoints with the given technology and
// returns a list of handles.
func (e *Endpoint) ListByTech(tech string) (ex []ari.EndpointHandle, err error) {
	err = e.client.get("/endpoints/"+tech, &ex)
	return
}

// Data retrieves the current state of the endpoint
func (e *Endpoint) Data(tech string, resource string) (ed *ari.EndpointData, err error) {
	ed = &ari.EndpointData{}
	err = e.client.get("/endpoints/"+tech+"/"+resource, ed)
	if err != nil {
		ed = nil
		err = dataGetError(err, "endpoint", "%v/%v", tech, resource)
	}

	return
}

// NewEndpointHandle creates a new EndpointHandle
func NewEndpointHandle(tech string, resource string, e *Endpoint) ari.EndpointHandle {
	return &EndpointHandle{
		tech:     tech,
		resource: resource,
	}
}

// An EndpointHandle is a reference to an endpoint attached to
// a transport to an asterisk server
type EndpointHandle struct {
	tech     string
	resource string
	e        *Endpoint
}

// ID returns the identifier for the endpoint
func (eh *EndpointHandle) ID() string {
	return eh.tech + "/" + eh.resource
}

// Data returns the state of the endpoint
func (eh *EndpointHandle) Data() (*ari.EndpointData, error) {
	return eh.e.Data(eh.tech, eh.resource)
}

// Match returns true if the event matches the bridge
func (eh *EndpointHandle) Match(e ari.Event) bool {
	en, ok := e.(ari.EndpointEvent)
	if !ok {
		return false
	}
	ids := en.GetEndpointIDs()
	for _, i := range ids {
		if i == eh.ID() {
			return true
		}
	}
	return false
}
