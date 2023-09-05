package ari

import (
	"errors"
	"strings"
)

// EndpointIDSeparator seperates the ID components of the endpoint ID
const EndpointIDSeparator = "|" // TODO: confirm separator isn't terrible

// Endpoint represents a communication path to an Asterisk server
// for endpoint resources
type Endpoint interface {
	// List lists the endpoints
	List(filter *Key) ([]*Key, error)

	// List available endpoints for a given endpoint technology
	ListByTech(tech string, filter *Key) ([]*Key, error)

	// Get returns a handle to the endpoint for further operations
	Get(key *Key) *EndpointHandle

	// Data returns the state of the endpoint
	Data(key *Key) (*EndpointData, error)
}

// NewEndpointKey returns the key for the given endpoint
func NewEndpointKey(tech, resource string, opts ...KeyOptionFunc) *Key {
	return NewKey(EndpointKey, endpointKeyID(tech, resource), opts...)
}

func endpointKeyID(tech, resource string) string {
	return tech + "/" + resource
}

// EndpointData describes an external device which may offer or accept calls
// to or from Asterisk.  Devices are defined by a technology/resource pair.
//
// Allowed states:  'unknown', 'offline', 'online'
type EndpointData struct {
	// Key is the cluster-unique identifier for this Endpoint
	Key *Key `json:"key"`

	ChannelIDs []string `json:"channel_ids"`     // List of channel Ids which are associated with this endpoint
	Resource   string   `json:"resource"`        // The endpoint's resource name
	State      string   `json:"state,omitempty"` // The state of the endpoint
	Technology string   `json:"technology"`      // The technology of the endpoint (e.g. SIP, PJSIP, DAHDI, etc)
}

// ID returns the ID for the endpoint
func (ed *EndpointData) ID() string {
	return ed.Technology + EndpointIDSeparator + ed.Resource
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

// NewEndpointHandle creates a new EndpointHandle
func NewEndpointHandle(key *Key, e Endpoint) *EndpointHandle {
	return &EndpointHandle{
		key: key,
		e:   e,
	}
}

// An EndpointHandle is a reference to an endpoint attached to
// a transport to an asterisk server
type EndpointHandle struct {
	key *Key
	e   Endpoint
}

// ID returns the identifier for the endpoint
func (eh *EndpointHandle) ID() string {
	return eh.key.ID
}

// Key returns the key for the endpoint
func (eh *EndpointHandle) Key() *Key {
	return eh.key
}

// Data returns the state of the endpoint
func (eh *EndpointHandle) Data() (*EndpointData, error) {
	return eh.e.Data(eh.key)
}
