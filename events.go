package ari

import "strings"

// Event is the top level event interface
type Event interface {
	// GetApplication returns the name of the ARI application to which this event is associated
	GetApplication() string

	// GetType returns the type name of this event
	GetType() string

	// Keys returns the related entity keys for the event
	Keys(parent *Key) []*Key
}

// EventData provides the basic metadata for an ARI event
type EventData struct {
	// Type is the type name of this event
	Type string `json:"type"`

	// AsteriskID indicates the unique identifier of the source Asterisk box for this event
	AsteriskID string `json:"asterisk_id,omitempty"`

	// Application indicates the ARI application which emitted this event
	Application string `json:"application"`

	// Timestamp indicates the time this event was generated
	Timestamp DateTime `json:"timestamp,omitempty"`
}

// GetApplication gets the application of the event
func (e *EventData) GetApplication() string {
	return e.Application
}

// GetType gets the type of the event
func (e *EventData) GetType() string {
	return e.Type
}

// implementations of events

// Keys returns the list of keys associated with this event
func (evt *ApplicationReplaced) Keys(parent *Key) (sx []*Key) {
	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeAttendedTransfer) Keys(parent *Key) (sx []*Key) {
	if id := evt.DestinationThreewayBridge.ID; id != "" {
		sx = append(sx, NewKey(BridgeKey, id, WithParent(parent)))
	}
	if id := evt.TransfererFirstLegBridge.ID; id != "" {
		sx = append(sx, NewKey(BridgeKey, id, WithParent(parent)))
	}
	if id := evt.TransfererSecondLegBridge.ID; id != "" {
		sx = append(sx, NewKey(BridgeKey, id, WithParent(parent)))
	}

	if id := evt.DestinationLinkFirstLeg.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.DestinationLinkSecondLeg.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.DestinationThreewayChannel.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.ReplaceChannel.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.Transferee.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.TransfererFirstLeg.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.TransfererSecondLeg.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.TransferTarget.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeBlindTransfer) Keys(parent *Key) (sx []*Key) {

	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))

	if id := evt.ReplaceChannel.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	if id := evt.Transferee.ID; id != "" {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeCreated) Keys(parent *Key) (sx []*Key) {

	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeDestroyed) Keys(parent *Key) (sx []*Key) {

	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeMerged) Keys(parent *Key) (sx []*Key) {

	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	sx = append(sx, NewKey(BridgeKey, evt.BridgeFrom.ID, WithParent(parent)))
	for _, id := range evt.BridgeFrom.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	return
}

// Keys returns the list of keys associated with this event
func (evt *BridgeVideoSourceChanged) Keys(parent *Key) (sx []*Key) {

	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))
	}

	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelCallerID) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelConnectedLine) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelCreated) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelDestroyed) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelDialplan) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelDtmfReceived) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelEnteredBridge) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))

	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelHangupRequest) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelHold) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelLeftBridge) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, id := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, id, WithParent(parent)))

	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelStateChange) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelTalkingFinished) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelTalkingStarted) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelUnhold) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelUserevent) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	sx = append(sx, NewKey(BridgeKey, evt.Bridge.ID, WithParent(parent)))
	for _, i := range evt.Bridge.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, i, WithParent(parent)))
	}
	sx = append(sx, NewEndpointKey(evt.Endpoint.Technology, evt.Endpoint.Resource, WithParent(parent)))
	for _, i := range evt.Endpoint.ChannelIDs {
		sx = append(sx, NewKey(ChannelKey, i, WithParent(parent)))
	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *ChannelVarset) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *ContactInfo) Keys(parent *Key) (sx []*Key) {
	return
}

// Keys returns the list of keys associated with this event
func (evt *ContactStatusChange) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewEndpointKey(evt.Endpoint.Technology, evt.Endpoint.Resource, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *DeviceStateChanged) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(DeviceStateKey, string(evt.DeviceState), WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *Dial) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Caller.ID, WithParent(parent)))
	sx = append(sx, NewKey(ChannelKey, evt.Peer.ID, WithParent(parent)))
	sx = append(sx, NewKey(ChannelKey, evt.Forwarded.ID, WithParent(parent)))

	return
}

// Keys returns the list of keys associated with this event
func (evt *EndpointStateChange) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewEndpointKey(evt.Endpoint.Technology, evt.Endpoint.Resource, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *MissingParams) Keys(parent *Key) (sx []*Key) {
	return
}

// Keys returns the list of keys associated with this event
func (evt *Peer) Keys(parent *Key) (sx []*Key) {
	return
}

// Keys returns the list of keys associated with this event
func (evt *PeerStatusChange) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewEndpointKey(evt.Endpoint.Technology, evt.Endpoint.Resource, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *PlaybackContinuing) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(PlaybackKey, evt.Playback.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *PlaybackFinished) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(PlaybackKey, evt.Playback.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *PlaybackStarted) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(PlaybackKey, evt.Playback.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *RecordingFailed) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(LiveRecordingKey, evt.Recording.Name, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *RecordingFinished) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(LiveRecordingKey, evt.Recording.Name, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *RecordingStarted) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(LiveRecordingKey, evt.Recording.Name, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *StasisEnd) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	return
}

// Keys returns the list of keys associated with this event
func (evt *StasisStart) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewKey(ChannelKey, evt.Channel.ID, WithParent(parent)))
	if evt.ReplaceChannel.ID != "" {
		sx = append(sx, NewKey(ChannelKey, evt.ReplaceChannel.ID, WithParent(parent)))
	}
	return
}

// Keys returns the list of keys associated with this event
func (evt *TextMessageReceived) Keys(parent *Key) (sx []*Key) {
	sx = append(sx, NewEndpointKey(evt.Endpoint.Technology, evt.Endpoint.Resource, WithParent(parent)))
	return
}

// ----------

// Created marks the BridgeCreated event that it created an event
func (evt *BridgeCreated) Created() (bridgeID string, related string) {
	bridgeID = evt.Bridge.ID
	if len(evt.Bridge.ChannelIDs) != 0 {
		related = evt.Bridge.ChannelIDs[0]
	} else {
		related = evt.Bridge.Creator
	}
	return
}

// Destroyed returns the bridge that was finished by this event.
// Used by the proxy to route events to dialogs.
func (evt *BridgeDestroyed) Destroyed() string {
	return evt.Bridge.ID
}

// GetChannelIDs gets the channel IDs for the event
func (evt *BridgeCreated) GetChannelIDs() (sx []string) {
	sx = evt.Bridge.ChannelIDs
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

// Created marks the event as creating a bridge for a channel and dialog
func (evt *ChannelEnteredBridge) Created() (o string, related string) {
	o = evt.Bridge.ID
	related = evt.Channel.ID
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

// GetEndpointIDs gets the endpoint IDs for the event
func (evt *EndpointStateChange) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *EndpointStateChange) GetChannelIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ChannelIDs...)
	return
}

// GetEndpointIDs gets the endpoint IDs for the event
func (evt *PeerStatusChange) GetEndpointIDs() (sx []string) {
	sx = append(sx, evt.Endpoint.ID())
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackContinuing) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *PlaybackContinuing) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Playback.TargetURI)
	if s == "" {
		return
	}

	sx = append(sx, s)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *PlaybackContinuing) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Playback.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackFinished) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *PlaybackFinished) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Playback.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *PlaybackFinished) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Playback.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// Destroyed returns the playbacK ID that was finished by this event.
// Used by the proxy to route events to dialogs.
func (evt *PlaybackFinished) Destroyed() (playbackID string) {
	playbackID = evt.Playback.ID
	return
}

// GetPlaybackIDs gets the playback IDs for the event
func (evt *PlaybackStarted) GetPlaybackIDs() (sx []string) {
	sx = append(sx, evt.Playback.ID)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *PlaybackStarted) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Playback.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *PlaybackStarted) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Playback.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// Created returns the playbacK ID that we created plus the ID that the playback
// is operating on (a bridge or channel).
// Used by the proxy to route events to dialogs
func (evt *PlaybackStarted) Created() (playbackID, otherID string) {
	playbackID = evt.Playback.ID
	items := strings.Split(evt.Playback.TargetURI, ":")
	if len(items) == 1 {
		otherID = items[0]
	} else {
		otherID = items[1]
	}

	return
}

// Destroyed returns the item that gets destroyed by this event
func (evt *RecordingFailed) Destroyed() string {
	return evt.Recording.ID()
}

// GetRecordingIDs gets the recording IDs for the event
func (evt *RecordingFailed) GetRecordingIDs() (sx []string) {
	sx = append(sx, evt.Recording.ID())
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *RecordingFailed) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *RecordingStarted) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetChannelIDs gets the channel IDs for the event
func (evt *RecordingFinished) GetChannelIDs() (sx []string) {
	s := resolveTarget("channel", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *RecordingFailed) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *RecordingStarted) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetBridgeIDs gets the bridge IDs for the event
func (evt *RecordingFinished) GetBridgeIDs() (sx []string) {
	s := resolveTarget("bridge", evt.Recording.TargetURI)
	if s == "" {
		return
	}
	sx = append(sx, s)
	return
}

// GetRecordingIDs gets the recording IDs for the event
func (evt *RecordingFinished) GetRecordingIDs() (sx []string) {
	sx = append(sx, evt.Recording.ID())
	return
}

// Destroyed returns the item that gets destroyed by this event
func (evt *RecordingFinished) Destroyed() string {
	return evt.Recording.ID()
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

func resolveTarget(typ string, targetURI string) (s string) {
	items := strings.Split(targetURI, ":")
	if items[0] != typ {
		return
	}
	if len(items) < 2 {
		return
	}

	s = strings.Join(items[1:], ":")
	return

}

// Header represents a set of key-value pairs to store transport-related metadata on Events
type Header map[string][]string

// Add appens the value to the list of values for the given header key.
func (h Header) Add(key, val string) {
	h[key] = append(h[key], val)
}

// Set sets the value for the given header key, replacing any existing values.
func (h Header) Set(key, val string) {
	h[key] = []string{val}
}

// Get returns the first value associated with the given header key.
func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}

	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Del deletes the values associated with the given header key.
func (h Header) Del(key string) {
	delete(h, key)
}
