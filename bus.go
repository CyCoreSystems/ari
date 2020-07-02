package ari

import (
	"context"
	"sync"
)

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription model.
type Bus interface {
	Close()
	Sender
	Subscriber
}

// A Sender is an entity which can send event bus messages
type Sender interface {
	Send(e Event)
}

// A Subscriber is an entity which can create ARI event subscriptions
type Subscriber interface {
	Subscribe(key *Key, n ...string) Subscription
}

// A Subscription is a subscription on series of ARI events
type Subscription interface {
	// Events returns a channel on which events related to this subscription are sent.
	Events() <-chan Event

	// Cancel terminates the subscription
	Cancel()

	// Set callback
	SetCallback(func(s Subscription))
}

// Once listens for the first event of the provided types,
// returning a channel which supplies that event.
func Once(ctx context.Context, bus Bus, key *Key, eTypes ...string) <-chan Event {
	s := bus.Subscribe(key, eTypes...)

	ret := make(chan Event)

	// Stop subscription after one event
	go func() {
		select {
		case ret <- <-s.Events():
		case <-ctx.Done():
		}
		close(ret)
		s.Cancel()
	}()

	return ret
}

// NewNullSubscription returns a subscription which never returns any events
func NewNullSubscription() *NullSubscription {
	return &NullSubscription{
		ch: make(chan Event),
	}
}

// NullSubscription is a subscription which never returns any events.
type NullSubscription struct {
	ch     chan Event
	closed bool
	mu     sync.RWMutex
}

// Events implements the Subscription interface
func (n *NullSubscription) Events() <-chan Event {
	if n.ch == nil {
		n.mu.Lock()
		n.closed = false
		n.ch = make(chan Event)
		n.mu.Unlock()
	}

	return n.ch
}

// Cancel implements the Subscription interface
func (n *NullSubscription) Cancel() {
	if n.closed {
		return
	}

	n.mu.Lock()

	n.closed = true
	if n.ch != nil {
		close(n.ch)
	}

	n.mu.Unlock()
}
