package ari

import (
	"sync"

	"golang.org/x/net/context"
)

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription
// model.
type Bus struct {
	subs []*Subscription // The list of subscriptions
	mu   sync.Mutex

	cancel context.CancelFunc
}

func (b *Bus) addSubscription(s *Subscription) {
	b.mu.Lock()
	b.subs = append(b.subs, s)
	b.mu.Unlock()
}

func (b *Bus) removeSubscription(s *Subscription) {
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

func (b *Bus) send(e *Event) {
	for _, s := range b.subs {
		for _, topic := range s.events {
			if topic == e.Type {
				select {
				case s.C <- e:
				default:
				}
			}
		}
	}
}

// StartBus creates and returns the event bus.
func StartBus(ctx context.Context) *Bus {
	bCtx, cancel := context.WithCancel(ctx)

	b := &Bus{
		cancel: cancel,
		subs:   []*Subscription{},
	}

	// Listen for stop and shut down subscriptions, as required
	go func() {
		<-bCtx.Done()
		b.Stop()
		return
	}()

	return b
}

// Stop the bus.  Cancels all subscriptions
// and stops listening for events.
func (b *Bus) Stop() {
	// Close all subscriptions
	b.mu.Lock()
	if b.subs != nil {
		for i, s := range b.subs {
			s.Cancel()
			b.subs[i] = nil
		}
		b.subs = nil
	}
	b.mu.Unlock()
	b.cancel()
}

type Subscription struct {
	b      *Bus        // reference to the event bus
	events []string    // list of events to listen for
	C      chan *Event // channel for sending events to the subscriber
	mu     sync.Mutex
}

// Subscribe returns a subscription to the given list
// of event types
func (b *Bus) Subscribe(eTypes ...string) *Subscription {
	s := &Subscription{
		b:      b,
		events: eTypes,
		C:      make(chan *Event),
	}
	b.addSubscription(s)
	return s
}

// Next blocks for the next event in the subscription,
// returning that event when it arrives or nil if
// the subscription is canceled.
// Normally, one would listen to subscription.C directly,
// but this is a convenience function for providing a
// context to alternately cancel.
func (s *Subscription) Next(ctx context.Context) *Event {
	select {
	case <-ctx.Done():
		return nil
	case e := <-s.C:
		return e
	}
}

func (s *Subscription) closeChan() {
	s.mu.Lock()
	if s.C != nil {
		close(s.C)
		s.C = nil
	}
	s.mu.Unlock()
}

// Cancel cancels the subscription and removes it from
// the event bus.
func (s *Subscription) Cancel() {
	s.b.removeSubscription(s)
	s.closeChan()
}

// Once listens for the first event of the provided types,
// returning a channel which supplies that event.
func (b *Bus) Once(eTypes ...string) <-chan *Event {
	s := b.Subscribe(eTypes...)
	ret := make(chan *Event, 1)

	// Stop subscription after one event
	go func() {
		ret <- <-s.C
		close(ret)
		s.Cancel()
	}()
	return ret
}
