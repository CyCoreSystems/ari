package stdbus

import (
	"testing"
	"time"

	"github.com/Amtelco-Software/ari/v6"
)

var dtmfTestEventData = `
{
  "channel": {
    "id": "9ae755c1-28a1-11e7-a1b1-0a580a480105",
    "dialplan": {
      "priority": 1,
      "context": "default",
      "exten": "9ae88e27-28a1-11e7-ba20-0a580a480707"
    },
    "creationtime": "2017-04-24T03:53:41.188+0000",
    "name": "Local/9ae88e27-28a1-11e7-ba20-0a580a480707@default-0000008b;1",
    "state": "Up",
    "connected": {
      "name": "",
      "number": ""
    },
    "caller": {
      "name": "",
      "number": ""
    },
    "accountcode": "",
    "language": "en"
  },
  "duration_ms": 240,
  "type": "ChannelDtmfReceived",
  "application": "sdp",
  "timestamp": "2017-04-24T03:53:42.155+0000",
  "digit": "1",
  "asterisk_id": "42:01:0a:64:00:06"
}
`

var dtmfTestEvent ari.Event

func init() {
	var err error

	dtmfTestEvent, err = ari.DecodeEvent([]byte(dtmfTestEventData))
	if err != nil {
		panic("failed to construct dtmf test event: " + err.Error())
	}
}

func TestSubscribe(t *testing.T) {
	b := &bus{
		subs: []*subscription{},
	}

	defer b.Close()

	sub := b.Subscribe(nil, ari.Events.ChannelDtmfReceived)

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
	sub := b.Subscribe(nil, ari.Events.ChannelDtmfReceived)
	sub.Cancel()
	sub.Cancel()

	sub2 := b.Subscribe(nil, ari.Events.ChannelDestroyed).(*subscription)

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

	sub := b.Subscribe(nil, ari.Events.All)
	defer sub.Cancel()

	b.Send(dtmfTestEvent)

	select {
	case <-time.After(time.Millisecond):
		t.Error("failed to receive event")
		return
	case e, ok := <-sub.Events():
		t.Log("event received")

		if !ok {
			t.Error("events channel was closed")
			return
		}

		if e == nil {
			t.Error("received empty event")
			return
		}

		dtmf, ok := e.(*ari.ChannelDtmfReceived)
		if !ok {
			t.Errorf("event is not a DTMF received event")
			return
		}

		if dtmf.Channel.ID != "9ae755c1-28a1-11e7-a1b1-0a580a480105" {
			t.Errorf("Failed to parse channel subentity on DTMF event")
			return
		}
	}
}

func TestEventsMultipleKeys(t *testing.T) {
	b := New()
	defer b.Close()

	sub := b.Subscribe(nil, ari.Events.All)
	defer sub.Cancel()

	multiKeyEvent := ari.BridgeCreated{
		Bridge: ari.BridgeData{
			ID:         "A",
			ChannelIDs: []string{"x", "y"},
		},
	}

	keys := multiKeyEvent.Keys()
	if len(keys) != 5 {
		t.Errorf("Expected BridgeCreated.Keys() to be 2, got '%v'", len(keys))
	}

	b.Send(&multiKeyEvent)

	eventCount := 0
L:
	for {
		select {
		case <-time.After(time.Millisecond):
			break L
		case <-sub.Events():
			eventCount++
		}
	}

	if eventCount != 1 {
		t.Errorf("Expected 1 event to be sent, got %v", eventCount)
	}
}
