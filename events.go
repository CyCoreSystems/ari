package ari

import (
	"encoding/json"

	"github.com/golang/glog"
)

//Websocket connection for events
//Equivalent to GET /events
func (c *Client) GetEvents(app string) (Message, error) {
	var m Message
	err := c.AriGet("/events/?app="+app, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Generate a new user event
//Equivalent to Post /events/user/{eventName}
func (c *Client) CreateUserEvent(eventName string, req CreateUserEventRequest) error {

	//TODO handle the Error responses individually (by code)

	//Send request
	err := c.AriPost("/events/user/"+eventName, nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Request structure for creating a user event. Only Application is required.
type CreateUserEventRequest struct {
	Application string `json:"application"`
	Source      string `json:"source,omitempty"`
	Variables   string `json:"variables,omitempty"`
}

//
//  Types of events
//

// Event is the base struct for all events
type Event struct {
	Message
	Application string       `json:"application"`
	Timestamp   AsteriskDate `json:"timestamp,omitempty"`
}

// NewEvent constructs an Event from a byte slice (unmarshals the event from Ari)
func NewEvent(raw []byte) (*Event, error) {
	var e Event
	err := json.Unmarshal(raw, &e)
	if err != nil {
		glog.Errorln("Failed to unmarshal new event", err.Error())
		return &e, err
	}

	// Set __raw to be our raw bytestream
	e.__raw = &raw

	return &e, nil
}

// ApplicationReplaced is a notification to the original
// application that notifications about this application
// are now going to a different WebSocket connection.
// There can only be a single websocket connection for each
// application at any given time.
type ApplicationReplaced Event

// BridgeAttendedTransfer events signify that an attended transfer has occurred
type BridgeAttendedTransfer struct {
	Event
	Destination_application      string  `json:"destination_application,omitempty"`
	Destination_bridge           string  `json:"desination_bridge,omitempty"`
	Destination_link_first_leg   Channel `json:"destination_link_first_leg,omitempty"`
	Destination_link_second_leg  Channel `json:"destination_link_second_leg,omitempty"`
	Destination_threeway_bridge  Bridge  `json:"destination_threeway_bridge,omitempty"`
	Destination_threeway_channel Channel `json:"destination_threeway_channel,omitempty"`
	Destination_type             string  `json:"destination_type"`
	Is_external                  bool    `json:"is_external"`
	Result                       string  `json:"result"`
	Transferer_first_leg         Channel `json:"transferer_first_leg"`
	Transferer_first_leg_bridge  Bridge  `json:"transferer_first_leg_bridge,omitempty"`
	Transferer_second_leg        Channel `json:"transferer_second_leg"`
	Transferer_second_leg_bridge Bridge  `json:"transferer_second_leg_bridge,omitempty"`
}

// BridgeBlindTransfer events signify that a blind transfer has occurred
type BridgeBlindTransfer struct {
	Event
	Bridge  Bridge  `json:"bridge,omitempty"`
	Channel Channel `json:"channel"`
	Context string  `json:"context"`
}

// BridgeCreated events indicate a bridge has been created
type BridgeCreated struct {
	Event
	Bridge Bridge `json:"bridge"`
}

// BridgeDestroyed events indicate a bridge has been destroyed
type BridgeDestroyed struct {
	Event
	Bridge Bridge `json:"bridge"`
}

// BridgeMerged events indicate a bridge has been merged into another (bridge)
type BridgeMerged struct {
	Event
	Bridge      Bridge `json:"bridge"`       // New bridge
	Bridge_from Bridge `json: "bridge_from"` // Old (independant) bridge -- TODO: verify this assumption
}

// ChannelCallerId events indicate a channel's caller Id information has changed
type ChannelCallerId struct {
	Event
	Caller_presentation     int     `json:"caller_presentation"`     // The numeric portion
	Caller_presentation_txt string  `json:"caller_presentation_txt"` // The textual portion
	Channel                 Channel `json:"channel"`
}

// ChannelCreated events indicate a channel has been created
type ChannelCreated struct {
	Event
	Channel Channel `json:"channel"`
}

// ChannelDestroyed events indicate a channel has been destroyed
type ChannelDestroyed struct {
	Event
	Cause     int     `json:"cause"`
	Cause_txt string  `json:"cause_txt"`
	Channel   Channel `json:"channel"`
}

// ChannelDialplan events indicate a channel has changed its location in the dialplan
//   NOTE: This event also most likely implies the channel is leaving the control of this
//   ARI application
type ChannelDialplan struct {
	Event
	Channel           Channel `json:"channel"`
	Dialplan_app      string  `json:"dialplan_app"`      // The application which is to be executed
	Dialplan_app_data string  `json:"dialplan_app_data"` // The data to be passed to the application
}

// ChannelDtmfReceived events indicate a channel has received a DTMF tone.
//   NOTE: this event is sent at the _end_ of the DTMF tone (and there is no indication
//   for the _start_ of the DTMF tone)
type ChannelDtmfReceived struct {
	Event
	Channel     Channel `json:"channel"`
	Digit       string  `json:"digit"`       // The DTMF digit which was received (0-9A-E*#)
	Duration_ms int     `json:"duration_ms"` // The duration of the DTMF tone, in milliseconds
}

// ChannelEngteredBridge events indicate that a channel has joined a bridge
type ChannelEnteredBridge struct {
	Event
	Bridge  Bridge  `json:"bridge"`
	Channel Channel `json:"channel"`
}

// ChannelHangupRequest events indicate that a channel has received a hangup
// request
//   TODO: find out if the application is supposed to act on this, or whether
//   this event is purely advisory
type ChannelHangupRequest struct {
	Event
	Cause   int     `json:"cause,omitempty"` // Integer cause code
	Channel Channel `json:"channel"`
	Soft    bool    `json:"soft,omitempty"` // Whether the request was "soft"
}

// ChannelLeftBridge events indicate that a channel has left a bridge
type ChannelLeftBridge struct {
	Event
	Bridge  Bridge  `json:"bridge"`
	Channel Channel `json:"channel"`
}

// ChannelStateChange events indicate that the state of a channel has changed
//   TODO: enumerate the possible channel states
type ChannelStateChange struct {
	Event
	Channel Channel `json:"channel"`
}

// ChannelTalkingFinished events indicate that previously-detected talking on
// a channel is now absent
type ChannelTalkingFinished struct {
	Event
	Channel  Channel `json:"channel"`
	Duration int     `json:"duration"` // Duration (in milliseconds) of the talking
}

// ChannelTalkingStarted events indicate that talking has been detected on a channel
type ChannelTalkingStarted struct {
	Event
	Channel Channel `json:"channel"`
}

// ChannelUserevent events are custom-formatted events which have been received
//   TODO: figure out if these include AMI user events, etc, and how the data is
//   formatted
type ChannelUserevent struct {
	Event
	Bridge    Bridge      `json:"bridge,omitempty"`
	Channel   Channel     `json:"channel,omitempty"`
	Endpoint  Endpoint    `json:"endpoint,omitempty"`
	Eventname string      `json:"eventname"` // Name of the user event
	Userevent interface{} `json:"userevent"` // Custome data sent with the user event
}

// ChannelVarset events indicate a channel variable has been set (or changed)
type ChannelVarset struct {
	Event
	Channel  Channel `json:"channel,omitempty"` // If not present, variable is global
	Value    string  `json:"value"`             // New value
	Variable string  `json:"variable"`          // Variable name
}

// DeviceStateChanged events indicate that a device state has changed
type DeviceStateChanged struct {
	Event
	Device_state DeviceState `json:"device_state"`
}

// Dial events indicate the dialing state for a channel has changed
type Dial struct {
	Event
	Caller     Channel `json:"caller,omitempty"`     // Dialing channel (if not system-originated)
	Dialstatus string  `json:"dialstatus"`           // Present status of the dial attempt
	Dialstring string  `json:"dialstring,omitempty"` // The string describing the dial
	Forward    string  `json:"forward,omitempty"`    // If present, indicates the forwarding target TODO: what is the format of this?
	Forwarded  Channel `json:"forwarded,omitempty"`  // If present, channel of forwarding target
	Peer       Channel `json:"peer"`                 // Dialed channel
}

// EndpointStateChange events indicate that an endpoint has changed state
type EndpointStateChange struct {
	Event
	Endpoint Endpoint `json:"endpoint"`
}

// PlaybackFinished events indicate that a media playback operation has completed
type PlaybackFinished struct {
	Event
	Playback Playback `json:"playback"`
}

// PlaybackStarted events indicate that a media playback operation has begun
type PlaybackStarted struct {
	Event
	Playback Playback `json:"playback"`
}

// RecordingFailed events indicate that a recording operation request has failed to complete
type RecordingFailed struct {
	Event
	Recording LiveRecording `json:"recording"`
}

// RecordingFinished events indicate that a recording operation has completed
type RecordingFinished struct {
	Event
	Recording LiveRecording `json:"recording"`
}

// RecordingStarted events indicate that a recording operation has begun
type RecordingStarted struct {
	Event
	Recording LiveRecording `json:"recording"`
}

// StasisEnd events indicate that a channel has left the Stasis (Ari) application
type StasisEnd struct {
	Event
	Args    []string `json:"args"`
	Channel Channel  `json:"channel"`
}

// StasisStart events indicate a channel has entered the Stasis (Ari) application
type StasisStart struct {
	Event
	Args           []string `json:"args"`
	Channel        Channel  `json:"channel"`
	ReplaceChannel Channel  `json:"replacechannel,omitempty"` // TODO: find out what this is
}

// GetChannel returns the channel for whom the StasisStart event occurred,
// optionally attaching a copy of the ARI client
func (s *StasisStart) GetChannel(c *Client) Channel {
	if c != nil {
		s.Channel.AttachClient(c)
	}
	return s.Channel
}

// TextMessageReceived events indicate that an endpoint has emitted a text message
type TextMessageReceived struct {
	Event
	Endpoint Endpoint    `json:"endpoint,omitempty"`
	Message  TextMessage `json:"message"`
}
