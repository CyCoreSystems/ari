package stdbus

import (
	"sync"

	"github.com/CyCoreSystems/ari"

	"golang.org/x/net/context"
)

// busChannelBuffer defines the buffer size of the subscription channels
var busChannelBuffer = 100

// bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription
// model.
type bus struct {
	subs []*subscription // The list of subscriptions

	mu sync.Mutex
}

func (b *bus) addSubscription(s *subscription) {
	b.mu.Lock()
	b.subs = append(b.subs, s)
	b.mu.Unlock()
}

func (b *bus) removeSubscription(s *subscription) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, si := range b.subs {
		if s == si {
			// Subs are pointers, so we have to explicitly remove them
			// to prevent memory leaks
			b.subs[i] = b.subs[len(b.subs)-1] // replace the current with the end
			b.subs[len(b.subs)-1] = nil       // remove the end
			b.subs = b.subs[:len(b.subs)-1]   // lop off the end
			return
		}
	}
}

// Send sends the message on the bus
func (b *bus) Send(msg *ari.Message) {
	b.send(msg)
}

func (b *bus) send(msg *ari.Message) {
	e := ari.Events.Parse(msg)

	//	Logger.Debug("Received event", "event", e)

	// Disseminate the message to the subscribers
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.subs {
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

// Start creates and returns the event bus.
func Start(ctx context.Context) ari.Bus {
	b := &bus{
		subs: []*subscription{},
	}

	// Listen for stop and shut down subscriptions, as required
	go func() {
		<-ctx.Done()
		b.Stop()
		return
	}()

	return b
}

// Stop the bus.  Cancels all subscriptions
// and stops listening for events.
func (b *bus) Stop() {
	// Close all subscriptions
	b.mu.Lock()
	if b.subs != nil {
		for i, s := range b.subs {
			s.closeChan()
			b.subs[i] = nil
		}
		b.subs = nil
	}
	b.mu.Unlock()
}

// A Subscription is a wrapped channel for receiving
// events from the ARI event bus.
type subscription struct {
	b      *bus           // reference to the event bus
	events []string       // list of events to listen for
	C      chan ari.Event // channel for sending events to the subscriber
	mu     sync.Mutex
	Closed bool
}

// newSubscription creates a new, unattached subscription
func newSubscription(eTypes ...string) *subscription {
	return &subscription{
		events: eTypes,
		C:      make(chan ari.Event, busChannelBuffer),
	}
}

// Subscribe returns a subscription to the given list
// of event types
func (b *bus) Subscribe(eTypes ...string) ari.Subscription {
	return b.subscribe(eTypes...)
}

// subscribe returns a subscription to the given list
// of event types
func (b *bus) subscribe(eTypes ...string) *subscription {
	s := newSubscription(eTypes...)
	s.b = b
	b.addSubscription(s)
	return s
}

// Events returns the events channel
func (s *subscription) Events() chan ari.Event {
	return s.C
}

// Next blocks for the next event in the subscription,
// returning that event when it arrives or nil if
// the subscription is canceled.
// Normally, one would listen to subscription.C directly,
// but this is a convenience function for providing a
// context to alternately cancel.
func (s *subscription) Next(ctx context.Context) ari.Event {
	select {
	case <-ctx.Done():
		return nil
	case e := <-s.C:
		return e
	}
}

func (s *subscription) closeChan() {
	s.mu.Lock()
	if s.C != nil {
		close(s.C)
		s.C = nil
	}
	s.Closed = true
	s.mu.Unlock()
}

// Cancel cancels the subscription and removes it from
// the event bus.
func (s *subscription) Cancel() {
	if s == nil {
		return
	}
	if s.b != nil {
		s.b.removeSubscription(s)
	}
	s.closeChan()
}

// Once listens for the first event of the provided types,
// returning a channel which supplies that event.
func (b *bus) Once(ctx context.Context, eTypes ...string) <-chan ari.Event {
	s := b.subscribe(eTypes...)

	ret := make(chan ari.Event, busChannelBuffer)

	// Stop subscription after one event
	go func() {
		select {
		case ret <- <-s.C:
		case <-ctx.Done():
		}
		close(ret)
		s.Cancel()
	}()
	return ret
}
