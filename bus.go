package ari

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription model.
type Bus interface {
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
	Events() chan Event
	Cancel()
}
