package native

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
)

type nativeApplication struct {
	conn *Conn
}

// Get returns a managed handle to an ARI application
func (a *nativeApplication) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

// List returns the list of applications managed by asterisk
func (a *nativeApplication) List() (ax []*ari.ApplicationHandle, err error) {
	var apps = []struct {
		Name string `json:"name"`
	}{}

	err = Get(a.conn, "/applications", &apps)

	for _, i := range apps {
		ax = append(ax, a.Get(i.Name))
	}

	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *nativeApplication) Data(name string) (d ari.ApplicationData, err error) {
	err = Get(a.conn, "/applications/"+name, &d)
	return
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *nativeApplication) Subscribe(name string, eventSource string) (err error) {
	var m ari.ApplicationData

	type request struct {
		EventSource string `json:"eventSource"`
	}

	req := request{EventSource: eventSource}
	err = Post(a.conn, "/applications/"+name+"/subscription", &m, &req)
	return err
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *nativeApplication) Unsubscribe(name string, eventSource string) (err error) {
	var m ari.ApplicationData

	// TODO: handle Error Responses individually

	// Make the request
	err = Delete(a.conn, "/applications/"+name+"/subscription", &m, fmt.Sprintf("eventSource=%s", eventSource))
	return
}
