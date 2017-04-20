package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

// Endpoint provides the ARI Endpoint accessors for the native client
type Endpoint struct {
	client *Client
}

// Get gets a lazy handle for the endpoint entity
func (e *Endpoint) Get(key *ari.Key) *ari.EndpointHandle {
	return ari.NewEndpointHandle(key, e)
}

// List lists the current endpoints and returns a list of handles
func (e *Endpoint) List(filter *ari.Key) (ex []*ari.Key, err error) {
	endpoints := []struct {
		Tech     string `json:"technology"`
		Resource string `json:"resource"`
	}{}
	if filter == nil {
		filter = ari.NodeKey(e.client.ApplicationName(), e.client.node)
	}
	err = e.client.get("/endpoints", &endpoints)
	for _, i := range endpoints {
		k := ari.NewEndpointKey(i.Tech, i.Resource, ari.WithApp(e.client.ApplicationName()), ari.WithNode(e.client.node))
		if filter.Match(k) {
			ex = append(ex, k)
		}
	}

	return
}

// ListByTech lists the current endpoints with the given technology and
// returns a list of handles.
func (e *Endpoint) ListByTech(tech string, filter *ari.Key) (ex []*ari.Key, err error) {
	endpoints := []struct {
		Tech     string `json:"technology"`
		Resource string `json:"resource"`
	}{}
	if filter == nil {
		filter = ari.NodeKey(e.client.ApplicationName(), e.client.node)
	}
	err = e.client.get("/endpoints/"+tech, &endpoints)
	for _, i := range endpoints {
		k := ari.NewEndpointKey(i.Tech, i.Resource, ari.WithApp(e.client.ApplicationName()), ari.WithNode(e.client.node))
		if filter.Match(k) {
			ex = append(ex, k)
		}
	}

	return
}

// Data retrieves the current state of the endpoint
func (e *Endpoint) Data(key *ari.Key) (*ari.EndpointData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("endpoint key not supplied")
	}
	if key.Kind != ari.EndpointKey {
		return nil, errors.New("wrong key type")
	}

	var data = new(ari.EndpointData)
	if err := e.client.get("/endpoints/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "endpoint", "%s", key.ID)
	}

	data.Key = e.client.stamp(key)
	return data, nil
}
