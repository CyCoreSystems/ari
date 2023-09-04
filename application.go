package ari

// Application represents a communication path interacting with an Asterisk
// server for application-level resources
type Application interface {
	// List returns the list of applications in Asterisk, optionally using the key for filtering
	List(*Key) ([]*Key, error)

	// Get returns a handle to the application for further interaction
	Get(key *Key) *ApplicationHandle

	// Data returns the applications data
	Data(key *Key) (*ApplicationData, error)

	// Subscribe subscribes the given application to an event source
	// event source may be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	Subscribe(key *Key, eventSource string) error

	// Unsubscribe unsubscribes (removes a subscription to) a given
	// ARI application from the provided event source
	// Equivalent to DELETE /applications/{applicationName}/subscription
	Unsubscribe(key *Key, eventSource string) error
}

// ApplicationData describes the data for a Stasis (Ari) application
type ApplicationData struct {
	// Key is the unique identifier for this application instance in the cluster
	Key *Key `json:"key"`

	BridgeIDs   []string `json:"bridge_ids"`   // Subscribed BridgeIds
	ChannelIDs  []string `json:"channel_ids"`  // Subscribed ChannelIds
	DeviceNames []string `json:"device_names"` // Subscribed Device names
	EndpointIDs []string `json:"endpoint_ids"` // Subscribed Endpoints (tech/resource format)
	Name        string   `json:"name"`         // Name of the application
}

// ApplicationHandle provides a wrapper to an Application interface for
// operations on a specific application
type ApplicationHandle struct {
	key *Key
	a   Application
}

// NewApplicationHandle creates a new handle to the application name
func NewApplicationHandle(key *Key, app Application) *ApplicationHandle {
	return &ApplicationHandle{
		key: key,
		a:   app,
	}
}

// ID returns the identifier for the application
func (ah *ApplicationHandle) ID() string {
	return ah.key.ID
}

// Key returns the key of the application
func (ah *ApplicationHandle) Key() *Key {
	return ah.key
}

// Data retrives the data for the application
func (ah *ApplicationHandle) Data() (ad *ApplicationData, err error) {
	ad, err = ah.a.Data(ah.key)
	return
}

// Subscribe subscribes the application to an event source
// event source may be one of:
//   - channel:<channelId>
//   - bridge:<bridgeId>
//   - endpoint:<tech>/<resource> (e.g. SIP/102)
//   - deviceState:<deviceName>
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
func (ah *ApplicationHandle) Match(e Event) bool {
	return e.GetApplication() == ah.key.ID
}
