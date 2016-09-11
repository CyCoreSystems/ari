package ari

type Event interface {

	// TODO(sea):  Delete Get(), Data(), add a Create() function
	// Get returns a handle pointer to the Event for further interaction
	Get(name string) *EventHandle

	// Data returns the Event's data
	Data(name string) (EventData, error)
}

/* --------- Event ---------------- */

// EventHandle provides a wrapper to an Event interface for
// operations on a specific Event
type EventHandle struct {
	name string
	e    Event
}

// Event is the base struct for all events
type EventData struct {
	Message
	Application string   `json:"application"`
	Timestamp   DateTime `json:"timestamp,omitempty"`
}

// NewEventHandle creates a new handle to the Event name
func NewEventHandle(name string, e Event) *EventHandle {
	return &EventHandle{
		name: name,
		e:    e,
	}
}

// Data retrieves the data for the Event
func (eh *EventHandle) Data() (ed EventData, err error) {
	ed, err = eh.e.Data(eh.name)
	return ed, err
}

/* --------- UserEvent ---------------- */

//Request structure for creating a user event. Only Application is required.
type CreateUserEventRequest struct {
	Application string `json:"application"`
	Source      string `json:"source,omitempty"`
	Variables   string `json:"variables,omitempty"`
}

/* --------- Meta Events ---------------- */

//
// Meta events
// These events are intermediate event types used for further
// classification of events
//

// BridgeEvent (meta) is an event which affects a bridge
type BridgeEvent struct {
	EventData
	BridgeData
}

// ChannelEvent (meta) is an event which affects a channel
type ChannelEvent struct {
	EventData
	ChannelData
}

// ApplicationReplaced is a notification to the original
// application that notifications about this application
// are now going to a different WebSocket connection.
// There can only be a single websocket connection for each
// application at any given time.
type ApplicationReplaced Event

/* --------- Bridge Events ---------------- */

// BridgeAttendedTransfer events signify that an attended transfer has occurred
type BridgeAttendedTransfer struct {
	EventData
	DestinationApplication     string  `json:"destination_application,omitempty"`
	DestinationBridge          string  `json:"desination_bridge,omitempty"`
	DestinationLinkFirstLeg    Channel `json:"destination_link_first_leg,omitempty"`
	DestinationLinkSecondLeg   Channel `json:"destination_link_second_leg,omitempty"`
	DestinationThreewayBridge  Bridge  `json:"destination_threeway_bridge,omitempty"`
	DestinationThreewayChannel Channel `json:"destination_threeway_channel,omitempty"`
	DestinationType            string  `json:"destination_type"`
	External                   bool    `json:"is_external"`
	Result                     string  `json:"result"`
	Transferee                 Channel `json:"transferee,omitempty"`
	TransferTarget             Channel `json:"transfer_target,omitempty"`
	TransfererFirstLeg         Channel `json:"transferer_first_leg"`
	TransfererFirstLegBridge   Bridge  `json:"transferer_first_leg_bridge,omitempty"`
	TransfererSecondLeg        Channel `json:"transferer_second_leg"`
	TransfererSecondLegBridge  Bridge  `json:"transferer_second_leg_bridge,omitempty"`
}

// BridgeBlindTransfer events signify that a blind transfer has occurred
type BridgeBlindTransfer struct {
	EventData
	BridgeData
	ChannelData
	Context        string  `json:"context"`
	Extension      string  `json:"exten"`
	External       bool    `json:"external"`
	ReplaceChannel Channel `json:"external,omitempty"`
	Result         string  `json:"result"`
	Transferee     Channel `json:"transferee,omitempty"`
}

// BridgeCreated events indicate a bridge has been created
type BridgeCreated struct {
	EventData
	BridgeData
}

// BridgeDestroyed events indicate a bridge has been destroyed
type BridgeDestroyed struct {
	EventData
	BridgeData
}

// BridgeMerged events indicate a bridge has been merged into another (bridge)
type BridgeMerged struct {
	EventData
	BridgeData
	Bridge_from Bridge `json:"bridge_from"` // Old (independant) bridge -- TODO: verify this assumption
}

/* --------- Channel Events ---------------- */

// ChannelCallerId events indicate a channel's caller Id information has changed
type ChannelCallerId struct {
	EventData
	Caller_presentation     int    `json:"caller_presentation"`     // The numeric portion
	Caller_presentation_txt string `json:"caller_presentation_txt"` // The textual portion
	ChannelData
}

// ChannelConnectedLine events indicate a channel has changed its connected line
type ChannelConnectedLine struct {
	EventData
	ChannelData
}

// ChannelCreated events indicate a channel has been created
type ChannelCreated struct {
	EventData
	ChannelData
}

// ChannelDestroyed events indicate a channel has been destroyed
type ChannelDestroyed struct {
	EventData
	Cause     int    `json:"cause"`
	Cause_txt string `json:"cause_txt"`
	ChannelData
}

// ChannelDialplan events indicate a channel has changed its location in the dialplan
//   NOTE: This event also most likely implies the channel is leaving the control of this
//   ARI application
type ChannelDialplan struct {
	EventData
	ChannelData
	Dialplan_app      string `json:"dialplan_app"`      // The application which is to be executed
	Dialplan_app_data string `json:"dialplan_app_data"` // The data to be passed to the application
}

// ChannelDtmfReceived events indicate a channel has received a DTMF tone.
//   NOTE: this event is sent at the _end_ of the DTMF tone (and there is no indication
//   for the _start_ of the DTMF tone)
type ChannelDtmfReceived struct {
	EventData
	ChannelData
	Digit       string `json:"digit"`       // The DTMF digit which was received (0-9A-E*#)
	Duration_ms int    `json:"duration_ms"` // The duration of the DTMF tone, in milliseconds
}

// ChannelEngteredBridge events indicate that a channel has joined a bridge
type ChannelEnteredBridge struct {
	EventData
	BridgeData
	ChannelData
}

// ChannelHangupRequest events indicate that a channel has received a hangup
// request
//   TODO: find out if the application is supposed to act on this, or whether
//   this event is purely advisory
type ChannelHangupRequest struct {
	EventData
	Cause int `json:"cause,omitempty"` // Integer cause code
	ChannelData
	Soft bool `json:"soft,omitempty"` // Whether the request was "soft"
}

// ChannelHold events indicate that a channel has been put on hold
type ChannelHold struct {
	EventData
	ChannelData
	MusicClass string `json:"musicclass"` // The music class being played to the held channel
}

// ChannelLeftBridge events indicate that a channel has left a bridge
type ChannelLeftBridge struct {
	EventData
	BridgeData
	ChannelData
}

// ChannelStateChange events indicate that the state of a channel has changed
//   TODO: enumerate the possible channel states
type ChannelStateChange struct {
	EventData
	ChannelData
}

// ChannelTalkingFinished events indicate that previously-detected talking on
// a channel is now absent
type ChannelTalkingFinished struct {
	EventData
	ChannelData
	Duration int `json:"duration"` // Duration (in milliseconds) of the talking
}

// ChannelTalkingStarted events indicate that talking has been detected on a channel
type ChannelTalkingStarted struct {
	EventData
	ChannelData
}

// ChannelUnhold events indicate that a channel has removed from Hold
type ChannelUnhold struct {
	EventData
	ChannelData
}

// ChannelUserevent events are custom-formatted events which have been received
//   TODO: figure out if these include AMI user events, etc, and how the data is
//   formatted
type ChannelUserevent struct {
	EventData
	BridgeData
	ChannelData
	EndpointData
	Eventname string      `json:"eventname"` // Name of the user event
	Userevent interface{} `json:"userevent"` // Custome data sent with the user event
}

// ChannelVarset events indicate a channel variable has been set (or changed)
type ChannelVarset struct {
	EventData
	ChannelData
	Value    string `json:"value"`    // New value
	Variable string `json:"variable"` // Variable name
}

/* --------- Contact Events ---------------- */

// ContactStatusChanged events indicate that a contact state has changed
type ContactStatusChanged struct {
	EventData
	EndpointData
	ContactInfo ContactInfo `json:"contact_info"`
}

// DeviceStateChanged events indicate that a device state has changed
type DeviceStateChanged struct {
	EventData
	Device_state DeviceState `json:"device_state"`
}

// Dial events indicate the dialing state for a channel has changed
type Dial struct {
	EventData
	Caller     Channel `json:"caller,omitempty"`     // Dialing channel (if not system-originated)
	Dialstatus string  `json:"dialstatus"`           // Present status of the dial attempt
	Dialstring string  `json:"dialstring,omitempty"` // The string describing the dial
	Forward    string  `json:"forward,omitempty"`    // If present, indicates the forwarding target TODO: what is the format of this?
	Forwarded  Channel `json:"forwarded,omitempty"`  // If present, channel of forwarding target
	Peer       Channel `json:"peer"`                 // Dialed channel
}

// EndpointStateChange events indicate that an endpoint has changed state
type EndpointStateChange struct {
	EventData
	EndpointData
}

/* --------- Peer Events ---------------- */

// PeerStatusChange events indicate that a peer has changed state
type PeerStatusChange struct {
	EventData
	EndpointData
	Peer     Peer     `json:"peer"`
}

/* --------- Playback Events ---------------- */

// PlaybackFinished events indicate that a media playback operation has completed
type PlaybackFinished struct {
	EventData
	Playback Playback `json:"playback"`
}

// PlaybackStarted events indicate that a media playback operation has begun
type PlaybackStarted struct {
	EventData
	Playback Playback `json:"playback"`
}

/* --------- Recording Events ---------------- */


// RecordingFailed events indicate that a recording operation request has failed to complete
type RecordingFailed struct {
	EventData
	Recording LiveRecordingData
}

// RecordingFinished events indicate that a recording operation has completed
type RecordingFinished struct {
	EventData
	Recording LiveRecordingData
}

// RecordingStarted events indicate that a recording operation has begun
type RecordingStarted struct {
	EventData
	Recording LiveRecordingData
}

/* --------- Stasis Events ---------------- */

// StasisEnd events indicate that a channel has left the Stasis (Ari) application
type StasisEnd struct {
	EventData
	Args []string `json:"args"`
	ChannelData
}

// StasisStart events indicate a channel has entered the Stasis (Ari) application
type StasisStart struct {
	EventData
	Args []string `json:"args"`
	ChannelData
	ReplaceChannel Channel `json:"replacechannel,omitempty"` // TODO: find out what this is
}

// R2 - SKIPPED FOR NOW
// GetChannel returns the channel for whom the StasisStart event occurred,
// optionally attaching a copy of the ARI client
/*
func (s *StasisStart) GetChannel(c *Client) Channel {
	if c != nil {
		s.Channel.AttachClient(c)
	}
return s.Channel
}
*/

/* --------- TextMessage Events ---------------- */

// TextMessageReceived events indicate that an endpoint has emitted a text message
type TextMessageReceived struct {
	EventData
	EndpointData
	Message  TextMessage `json:"message"`
}

