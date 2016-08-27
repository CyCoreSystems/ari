package ari

// Channel represents a communication path interacting with an Asterisk server.
type Channel interface {
	// Get returns a handle to a channel for further interaction
	Get(id string) *ChannelHandle

	// Create creates a new channel, returning a handle to it or an
	// error, if the creation failed
	Create() (*ChannelHandle, error)

	// Data returns the channel data for a given channel
	Data(id string) ChannelData

	// Continue tells Asterisk to return a channel to the dialplan
	Continue(id, context, extension, priority string) error

	// Busy hangs up the channel with the "busy" cause code
	Busy(id string) error

	// TODO: rest of interface
}

// NewChannelHandle returns a handle to the given ARI channel
func NewChannelHandle(id string, c Channel) *ChannelHandle {
	return &ChannelHandle{
		id: id,
		c:  c,
	}
}

// ChannelHandle provides a wrapper to a Channel interface for
// operations on a particular channel ID
type ChannelHandle struct {
	id string  // id of the channel on which we are operating
	c  Channel // the Channel interface with which we are operating
}

// Data returns the channel's data
func (ch *ChannelHandle) Data() ChannelData {
	return ch.c.Data(ch.id)
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension, priority string) error {
	return ch.c.Continue(ch.id, context, extension, priority)
}

// Busy hangs up the channel with the "busy" cause code
func (ch *ChannelHandle) Busy() error {
	return ch.c.Busy(ch.id)
}

// TODO: rest of ChannelHandle
