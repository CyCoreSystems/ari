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
func (a *Application) Get(key *ari.Key) ari.ApplicationHandle {
	return NewApplicationHandle(key, a)
}

// List returns the list of applications managed by asterisk
func (a *Application) List() (ax []*ari.Key, err error) {
	var apps = []struct {
		Name string `json:"name"`
	}{}

	err = a.client.get("/applications", &apps)

	for _, i := range apps {
		ax = append(ax, ari.NewKey(ari.ApplicationKey, i.Name))
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

// ApplicationHandle provides a wrapper to an Application interface for
// operations on a specific application
type ApplicationHandle struct {
	key *ari.Key
	a   *Application
}

// NewApplicationHandle creates a new handle to the application name
func NewApplicationHandle(key *ari.Key, app *Application) ari.ApplicationHandle {
	return &ApplicationHandle{
		key: key,
		a:   app,
	}
}

// ID returns the identifier for the application
func (ah *ApplicationHandle) ID() string {
	return ah.key.ID
}

// Data retrives the data for the application
func (ah *ApplicationHandle) Data() (ad *ari.ApplicationData, err error) {
	ad, err = ah.a.Data(ah.key)
	return
}

// Subscribe subscribes the application to an event source
// event source may be one of:
//  - channel:<channelId>
//  - bridge:<bridgeId>
//  - endpoint:<tech>/<resource> (e.g. SIP/102)
//  - deviceState:<deviceName>
func (ah *ApplicationHandle) Subscribe(eventSource string) (err error) {
	err = ah.a.Subscribe(ah.key, eventSource)
	return
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (ah *ApplicationHandle) Unsubscribe(eventSource string) (err error) {
	err = ah.a.Unsubscribe(ah.key, eventSource)
	return
}

// Match returns true fo the event matches the application
func (ah *ApplicationHandle) Match(e ari.Event) bool {
	return e.GetApplication() == ah.key.ID
}
