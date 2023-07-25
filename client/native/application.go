package native

import (
	"fmt"

	"github.com/Amtelco-Software/ari/v6"
	"github.com/rotisserie/eris"
)

// Application is a native implementation of ARI's Application functions
type Application struct {
	client *Client
}

// Get returns a managed handle to an ARI application
func (a *Application) Get(key *ari.Key) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(a.client.stamp(key), a)
}

// List returns the list of applications managed by asterisk
func (a *Application) List(filter *ari.Key) (ax []*ari.Key, err error) {
	if filter == nil {
		filter = ari.NewKey(ari.ApplicationKey, "")
	}

	apps := []struct {
		Name string `json:"name"`
	}{}

	err = a.client.get("/applications", &apps)

	for _, i := range apps {
		k := a.client.stamp(ari.NewKey(ari.ApplicationKey, i.Name))
		if filter.Match(k) {
			ax = append(ax, k)
		}
	}

	err = eris.Wrap(err, "Error listing applications")

	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *Application) Data(key *ari.Key) (*ari.ApplicationData, error) {
	if key == nil || key.ID == "" {
		return nil, eris.New("application key not supplied")
	}

	data := new(ari.ApplicationData)
	if err := a.client.get("/applications/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "application", "%v", key.ID)
	}

	data.Key = a.client.stamp(key)

	return data, nil
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *Application) Subscribe(key *ari.Key, eventSource string) error {
	req := struct {
		EventSource string `json:"eventSource"`
	}{
		EventSource: eventSource,
	}

	err := a.client.post("/applications/"+key.ID+"/subscription", nil, &req)

	return eris.Wrapf(err, "Error subscribing application '%v' for event source '%v'", key.ID, eventSource)
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *Application) Unsubscribe(key *ari.Key, eventSource string) error {
	name := key.ID
	err := a.client.del("/applications/"+name+"/subscription", nil, fmt.Sprintf("eventSource=%s", eventSource))

	return eris.Wrapf(err, "Error unsubscribing application '%v' for event source '%v'", name, eventSource)
}
