package ari

// Client represents a set of operations to interact
// with an Asterisk ARI server.  It is agnostic to transport
// and implementation.
type Client struct {
	// Namespaced Interfaces
	Application Application
	Asterisk    Asterisk
	Channel     Channel
	Bridge      Bridge
	Playback    Playback
	Mailbox     Mailbox
	Endpoint    Endpoint
	TextMessage TextMessage
	Bus         Bus

	// TODO: other interfaces
}
