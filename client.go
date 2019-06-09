package ari

// Client represents a set of operations to interact
// with an Asterisk ARI server.  It is agnostic to transport
// and implementation.
type Client interface {

	// ApplicationName returns the ARI application name by which this client is connected
	ApplicationName() string

	// Bus returns the event bus of the client
	Bus() Bus

	// Connected indicates whether the Websocket is connected
	Connected() bool

	// Close shuts down the client
	Close()

	//
	//  ARI Namespaces
	//

	// Application access the Application ARI namespace
	Application() Application

	// Asterisk accesses the Asterisk ARI namespace
	Asterisk() Asterisk

	// Bridge accesses the Bridge ARI namespace
	Bridge() Bridge

	// Channel accesses the Channel ARI namespace
	Channel() Channel

	// DeviceState accesses the DeviceState ARI namespace
	DeviceState() DeviceState

	// Endpoint accesses the Endpoint ARI namespace
	Endpoint() Endpoint

	// LiveRecording accesses the LiveRecording ARI namespace
	LiveRecording() LiveRecording

	// Mailbox accesses the Mailbox ARI namespace
	Mailbox() Mailbox

	// Playback accesses the Playback ARI namespace
	Playback() Playback

	// Sound accesses the Sound ARI namespace
	Sound() Sound

	// StoredRecording accesses the StoredRecording ARI namespace
	StoredRecording() StoredRecording

	// TextMessage accesses the TextMessage ARI namespace
	TextMessage() TextMessage
}
