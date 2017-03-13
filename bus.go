package ari

import "context"

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription model.
type Bus interface {
	Close()
	Sender
	Subscriber
}

// A Sender is an entity which can send event bus messages
type Sender interface {
	Send(m *Message)
}

// A Subscriber is an entity which can create ARI event subscriptions
type Subscriber interface {
	Subscribe(n ...string) Subscription
}

// A Subscription is a subscription on series of ARI events
type Subscription interface {
	// Events returns a channel on which events related to this subscription are sent.
	Events() <-chan Event

	// Cancel terminates the subscription
	Cancel()
}

// Once listens for the first event of the provided types,
// returning a channel which supplies that event.
func Once(ctx context.Context, bus Bus, eTypes ...string) <-chan Event {
	s := bus.Subscribe(eTypes...)

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
