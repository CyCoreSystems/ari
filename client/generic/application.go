package generic

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
)

type Application struct {
	Conn Conn
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

	err = a.Conn.Get("/applications", nil, &apps)

	for _, i := range apps {
		ax = append(ax, a.Get(i.Name))
	}

	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *Application) Data(name string) (d ari.ApplicationData, err error) {
	err = a.Conn.Get("/applications/%s", []interface{}{name}, &d)
	return
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *Application) Subscribe(name string, eventSource string) (err error) {
	var m ari.ApplicationData

	type request struct {
		EventSource string `json:"eventSource"`
	}

	req := request{EventSource: eventSource}
	err = a.Conn.Post("/applications/%s/subscription", []interface{}{name}, &m, &req)
	return err
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *Application) Unsubscribe(name string, eventSource string) (err error) {
	var m ari.ApplicationData

	// TODO: handle Error Responses individually

	// Make the request
	err = a.Conn.Delete("/applications/%s/subscription", []interface{}{name}, &m, fmt.Sprintf("eventSource=%s", eventSource))
	return
}
