package ari

import (
	"sync"

	"golang.org/x/net/context"
)

// ALL signifies that the subscriber wants all events
const ALL string = "all"

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription
// model.
type Bus struct {
	subs []*Subscription // The list of subscriptions

	mu sync.Mutex
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

func (b *Bus) send(msg *Message) {
	var e Eventer
	switch msg.Type {
	case "BridgeAttendedTransfer":
		e = &BridgeAttendedTransfer{}
	case "BridgeBlindTransfer":
		e = &BridgeBlindTransfer{}
	case "BridgeCreated":
		e = &BridgeCreated{}
	case "BridgeDestroyed":
		e = &BridgeDestroyed{}
	case "BridgeMerged":
		e = &BridgeMerged{}
	case "ChannelCallerId":
		e = &ChannelCallerId{}
	case "ChannelConnectedLine":
		e = &ChannelConnectedLine{}
	case "ChannelCreated":
		e = &ChannelCreated{}
	case "ChannelDestroyed":
		e = &ChannelDestroyed{}
	case "ChannelDialplan":
		e = &ChannelDialplan{}
	case "ChannelDtmfReceived":
		e = &ChannelDtmfReceived{}
	case "ChannelEnteredBridge":
		e = &ChannelEnteredBridge{}
	case "ChannelHangupRequest":
		e = &ChannelHangupRequest{}
	case "ChannelHold":
		e = &ChannelHold{}
	case "ChannelLeftBridge":
		e = &ChannelLeftBridge{}
	case "ChannelStateChange":
		e = &ChannelStateChange{}
	case "ChannelTalkingFinished":
		e = &ChannelTalkingFinished{}
	case "ChannelTalkingStarted":
		e = &ChannelTalkingStarted{}
	case "ChannelUnhold":
		e = &ChannelUnhold{}
	case "ChannelUserevent":
		e = &ChannelUserevent{}
	case "ChannelVarset":
		e = &ChannelVarset{}
	case "ContactStatusChanged":
		e = &ContactStatusChanged{}
	case "DeviceStateChanged":
		e = &DeviceStateChanged{}
	case "Dial":
		e = &Dial{}
	case "EndpointStateChange":
		e = &EndpointStateChange{}
	case "PeerStatusChange":
		e = &PeerStatusChange{}
	case "PlaybackFinished":
		e = &PlaybackFinished{}
	case "PlaybackStarted":
		e = &PlaybackStarted{}
	case "RecordingFailed":
		e = &RecordingFailed{}
	case "RecordingFinished":
		e = &RecordingFinished{}
	case "RecordingStarted":
		e = &RecordingStarted{}
	case "StasisEnd":
		e = &StasisEnd{}
	case "StasisStart":
		e = &StasisStart{}
	case "TextMessageReceived":
		e = &TextMessageReceived{}
	default:
		Logger.Debug("Unhandled event received:", "type", msg.Type)
		e = &Event{}
	}
	err := msg.DecodeAs(e)
	if err != nil {
		Logger.Error("Failed to decode message", "error", err)
		return
	}
	Logger.Debug("Received event", "event", e)

	// Disseminate the message to the subscribers
	for _, s := range b.subs {
		for _, topic := range s.events {
			if topic == e.GetType() || topic == ALL {
				select {
				case s.C <- e:
				default: // never block
				}
			}
		}
	}
}

// StartBus creates and returns the event bus.
func StartBus(ctx context.Context) *Bus {
	b := &Bus{
		subs: []*Subscription{},
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
}

// A Subscription is a wrapped channel for receiving
// events from the ARI event bus.
type Subscription struct {
	b      *Bus         // reference to the event bus
	events []string     // list of events to listen for
	C      chan Eventer // channel for sending events to the subscriber
	mu     sync.Mutex
}

// Subscribe returns a subscription to the given list
// of event types
func (b *Bus) Subscribe(eTypes ...string) *Subscription {
	s := &Subscription{
		b:      b,
		events: eTypes,
		C:      make(chan Eventer, 1),
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
func (s *Subscription) Next(ctx context.Context) Eventer {
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
func (b *Bus) Once(ctx context.Context, eTypes ...string) <-chan Eventer {
	s := b.Subscribe(eTypes...)
	ret := make(chan Eventer, 1)

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
