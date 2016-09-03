// +build test

package testutils

import (
	"time"

	v2 "github.com/CyCoreSystems/ari/v2"
)

// DelayedBus is a version of bus where the event send methods are delayed
type DelayedBus struct {
	bus   *Bus
	delay time.Duration
}

// NewDelayedBus that adds a delay between events sent
func NewDelayedBus(delay time.Duration) *DelayedBus {
	return &DelayedBus{
		bus:   NewBus(),
		delay: delay,
	}
}

// Subscribe returns a subscription to the given list of events
func (bus *DelayedBus) Subscribe(nx ...string) (a *v2.Subscription) {
	a = bus.bus.Subscribe(nx...)
	return
}

// Send sends an event to the event name
func (bus *DelayedBus) Send(evt v2.Eventer) {
	<-time.After(bus.delay)
	bus.bus.Send(evt)
}

// SendTo sends an event to the event name
func (bus *DelayedBus) SendTo(n string, evt v2.Eventer) {
	<-time.After(bus.delay)
	bus.bus.SendTo(n, evt)
}
