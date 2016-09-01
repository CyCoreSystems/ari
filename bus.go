package ari

import (
	v2 "github.com/CyCoreSystems/ari/v2"
)

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription model.
type Bus interface {
	Sender
	Subscriber
}

// A Sender is an entity which can send event bus messages
type Sender interface {
	Send(m *v2.Message)
}

// A Subscriber is an entity which can create ARI event subscriptions
type Subscriber interface {
	Subscribe(n ...string) *v2.Subscription
}
