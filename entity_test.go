package ari

import (
	"fmt"
	"testing"
	"time"
)

type EntityFromEventTest struct {
	Desc   string
	Input  Event
	Output []Entity
}

func (efe *EntityFromEventTest) Test() string {
	list := EntitiesFromEvent(efe.Input)
	list = []Entity(EntitySlice(list).SortBy(byID))
	efe.Output = []Entity(EntitySlice(efe.Output).SortBy(byID))

	failed := !EntitySlice(list).Eqv(EntitySlice(efe.Output))

	if failed {
		return fmt.Sprintf("%s: expected %s; got %s", efe.Desc, efe.Output, list)
	}

	return ""
}

func (e EntitySlice) Eqv(r EntitySlice) bool {
	if len(e) != len(r) {
		return false
	}

	for i, v := range e {
		if v.ID != r[i].ID || v.Type != r[i].Type {
			return false
		}
	}

	return true
}

func byID(l Entity, r Entity) bool {
	return l.ID < r.ID
}

func channelData(id string, name string) ChannelData {
	return ChannelData{ID: id, Name: name}
}

func TestEntitesFrom(t *testing.T) {

	channelDtmfReceivedEvent := &ChannelDtmfReceived{
		EventData: EventData{
			Application: "testApp",
			Timestamp:   DateTime(time.Now()),
		},
		Channel: channelData("testChannelEvent", "Local/testChannel"),
	}
	bridgeCreatedEvent := &BridgeCreated{
		Bridge: BridgeData{
			ID:         "testBridgeEvent",
			ChannelIDs: []string{"channel1", "channel2"},
		},
	}

	endpointStateChangeEvent := &EndpointStateChange{
		Endpoint: EndpointData{
			ChannelIDs: []string{"channel1", "channel2"},
			Resource:   "rsc",
			State:      "state1",
			Technology: "tech1",
		},
	}

	playbackStartedOnChannelEvent := &PlaybackStarted{
		Playback: PlaybackData{
			ID:        "playback1",
			TargetURI: "channel:channel1",
		},
	}

	playbackStartedOnBridgeEvent := &PlaybackStarted{
		Playback: PlaybackData{
			ID:        "playback1",
			TargetURI: "bridge:bridge1",
		},
	}

	recordingStartedOnChannelEvent := &RecordingStarted{
		Recording: LiveRecordingData{
			Name:      "rec1",
			TargetURI: "channel:channel1",
		},
	}

	recordingFinishedOnChannelEvent := &RecordingStarted{
		Recording: LiveRecordingData{
			Name:      "rec1",
			TargetURI: "channel:channel1",
		},
	}

	recordingStartedOnBridgeEvent := &RecordingStarted{
		Recording: LiveRecordingData{
			Name:      "rec1",
			TargetURI: "bridge:bridge1",
		},
	}

	recordingFinishedOnBridgeEvent := &RecordingStarted{
		Recording: LiveRecordingData{
			Name:      "rec1",
			TargetURI: "bridge:bridge1",
		},
	}

	var tests = []EntityFromEventTest{
		{Desc: "Channel DTMF Received",
			Input:  channelDtmfReceivedEvent,
			Output: []Entity{{Type: "channel", ID: "testChannelEvent"}},
		},
		{Desc: "Bridge Created Event",
			Input: bridgeCreatedEvent,
			Output: []Entity{
				{Type: "bridge", ID: "testBridgeEvent"},
				{Type: "channel", ID: "channel1"},
				{Type: "channel", ID: "channel2"}},
		},
		{Desc: "Endpoint State Change Event",
			Input: endpointStateChangeEvent,
			Output: []Entity{
				{Type: "endpoint", ID: "tech1/rsc"},
				{Type: "channel", ID: "channel1"},
				{Type: "channel", ID: "channel2"}},
		},
		{Desc: "Playback Started On Channel Event",
			Input: playbackStartedOnChannelEvent,
			Output: []Entity{
				{Type: "playback", ID: "playback1"},
				{Type: "channel", ID: "channel1"}},
		},
		{Desc: "Playback Started On Bridge Event",
			Input: playbackStartedOnBridgeEvent,
			Output: []Entity{
				{Type: "playback", ID: "playback1"},
				{Type: "bridge", ID: "bridge1"}},
		},
		{Desc: "Recording Started On Channel Event",
			Input: recordingStartedOnChannelEvent,
			Output: []Entity{
				{Type: "recording", ID: "rec1"},
				{Type: "channel", ID: "channel1"}},
		},
		{Desc: "Recording Finished On Channel Event",
			Input: recordingFinishedOnChannelEvent,
			Output: []Entity{
				{Type: "recording", ID: "rec1"},
				{Type: "channel", ID: "channel1"}},
		},
		{Desc: "Recording Started On Bridge Event",
			Input: recordingStartedOnBridgeEvent,
			Output: []Entity{
				{Type: "recording", ID: "rec1"},
				{Type: "bridge", ID: "bridge1"}},
		},

		{Desc: "Recording Finished On Bridge Event",
			Input: recordingFinishedOnBridgeEvent,
			Output: []Entity{
				{Type: "recording", ID: "rec1"},
				{Type: "bridge", ID: "bridge1"}},
		},
		{Desc: "Recording Started On Bridge Event",
			Input: recordingStartedOnBridgeEvent,
			Output: []Entity{
				{Type: "recording", ID: "rec1"},
				{Type: "bridge", ID: "bridge1"}},
		},
	}

	for _, tx := range tests {
		res := tx.Test()
		if res != "" {
			t.Errorf(res)
		}
	}
}
