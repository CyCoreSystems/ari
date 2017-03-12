package native

import "github.com/CyCoreSystems/ari"

// Endpoint provides the ARI Endpoint accessors for the native client
type Endpoint struct {
	client *Client
}

// Get gets a lazy handle for the endpoint entity
func (e *Endpoint) Get(tech string, resource string) *ari.EndpointHandle {
	return ari.NewEndpointHandle(tech, resource, e)
}

// List lists the current endpoints and returns a list of handles
func (e *Endpoint) List() (ex []*ari.EndpointHandle, err error) {
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
func (e *Endpoint) ListByTech(tech string) (ex []*ari.EndpointHandle, err error) {
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
