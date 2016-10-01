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
	DeviceState DeviceState
	TextMessage TextMessage
	Sound       Sound
	Recording   *Recording
	Bus         Bus

	// TODO: other interfaces

	Cleanup func() error
}

// Close closes the client and calls any implementation specific cleanup code
func (cl *Client) Close() error {
	if cl.Cleanup != nil {
		return cl.Cleanup()
	}
	return nil
}
