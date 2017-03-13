package stdbus

import (
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
)

var dtmfTestEventData = `
{
	"asterisk_id": "aa:bb:cc:dd:ee:ff",
	"type": "ChannelDtmfReceived",
	"application": "test",
	"timestamp": "2017-03-13T17:15:40.000-0400",
	"channel": {
		"accountcode": "testAccount",
		"caller": {
			"name": "testCallerName",
			"number": "testCallerNumber"
		},
		"connected": {
			"name": "testConnectedName",
			"number": "testConnectedNumber"
		},
		"creationtime": "2017-03-13T17:15:39.000-0400",
		"dialplan": {

		},
		"id": "testChannelID"
	},
	"digit": "1",
	"duration_ms": "100"
}
`

var dtmfTestEvent *ari.Message

func init() {
	var err error
	dtmfTestEvent, err = ari.NewMessage([]byte(dtmfTestEventData))
	if err != nil {
		panic("failed to construct dtmf test event")
	}
}

func TestSubscribe(t *testing.T) {
	b := &bus{
		subs: []*subscription{},
	}

	defer b.Close()

	sub := b.Subscribe(ari.Events.ChannelDtmfReceived)
	if len(b.subs) != 1 {
		t.Error("failed to add subscription to bus")
	}
	sub.Cancel()
}

func TestClose(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Close caused a panic")
		}
	}()

	b := New()
	sub := b.Subscribe(ari.Events.ChannelDtmfReceived)
	sub.Cancel()
	sub.Cancel()

	sub2 := b.Subscribe(ari.Events.ChannelDestroyed).(*subscription)

	b.Close()
	b.Close()

	if !sub2.closed {
		t.Error("subscription was not marked as closed")
		return
	}

	select {
	case _, ok := <-sub2.C:
		if ok {
			t.Error("subscription channel is not closed")
			return
		}
	default:
	}

}

func TestEvents(t *testing.T) {
	b := New()
	defer b.Close()

	sub := b.Subscribe(ari.Events.ChannelDtmfReceived)
	defer sub.Cancel()

	b.Send(dtmfTestEvent)

	select {
	case <-time.After(time.Millisecond):
		t.Error("failed to receive event")
		return
	case e, ok := <-sub.Events():
		if !ok {
			t.Error("events channel was closed")
			return
		}
		if e == nil {
			t.Error("received empty event")
			return
		}
	}
}
