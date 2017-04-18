package ari

// Application represents a communication path interacting with an Asterisk
// server for application-level resources
type Application interface {

	// List returns the list of applications in Asterisk
	List() ([]*Key, error)

	// Get returns a handle to the application for further interaction
	Get(key *Key) ApplicationHandle

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
	BridgeIDs   []string `json:"bridge_ids"`   // Subscribed BridgeIds
	ChannelIDs  []string `json:"channel_ids"`  // Subscribed ChannelIds
	DeviceNames []string `json:"device_names"` // Subscribed Device names
	EndpointIDs []string `json:"endpoint_ids"` // Subscribed Endpoints (tech/resource format)
	Name        string   `json:"name"`         // Name of the application
}

// ApplicationHandle provides a wrapper to an Application interface for
// operations on a specific application
type ApplicationHandle interface {
	// ID returns the identifier for the application
	ID() string

	// Data retrives the data for the application
	Data() (ad *ApplicationData, err error)

	// Subscribe subscribes the application to an event source
	// event source may be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	Subscribe(eventSource string) (err error)

	// Unsubscribe unsubscribes (removes a subscription to) a given
	// ARI application from the provided event source
	// Equivalent to DELETE /applications/{applicationName}/subscription
	Unsubscribe(eventSource string) (err error)

	// Match returns true fo the event matches the application
	Match(evt Event) bool
}
