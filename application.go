package ari

// Application represents a communication path interacting with an Asterisk
// server for application-level resources
type Application interface {

	// List returns the list of applications in Asterisk
	List() ([]*ApplicationHandle, error)

	// Get returns a handle to the application for further interaction
	Get(name string) *ApplicationHandle

	// Data returns the applications data
	Data(name string) (ApplicationData, error)

	// Subscribe subscribes the given application to an event source
	// event source may be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	Subscribe(name string, eventSource string) error

	// Unsubscribe unsubscribes (removes a subscription to) a given
	// ARI application from the provided event source
	// Equivalent to DELETE /applications/{applicationName}/subscription
	Unsubscribe(name string, eventSource string) error
}

// ApplicationData describes the data for a Stasis (Ari) application
type ApplicationData struct {
	BridgeIDs   []string `json:"bridge_ids"`   // Subscribed BridgeIds
	ChannelIDs  []string `json:"channel_ids"`  // Subscribed ChannelIds
	DeviceNames []string `json:"device_names"` // Subscribed Device names
	EndpointIDs []string `json:"endpoint_ids"` // Subscribed Endpoints (tech/resource format)
	Name        string   `json:"name"`         // Name of the application
}

// NewApplicationHandle creates a new handle to the application name
func NewApplicationHandle(name string, app Application) *ApplicationHandle {
	return &ApplicationHandle{
		name: name,
		a:    app,
	}
}

// ApplicationHandle provides a wrapper to an Application interface for
// operations on a specific application
type ApplicationHandle struct {
	name string
	a    Application
}

// ID returns the identifier for the application
func (ah *ApplicationHandle) ID() string {
	return ah.name
}

// Data retrives the data for the application
func (ah *ApplicationHandle) Data() (ad ApplicationData, err error) {
	ad, err = ah.a.Data(ah.name)
	return
}

// Subscribe subscribes the application to an event source
// event source may be one of:
//  - channel:<channelId>
//  - bridge:<bridgeId>
//  - endpoint:<tech>/<resource> (e.g. SIP/102)
//  - deviceState:<deviceName>
func (ah *ApplicationHandle) Subscribe(eventSource string) (err error) {
	err = ah.a.Subscribe(ah.name, eventSource)
	return
}

// Unsubscribe unsubscribes (removes a subscription to) a given
// ARI application from the provided event source
// Equivalent to DELETE /applications/{applicationName}/subscription
func (ah *ApplicationHandle) Unsubscribe(eventSource string) (err error) {
	err = ah.a.Unsubscribe(ah.name, eventSource)
	return
}

// Match returns true fo the event matches the application
func (ah *ApplicationHandle) Match(evt Event) bool {
	applicationEvent, ok := evt.(ApplicationEvent)
	if !ok {
		return false
	}
	if applicationEvent.GetApplication() == ah.name {
		return true
	}
	return false
}
