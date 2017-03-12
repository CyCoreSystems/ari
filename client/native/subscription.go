package native

import "github.com/CyCoreSystems/ari"

type nativeSubscription struct {
	closeChan chan struct{}
	events    chan ari.Event
}

func (ns *nativeSubscription) Events() chan ari.Event {
	return ns.events
}

func (ns *nativeSubscription) Cancel() {
	if ns.closeChan != nil {
		close(ns.closeChan)
	}
}
