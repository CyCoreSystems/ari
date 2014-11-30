package ari

import (
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

type Distributor interface {
	Publish(*Event)
	Subscribe() *Subscription
}

// Distribution is a convenience tool for distribution
// of ARI events to subscribers.  It is essentially a
// pubsub interface for ARI events.  It contains no filtering
// and does not automatically attach to ARI for event sourcing.
// Instead, it merely provides a tool to redistribute received
// ARI events.
type Distribution struct {
	subscribers     map[string]chan *Event
	subscribersLock sync.Mutex
}

// Publish pushes an event to all subscribers to this distribution
func (d *Distribution) Publish(e *Event) {
	// Do nothing if we have no subscribers
	if d.subscribers == nil {
		return
	}

	d.subscribersLock.Lock()
	defer d.subscribersLock.Unlock()

	for _, s := range d.subscribers {
		s <- e
	}
}

// Subscribe returns a subscription to this distribution
// NOTE: it is the duty of subscribers to `Cancel` their own
// subscriptions.
func (d *Distribution) Subscribe() *Subscription {
	eventChan := make(chan *Event)
	return d.SubscribeChan(eventChan)
}

// SubscribeChan returns a subscription to this distribution
// using the provided channel as the event sink
// NOTE: it is the duty of subscribers to `Cancel` their own
// subscriptions.
func (d *Distribution) SubscribeChan(eventChan chan *Event) *Subscription {
	id := uuid.New()

	d.subscribersLock.Lock()
	defer d.subscribersLock.Unlock()

	if d.subscribers == nil {
		d.subscribers = make(map[string]chan *Event)
	}

	d.subscribers[id] = eventChan

	s := &Subscription{
		id:           id,
		distribution: d,
		eventChan:    eventChan,
	}
	return s
}

// Unsubscribe removes a subscription from the distribution.  This
// should not be called directly, but rather, subscribers should call
// their subscription's `Cancel` method.
func (d *Distribution) Unsubscribe(id string) {
	// Do nothing if we have no subscribers
	if d.subscribers == nil {
		return
	}

	d.subscribersLock.Lock()
	defer d.subscribersLock.Unlock()

	_, exists := d.subscribers[id]
	if exists {
		delete(d.subscribers, id)
	}
}

type Subscription struct {
	distribution *Distribution // subscribed distribution
	id           string        // id of this subscription
	eventChan    chan *Event   // chan for receiving events
}

// Next blocks to return the next event
func (s *Subscription) Next() (*Event, bool) {
	e, open := <-s.eventChan
	return e, open
}

// Cancel removes a subscription and closes
// its channel
func (s *Subscription) Cancel() {
	s.distribution.Unsubscribe(s.id)
	close(s.eventChan)
}

//Allows for the injection of an event channel into a Subscription method
func (s *Subscription) Inject(eChan chan *Event) {
	s.eventChan = eChan
}
