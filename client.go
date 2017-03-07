package ari

// Client represents a set of operations to interact
// with an Asterisk ARI server.  It is agnostic to transport
// and implementation.
type Client interface {
	ApplicationName() string
	Close()

	Application() Application
	Asterisk() Asterisk
	Bridge() Bridge
	Bus() Bus
	Channel() Channel
	DeviceState() DeviceState
	Endpoint() Endpoint
	LiveRecording() LiveRecording
	Mailbox() Mailbox
	Playback() Playback
	Sound() Sound
	StoredRecording() StoredRecording
	TextMessage() TextMessage
}
