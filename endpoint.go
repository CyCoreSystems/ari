package ari

import (
	"errors"
	"strings"
)

// EndpointIDSeparator seperates the ID components of the endpoint ID
const EndpointIDSeparator = "|" //TODO: confirm separator isn't terrible

// Endpoint represents a communication path to an Asterisk server
// for endpoint resources
type Endpoint interface {

	// List lists the endpoints
	// TODO: associated with the application, or on the entire system?
	List() ([]*EndpointHandle, error)

	// List available endpoints for a given endpoint technology
	ListByTech(tech string) ([]*EndpointHandle, error)

	// Get returns a handle to the endpoint for further operations
	Get(tech string, resource string) *EndpointHandle

	// Data returns the state of the endpoint
	Data(tech string, resource string) (EndpointData, error)
}

// EndpointData describes an external device which may offer or accept calls
// to or from Asterisk.  Devices are defined by a technology/resource pair.
//
// Allowed states:  'unknown', 'offline', 'online'
type EndpointData struct {
	ChannelIDs []string `json:"channel_ids"`     // List of channel Ids which are associated with this endpoint
	Resource   string   `json:"resource"`        // The endpoint's resource name
	State      string   `json:"state,omitempty"` // The state of the endpoint
	Technology string   `json:"technology"`      // The technology of the endpoint (e.g. SIP, PJSIP, DAHDI, etc)
}

// ID returns the ID for the endpoint
func (ed *EndpointData) ID() string {
	return ed.Technology + "/" + ed.Resource
}

// NewEndpointHandle creates a new EndpointHandle
func NewEndpointHandle(tech string, resource string, e Endpoint) *EndpointHandle {
	return &EndpointHandle{
		tech:     tech,
		resource: resource,
	}
}

// An EndpointHandle is a reference to an endpoint attached to
// a transport to an asterisk server
type EndpointHandle struct {
	tech     string
	resource string
	e        Endpoint
}

// ID returns the identifier for the endpoint
func (eh *EndpointHandle) ID() string {
	return eh.tech + "/" + eh.resource
}

// Data returns the state of the endpoint
func (eh *EndpointHandle) Data() (EndpointData, error) {
	return eh.e.Data(eh.tech, eh.resource)
}

// FromEndpointID converts the endpoint ID to the tech, resource pair.
func FromEndpointID(id string) (tech string, resource string, err error) {
	items := strings.Split(id, EndpointIDSeparator)
	if len(items) < 2 {
		err = errors.New("Endpoint ID is not in tech" + EndpointIDSeparator + "resource format")
		return
	}

	if len(items) > 2 {
		// huge programmer error here, we want to handle it
		// tempted to panic here...
		err = errors.New("EndpointIDSeparator is conflicting with tech and resource identifiers")
		return
	}

	tech = items[0]
	resource = items[1]
	return
}

// Match returns true if the event matches the bridge
func (eh *EndpointHandle) Match(e Event) bool {
	en, ok := e.(EndpointEvent)
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
