package ari

// The event interfaces are most useful when
// checking whether a random "event" type
// is under a specific group:
//
//		ch, ok := evt.(ChannelEvent)
//		if ok { // event is for a channel }
//

// A ChannelEvent is an event with a channel ID
type ChannelEvent interface {
	ChannelID() string
}
