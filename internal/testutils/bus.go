package testutils

import (
	v2 "github.com/CyCoreSystems/ari/v2"
)

// Bus is a testing version of the event bus
type Bus struct {
	subs map[string]*v2.Subscription
}

// NewBus creates a new bus
func NewBus() *Bus {
	return &Bus{
		subs: make(map[string]*v2.Subscription),
	}
}

// Subscribe returns a subscription to the given list of events
func (bus *Bus) Subscribe(nx ...string) *v2.Subscription {
	a := v2.NewSubscription("")

	for _, n := range nx {
		bus.subs[n] = a
	}

	return a
}

// Send sends an event to the event name
func (bus *Bus) Send(evt v2.Eventer) {
	bus.subs[evt.GetType()].C <- evt
}

// SendTo sends an event to the event name
func (bus *Bus) SendTo(n string, evt v2.Eventer) {
	bus.subs[n].C <- evt
}
