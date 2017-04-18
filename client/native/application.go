package native

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
)

// Application is a native implementation of ARI's Application functions
type Application struct {
	client *Client
}

// Get returns a managed handle to an ARI application
func (a *Application) Get(key *ari.Key) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(key, a)
}

// List returns the list of applications managed by asterisk
func (a *Application) List(filter *ari.Key) (ax []*ari.Key, err error) {

	if filter == nil {
		filter = ari.NodeKey(a.client.ApplicationName(), a.client.node)
	}

	var apps = []struct {
		Name string `json:"name"`
	}{}

	err = a.client.get("/applications", &apps)

	for _, i := range apps {
		k := ari.NewKey(ari.ApplicationKey, i.Name, ari.WithApp(a.client.ApplicationName()), ari.WithNode(a.client.node))
		if filter.Match(k) {
			ax = append(ax, k)
		}
	}

	err = errors.Wrap(err, "Error listing applications")
	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *Application) Data(key *ari.Key) (d *ari.ApplicationData, err error) {
	d = &ari.ApplicationData{}
	name := key.ID
	err = a.client.get("/applications/"+name, d)
	if err != nil {
		d = nil
		err = dataGetError(err, "application", "%v", name)
	}
	return
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *Application) Subscribe(key *ari.Key, eventSource string) (err error) {
	req := struct {
		EventSource string `json:"eventSource"`
	}{
		EventSource: eventSource,
	}
	name := key.ID
	err = a.client.post("/applications/"+name+"/subscription", nil, &req)
	err = errors.Wrapf(err, "Error subscribing application '%v' for event source '%v'", name, eventSource)
	return
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *Application) Unsubscribe(key *ari.Key, eventSource string) (err error) {
	name := key.ID
	err = a.client.del("/applications/"+name+"/subscription", nil, fmt.Sprintf("eventSource=%s", eventSource))
	err = errors.Wrapf(err, "Error unsubscribing application '%v' for event source '%v'", name, eventSource)
	return
}
