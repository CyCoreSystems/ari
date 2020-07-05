package stdbus

import (
	"sync"

	"github.com/CyCoreSystems/ari/v5"
)

// subscriptionEventBufferSize defines the number of events that each
// subscription will queue before accepting more events.
var subscriptionEventBufferSize = 100

// bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription
// model.
type bus struct {
	subs []*subscription // The list of subscriptions

	rwMux sync.RWMutex

	closed bool
}

// New creates and returns the event bus.
func New() ari.Bus {
	b := &bus{
		subs: []*subscription{},
	}

	return b
}

// Close closes out all subscriptions in the bus.
func (b *bus) Close() {
	if b.closed {
		return
	}

	b.closed = true

	for _, s := range b.subs {
		s.Cancel()
	}
}

// Send sends the message to the bus
func (b *bus) Send(e ari.Event) {
	var matched bool

	b.rwMux.RLock()

	// Disseminate the message to the subscribers
	for _, s := range b.subs {
		if s.closed {
			continue
		}
		matched = false
		for _, k := range e.Keys() {
			if matched {
				break
			}

			if s.key.Match(k) {
				matched = true

				for _, topic := range s.events {
					if topic == e.GetType() || topic == ari.Events.All {
						select {
						case s.C <- e:
						default: // never block
						}
					}
				}
			}
		}
	}

	b.rwMux.RUnlock()
}

// Subscribe returns a subscription to the given list
// of event types
func (b *bus) Subscribe(key *ari.Key, eTypes ...string) ari.Subscription {
	s := newSubscription(key, eTypes...)
	s.AddCancelCallback(b.remove)
	b.add(s)

	return s
}

// add appends a new subscription to the bus
func (b *bus) add(s *subscription) {
	b.rwMux.Lock()
	b.subs = append(b.subs, s)
	b.rwMux.Unlock()
}

// remove deletes the given subscription from the bus
func (b *bus) remove(s interface{}) {
	b.rwMux.Lock()
	for i, si := range b.subs {
		if s == si {
			// Subs are pointers, so we have to explicitly remove them
			// to prevent memory leaks
			b.subs[i] = b.subs[len(b.subs)-1] // replace the current with the end
			b.subs[len(b.subs)-1] = nil       // remove the end
			b.subs = b.subs[:len(b.subs)-1]   // lop off the end

			break
		}
	}
	b.rwMux.Unlock()
}

// A Subscription is a wrapped channel for receiving
// events from the ARI event bus.
type subscription struct {
	key    *ari.Key
	events []string // list of events to listen for

	mu     sync.Mutex
	closed bool           // channel closure protection flag
	C      chan ari.Event // channel for sending events to the subscriber
	cb     []func(d interface{}) // slice of callbacks to execute on Cancel()
}

// newSubscription creates a new, unattached subscription
func newSubscription(key *ari.Key, eTypes ...string) *subscription {
	return &subscription{
		key:    key,
		events: eTypes,
		C:      make(chan ari.Event, subscriptionEventBufferSize),
	}
}

// Events returns the events channel
func (s *subscription) Events() <-chan ari.Event {
	return s.C
}

// Cancel cancels the subscription and removes it from
// the event bus.
func (s *subscription) Cancel() {
	if s == nil {
		return
	}

	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return
	}

	s.closed = true

	s.mu.Unlock()

	// Iterate through callbacks
	for _,cb := range s.cb {
		cb(s)
	}

	// Close the subscription's deliver channel
	if s.C != nil {
		close(s.C)
	}
}

// Add cancel callback function
func (s *subscription) AddCancelCallback(cb func(d interface{})) {
	s.cb = append(s.cb, cb)
}