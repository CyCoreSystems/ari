package ari

type Event interface {

	// Get returns a handle pointer to the Event for further interaction
	Get(name string) *EventHandle

	// Data returns the Event's data
	Data(name string) (EventData, error)
}

// EventHandle provides a wrapper to an Event interface for
// operations on a specific Event
type EventHandle struct {
	name string
	e    Event
}

// Event is the base struct for all events
type EventData struct {
	Message
	Application string   `json:"application"`
	Timestamp   DateTime `json:"timestamp,omitempty"`
}

//Request structure for creating a user event. Only Application is required.
type CreateUserEventRequest struct {
	Application string `json:"application"`
	Source      string `json:"source,omitempty"`
	Variables   string `json:"variables,omitempty"`
}

// NewEventHandle creates a new handle to the Event name
func NewEventHandle(name string, e Event) *EventHandle {
	return &EventHandle{
		name: name,
		e:    e,
	}
}

// Data retrieves the data for the Event
func (eh *EventHandle) Data() (ed EventData, err error) {
	ed, err = eh.e.Data(eh.name)
	return ed, err
}
