package ari

// Event is the top level event interface
type Event interface {
	MessageRawer
	ApplicationEvent
	GetType() string
}

// A Matcher is an entity which can query an event
type Matcher interface {
	Matches(evt Event) bool
}

// EventData is the base struct for all events
type EventData struct {
	Message
	Application string   `json:"application"`
	Timestamp   DateTime `json:"timestamp,omitempty"`
}

// GetApplication gets the application of the event
func (e *EventData) GetApplication() string {
	return e.Application
}

// GetType gets the type of the event
func (e *EventData) GetType() string {
	return e.Type
}

// The event interfaces are most useful when
// checking whether a random "event" type
// is under a specific group:
//
//		ch, ok := evt.(ChannelEvent)
//		if ok { // event is for a channel }
//

// An ApplicationEvent is an event with an application (which is every event actually)
type ApplicationEvent interface {
	GetApplication() string
}

// A ChannelEvent is an event with one or more channel IDs
type ChannelEvent interface {
	GetChannelIDs() []string
}

// A BridgeEvent is an event with one or more Bridge IDs
type BridgeEvent interface {
	GetBridgeIDs() []string
}

// An EndpointEvent is an event with one or more endpoint IDs
type EndpointEvent interface {
	GetEndpointIDs() []string
}

// A PlaybackEvent is an event with one or more playback IDs
type PlaybackEvent interface {
	GetPlaybackIDs() []string
}

// A RecordingEvent is an event with one or more recording IDs
type RecordingEvent interface {
	GetRecordingIDs() []string
}

// implementations of events

// GetBridgeIDs gets the bridge IDs for the event
func (evt *BridgeAttendedTransfer) GetBridgeIDs() (sx []string) {
	if id := evt.DestinationThreewayBridge.ID; id != "" {
		sx = append(sx, id)
	}

	if id := evt.TransfererFirstLegBridge.ID; id != "" {
		sx = append(sx, id)
	}

	if id := evt.TransfererSecondLegBridge.ID; id != "" {
		sx = append(sx, id)
	}

	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *BridgeAttendedTransfer) GetChannelIDs() (sx []string) {
	if id := evt.DestinationLinkFirstLeg.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.DestinationLinkSecondLeg.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.DestinationThreewayChannel.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.ReplaceChannel.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.Transferee.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.TransfererFirstLeg.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.TransfererSecondLeg.ID; id != "" {
		sx = append(sx, id)
	}
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *BridgeBlindTransfer) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *BridgeBlindTransfer) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	if id := evt.ReplaceChannel.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.Transferee.ID; id != "" {
		sx = append(sx, id)
	}
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *BridgeCreated) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *BridgeDestroyed) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *BridgeMerged) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	sx = append(sx, evt.BridgeFrom.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelCallerID) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelCreated) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelDialplan) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelDtmfReceived) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelEnteredBridge) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *ChannelEnteredBridge) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelHangupRequest) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelHold) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelLeftBridge) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *ChannelLeftBridge) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelStateChange) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelTalkingStarted) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelUnhold) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelUserevent) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *ChannelUserevent) GetBridgeIDs() (sx []string) {
	sx = append(sx, evt.Bridge.ID)
	return
}

// GetEndpointIDs gets the bridge IDs for the event
func (evt *ChannelUserevent) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *ChannelVarset) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetEndpointIDs gets the bridge IDs for the event
func (evt *ContactStatusChange) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetChannelIDs gets the bridge IDs for the event
func (evt *Dial) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Caller.ID)
	if id := evt.Forwarded.ID; id != "" {
		sx = append(sx, id)
	}
	if id := evt.Peer.ID; id != "" {
		sx = append(sx, id)
	}

	return
}

// GetEndpointIDs gets the bridge IDs for the event
func (evt *EndpointStateChange) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetEndpointIDs gets the bridge IDs for the event
func (evt *PeerStatusChange) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackContinuing) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackFinished) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackStarted) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetRecordingIDs gets the recording IDs for the event
func (evt *RecordingFailed) GetRecordingIDs() (sx []string) {
	sx = append(sx, evt.Recording.ID())
	return
}

// GetRecordingIDs gets the recording IDs for the event
func (evt *RecordingFinished) GetRecordingIDs() (sx []string) {
	sx = append(sx, evt.Recording.ID())
	return
}

// GetRecordingIDs gets the recording IDs for the event
func (evt *RecordingStarted) GetRecordingIDs() (sx []string) {
	sx = append(sx, evt.Recording.ID())
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *StasisEnd) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *StasisStart) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Channel.ID)
	if id := evt.ReplaceChannel.ID; id != "" {
		sx = append(sx, id)
	}
	return
}

// GetEndpointIDs gets the bridge IDs for the event
func (evt *TextMessageReceived) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}
