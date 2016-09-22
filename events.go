package ari

// Event is the top level event interface
type Event interface {
	MessageRawer
	ApplicationEvent
	GetType() string
}

// EventData is the base struct for all events
type EventData struct {
	Message
	Application string   `json:"application"`
	Timestamp   DateTime `json:"timestamp,omitempty"`
}

// GetApplication gets the application of the event
func (e *EventData) GetApplication() string {
	return e.Application
}

// GetType gets the type of the event
func (e *EventData) GetType() string {
	return e.Type
}

// The event interfaces are most useful when
// checking whether a random "event" type
// is under a specific group:
//
//		ch, ok := evt.(ChannelEvent)
//		if ok { // event is for a channel }
//

// An ApplicationEvent is an event with an application (which is every event actually)
type ApplicationEvent interface {
	GetApplication() string
}

// A ChannelEvent is an event with a channel ID
type ChannelEvent interface {
	GetChannelID() string
}

// A BridgeEvent is an event with a Bridge ID
type BridgeEvent interface {
	GetBridgeID() string
}
