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
func (a *Application) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

// List returns the list of applications managed by asterisk
func (a *Application) List() (ax []*ari.ApplicationHandle, err error) {
	var apps = []struct {
		Name string `json:"name"`
	}{}

	err = a.client.conn.Get("/applications", &apps)

	for _, i := range apps {
		ax = append(ax, a.Get(i.Name))
	}

	err = errors.Wrap(err, "Error listing applications")
	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *Application) Data(name string) (d *ari.ApplicationData, err error) {
	d = &ari.ApplicationData{}
	err = a.client.conn.Get("/applications/"+name, d)
	if err != nil {
		d = nil
		err = dataGetError(err, "application", "%v", name)
	}
	return
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *Application) Subscribe(name string, eventSource string) (err error) {
	req := struct {
		EventSource string `json:"eventSource"`
	}{
		EventSource: eventSource,
	}
	err = a.client.conn.Post("/applications/"+name+"/subscription", nil, &req)
	err = errors.Wrapf(err, "Error subscribing application '%v' for event source '%v'", name, eventSource)
	return
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *Application) Unsubscribe(name string, eventSource string) (err error) {
	// TODO: handle Error Responses individually
	err = a.client.conn.Delete("/applications/"+name+"/subscription", nil, fmt.Sprintf("eventSource=%s", eventSource))
	err = errors.Wrapf(err, "Error unsubscribing application '%v' for event source '%v'", name, eventSource)
	return
}
