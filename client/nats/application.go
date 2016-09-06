package nats

import (
	"errors"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsApplication struct {
	conn *nats.EncodedConn
}

// Get returns a managed handle to an ARI application
func (a *natsApplication) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

// List returns the list of applications managed by asterisk
func (a *natsApplication) List() (ax []*ari.ApplicationHandle, err error) {
	var apps []string
	err = a.conn.Request("ari.applications", "", &apps, DefaultRequestTimeout)
	for _, app := range apps {
		ax = append(ax, a.Get(app))
	}
	return
}

// Data returns the details of a given ARI application
// Equivalent to GET /applications/{applicationName}
func (a *natsApplication) Data(name string) (d ari.ApplicationData, err error) {
	err = a.conn.Request("ari.applications.data."+name, "", &d, DefaultRequestTimeout)
	return
}

// Subscribe subscribes the given application to an event source
// Equivalent to POST /applications/{applicationName}/subscription
func (a *natsApplication) Subscribe(name string, eventSource string) (err error) {
	var response string
	err = a.conn.Request("ari.applications.subscribe."+name, eventSource, &response, DefaultRequestTimeout)
	if err == nil && response != "OK" {
		err = errors.New(response)
	}
	return err
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (a *natsApplication) Unsubscribe(name string, eventSource string) (err error) {
	var response string
	err = a.conn.Request("ari.applications.unsubscribe."+name, eventSource, &response, DefaultRequestTimeout)
	if err == nil && response != "OK" {
		err = errors.New(response)
	}
	return
}
