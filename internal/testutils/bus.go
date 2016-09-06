package testutils

import (
	"sync"

	v2 "github.com/CyCoreSystems/ari/v2"
)

// Bus is a testing version of the event bus
type Bus struct {
	mu      sync.RWMutex
	subs    map[string][]*v2.Subscription
	expects map[string]chan struct{}
}

// NewBus creates a new bus
func NewBus() *Bus {
	return &Bus{
		subs:    make(map[string][]*v2.Subscription),
		expects: make(map[string]chan struct{}),
	}
}

// Expect returns a channel that will be closed when a subscription occurs
func (bus *Bus) Expect(n string) (ch chan struct{}) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	ch = make(chan struct{})
	bus.expects[n] = ch

	return
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

		if ch, ok := bus.expects[n]; ok {
			close(ch)
			delete(bus.expects, n)
		}

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
