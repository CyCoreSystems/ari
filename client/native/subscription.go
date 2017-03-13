package native

import "github.com/CyCoreSystems/ari"

var subscriptionEventBufferLength = 100

type nativeSubscription struct {
	closed bool

	closedChan chan struct{}
	events     chan ari.Event
}

func newSubscription() *nativeSubscription {
	return &nativeSubscription{
		closedChan: make(chan struct{}),
		events:     make(chan ari.Event, subscriptionEventBufferLength),
	}
}

func (ns *nativeSubscription) Events() <-chan ari.Event {
	return ns.events
}

func (ns *nativeSubscription) Cancel() {
	if !ns.closed {
		ns.closed = true
		close(ns.closedChan)
		close(ns.events)
	}
}
