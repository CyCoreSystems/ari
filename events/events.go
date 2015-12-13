// Events is an ARI event pub-sub service
package events

import (
	"sync"

	"github.com/CyCoreSystems/ari"
	"golang.org/x/net/context"
)

// Bus is an event bus.  It receives and
// redistributes events based on a subscription
// model.
type Bus struct {
	a *ari.Client // The ARI client connection

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
			b.subs[len(b.subs)-1] = nil       // remvoe the end
			b.subs = b.subs[:len(b.subs)-1]   // lop off the end
			return
		}
	}
}

func (b *Bus) send(e *ari.Event) {
	for _, s := range b.subs {
		for _, topic := range s.events {
			if topic == e.Type {
				select {
				case s.eventsChan <- e:
				default:
				}
			}
		}
	}
}

// Start creates and returns the event bus.
func Start(ctx context.Context, a *ari.Client) *Bus {
	bCtx, cancel := context.WithCancel(ctx)

	b := &Bus{
		a:      a,
		cancel: cancel,
		subs:   []*Subscription{},
	}

	// Listen for events
	go func() {
		for {
			select {
			case <-bCtx.Done():
				b.Stop()
				return
			case e := <-a.Events:
				b.send(e)
			}
		}
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
			s.stop()
			b.subs[i] = nil
		}
		b.subs = nil
	}
	b.mu.Unlock()
	b.cancel()
}

type Subscription struct {
	events     []string        // list of events to listen for
	eventsChan chan *ari.Event // channel for sending events to the subscriber
	mu         sync.Mutex
}

// Subscribe returns a subscription to the given list
// of event types
func (b *Bus) Subscribe(eTypes ...string) *Subscription {
	s := &Subscription{
		events:     eTypes,
		eventsChan: make(chan *ari.Event),
	}
	b.addSubscription(s)
	return s
}

func (s *Subscription) closeChan() {
	s.mu.Lock()
	if s.eventsChan != nil {
		close(s.eventsChan)
		s.eventsChan = nil
	}
	s.mu.Unlock()
}

// Cancel cancels the subscription and removes it from
// the event bus.
func (s *Subscription) Cancel() {
	b.removeSubscription(s)
	s.closeChan()
}

// Once listens for the first event of the provided types,
// returning a channel which supplies that event.
func (b *Bus) Once(eTypes ...string) <-chan *ari.Event {
	s := b.Subscribe(eTypes)
	ret := make(chan *ari.Event, 1)

	// Stop subscription after one event
	go func() {
		ret <- s.events
		close(ret)
		s.Cancel()
	}()
	return ret
}
