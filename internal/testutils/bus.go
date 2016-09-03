// +build test

package testutils

import (
	"sync"

	v2 "github.com/CyCoreSystems/ari/v2"
)

// Bus is a testing version of the event bus
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]*v2.Subscription
}

// NewBus creates a new bus
func NewBus() *Bus {
	return &Bus{
		subs: make(map[string][]*v2.Subscription),
	}
}

// Subscribe returns a subscription to the given list of events
func (bus *Bus) Subscribe(nx ...string) (a *v2.Subscription) {

	a = v2.NewSubscription("")

	bus.mu.Lock()

	for _, n := range nx {
		if _, ok := bus.subs[n]; !ok {
			bus.subs[n] = make([]*v2.Subscription, 0)
		}
		bus.subs[n] = append(bus.subs[n], a)
	}

	bus.mu.Unlock()

	return a
}

// Send sends an event to the event name
func (bus *Bus) Send(evt v2.Eventer) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	for _, l := range bus.subs[evt.GetType()] {
		if !l.Closed {
			l.C <- evt
		}
	}
}

// SendTo sends an event to the event name
func (bus *Bus) SendTo(n string, evt v2.Eventer) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	for _, l := range bus.subs[n] {
		if !l.Closed {
			l.C <- evt
		}
	}
}
