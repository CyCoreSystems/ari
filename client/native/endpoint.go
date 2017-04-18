package native

import (
	"errors"
	"strings"

	"github.com/CyCoreSystems/ari"
)

// Endpoint provides the ARI Endpoint accessors for the native client
type Endpoint struct {
	client *Client
}

// Get gets a lazy handle for the endpoint entity
func (e *Endpoint) Get(key *ari.Key) ari.EndpointHandle {
	return NewEndpointHandle(key, e)
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
func (e *Endpoint) Data(key *ari.Key) (ed *ari.EndpointData, err error) {
	if key.Kind != ari.EndpointKey {
		err = errors.New("wrong key type")
		return
	}
	items := strings.Split(key.ID, "/")
	if len(items) != 2 {
		err = errors.New("malformed key")
		return
	}
	tech := items[0]
	resource := items[1]
	ed = &ari.EndpointData{}
	err = e.client.get("/endpoints/"+tech+"/"+resource, ed)
	if err != nil {
		ed = nil
		err = dataGetError(err, "endpoint", "%v/%v", tech, resource)
	}

	return
}

// NewEndpointHandle creates a new EndpointHandle
func NewEndpointHandle(key *ari.Key, e *Endpoint) ari.EndpointHandle {
	return &EndpointHandle{
		key: key,
		e:   e,
	}
}

// An EndpointHandle is a reference to an endpoint attached to
// a transport to an asterisk server
type EndpointHandle struct {
	key *ari.Key
	e   *Endpoint
}

// ID returns the identifier for the endpoint
func (eh *EndpointHandle) ID() string {
	return eh.key.ID
}

// Data returns the state of the endpoint
func (eh *EndpointHandle) Data() (*ari.EndpointData, error) {
	return eh.e.Data(eh.key)
}

// Match returns true if the event matches the endpoint
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
