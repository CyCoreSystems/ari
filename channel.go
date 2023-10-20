package ari

import (
	"encoding/json"
	"strings"
	"time"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/rotisserie/eris"
)

// Channel represents a communication path interacting with an Asterisk server.
type Channel interface {
	// Get returns a handle to a channel for further interaction
	Get(key *Key) *ChannelHandle

	// GetVariable retrieves the value of a channel variable
	GetVariable(*Key, string) (string, error)

	// List lists the channels in asterisk, optionally using the key for filtering
	List(*Key) ([]*Key, error)

	// Originate creates a new channel, returning a handle to it or an error, if
	// the creation failed.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	Originate(*Key, OriginateRequest) (*ChannelHandle, error)

	// StageOriginate creates a new Originate, created when the `Exec` method
	// on `ChannelHandle` is invoked.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	StageOriginate(*Key, OriginateRequest) (*ChannelHandle, error)

	// Create creates a new channel, returning a handle to it or an
	// error, if the creation failed. Create is already Staged via `Dial`.
	// The Key should be that of the linked channel, if one exists, so that the
	// Node can be matches to it.
	Create(*Key, ChannelCreateRequest) (*ChannelHandle, error)

	// Data returns the channel data for a given channel
	Data(key *Key) (*ChannelData, error)

	// Continue tells Asterisk to return a channel to the dialplan
	Continue(key *Key, context, extension string, priority int) error

	// Move moves the channel into another Stasis application
	Move(key *Key, app string, appArgs string) error

	// Busy hangs up the channel with the "busy" cause code
	Busy(key *Key) error

	// Congestion hangs up the channel with the "congestion" cause code
	Congestion(key *Key) error

	// Answer answers the channel
	Answer(key *Key) error

	// Hangup hangs up the given channel
	Hangup(key *Key, reason string) error

	// Ring indicates ringing to the channel
	Ring(key *Key) error

	// StopRing stops ringing on the channel
	StopRing(key *Key) error

	// SendDTMF sends DTMF to the channel
	SendDTMF(key *Key, dtmf string, opts *DTMFOptions) error

	// Hold puts the channel on hold
	Hold(key *Key) error

	// StopHold retrieves the channel from hold
	StopHold(key *Key) error

	// Mute mutes a channel in the given direction (in,out,both)
	Mute(key *Key, dir Direction) error

	// Unmute unmutes a channel in the given direction (in,out,both)
	Unmute(key *Key, dir Direction) error

	// MOH plays music on hold
	MOH(key *Key, moh string) error

	// SetVariable sets a channel variable
	SetVariable(key *Key, name, value string) error

	// StopMOH stops music on hold
	StopMOH(key *Key) error

	// Silence plays silence to the channel
	Silence(key *Key) error

	// StopSilence stops the silence on the channel
	StopSilence(key *Key) error

	// Play plays the media URI to the channel
	Play(key *Key, playbackID string, mediaURI ...string) (*PlaybackHandle, error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	StagePlay(key *Key, playbackID string, mediaURI ...string) (*PlaybackHandle, error)

	// Record records the channel
	Record(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	StageRecord(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// Dial dials a created channel
	Dial(key *Key, caller string, timeout time.Duration) error

	// Snoop spies on a specific channel, creating a new snooping channel
	Snoop(key *Key, snoopID string, opts *SnoopOptions) (*ChannelHandle, error)

	// StageSnoop creates a new `ChannelHandle`, when `Exec`ed, snoops on the given channel ID and
	// creates a new snooping channel.
	StageSnoop(key *Key, snoopID string, opts *SnoopOptions) (*ChannelHandle, error)

	// StageExternalMedia creates a new non-telephony external media channel,
	// when `Exec`ed, by which audio may be sent or received.  The stage version
	// of this command will not actually communicate with Asterisk until Exec is
	// called on the returned ExternalMedia channel.
	StageExternalMedia(key *Key, opts ExternalMediaOptions) (*ChannelHandle, error)

	// ExternalMedia creates a new non-telephony external media channel by which audio may be sent or received
	ExternalMedia(key *Key, opts ExternalMediaOptions) (*ChannelHandle, error)

	// Subscribe subscribes on the channel events
	Subscribe(key *Key, n ...string) Subscription

	// UserEvent Sends user-event to AMI channel subscribers
	UserEvent(key *Key, ue *ChannelUserevent) error

	// Redirect tells Asterisk to redirect the channel to a different location
	Redirect(key *Key, endpoint string) error
}

// channelDataJSON is the data for a specific channel
type channelDataJSON struct {
	// Key is the unique identifier for a channel in the cluster
	Key *Key `json:"key,omitempty"`

	ID           string            `json:"id"`    // Unique id for this channel (same as for AMI)
	Name         string            `json:"name"`  // Name of this channel (tech/name-id format)
	State        string            `json:"state"` // State of the channel
	Accountcode  string            `json:"accountcode"`
	Caller       *CallerID         `json:"caller"`    // CallerId of the calling endpoint
	Connected    *CallerID         `json:"connected"` // CallerId of the connected line
	Creationtime DateTime          `json:"creationtime"`
	Dialplan     *DialplanCEP      `json:"dialplan"` // Current location in the dialplan
	ChannelVars  map[string]string `json:"channelvars"`
}

// MarshalJSON encodes ChannelData to JSON
func (d *ChannelData) MarshalJSON() ([]byte, error) {
	t, err := ptypes.TimestampFromProto(d.Creationtime)
	if err != nil {
		t = time.Now()
	}

	return json.Marshal(&channelDataJSON{
		Key:          d.Key,
		ID:           d.ID,
		Name:         d.Name,
		State:        d.State,
		Accountcode:  d.Accountcode,
		Caller:       d.Caller,
		Connected:    d.Connected,
		Creationtime: DateTime(t),
		Dialplan:     d.Dialplan,
		ChannelVars:  d.ChannelVars,
	})
}

// UnmarshalJSON decodes ChannelData from JSON
func (d *ChannelData) UnmarshalJSON(data []byte) error {
	in := new(channelDataJSON)

	if err := json.Unmarshal(data, in); err != nil {
		return err
	}

	t, err := ptypes.TimestampProto(time.Time(in.Creationtime))
	if err != nil {
		t = &ptypes.Timestamp{
			Seconds: time.Now().Unix(),
		}
	}

	*d = ChannelData{
		Key:          in.Key,
		ID:           in.ID,
		Name:         in.Name,
		State:        in.State,
		Accountcode:  in.Accountcode,
		Caller:       in.Caller,
		Connected:    in.Connected,
		Creationtime: t,
		Dialplan:     in.Dialplan,
		ChannelVars:  in.ChannelVars,
	}

	return nil
}

// ChannelCreateRequest describes how a channel should be created, when
// using the separate Create and Dial calls.
type ChannelCreateRequest struct {
	// Endpoint is the target endpoint for the dial
	Endpoint string `json:"endpoint"`

	// App is the name of the Stasis application to execute on connection
	App string `json:"app"`

	// AppArgs is the set of (comma-separated) arguments for the Stasis App
	AppArgs string `json:"appArgs,omitempty"`

	// ChannelID is the ID to give to the newly-created channel
	ChannelID string `json:"channelId,omitempty"`

	// OtherChannelID is the ID of the second created channel (when creating Local channels)
	OtherChannelID string `json:"otherChannelId,omitempty"`

	// Originator is the unique ID of the calling channel, for which this new channel-dial is being created
	Originator string `json:"originator,omitempty"`

	// Formats is the comma-separated list of valid codecs to allow for the new channel, in the case that
	// the Originator is not specified
	Formats string `json:"formats,omitempty"`
}

// SnoopOptions enumerates the non-required arguments for the snoop operation
type SnoopOptions struct {
	// App is the ARI application into which the newly-created Snoop channel should be dropped.
	App string `json:"app"`

	// AppArgs is the set of arguments to pass with the newly-created Snoop channel's entry into ARI.
	AppArgs string `json:"appArgs,omitempty"`

	// Spy describes the direction of audio on which to spy (none, in, out, both).
	// The default is 'none'.
	Spy Direction `json:"spy,omitempty"`

	// Whisper describes the direction of audio on which to send (none, in, out, both).
	// The default is 'none'.
	Whisper Direction `json:"whisper,omitempty"`
}

// ExternalMediaOptions describes the parameters to the externalMedia channel creation operation
type ExternalMediaOptions struct {
	// ChannelID specifies the channel ID to be used for the external media channel.  This parameter is optional and if not specified, a randomly-generated channel ID will be used.
	ChannelID string `json:"channelId"`

	// App is the ARI Application to which the newly-created external media channel should be placed.  This parameter is optional and if not specified, the current application will be used.
	App string `json:"app"`

	// ExternalHost specifies the <host>:<port> of the external host to which the external media channel will be connected.  This parameter is MANDATORY and has no default.
	ExternalHost string `json:"external_host"`

	// Encapsulation specifies the payload encapsulation which should be used.  Options include:  'rtp'.  This parameter is optional and if not specified, 'rtp' will be used.
	Encapsulation string `json:"encapsulation"`

	// Transport specifies the connection type to be used to communicate to the external server.  Options include 'udp'.  This parameter is optional and if not specified, 'udp' will be used.
	Transport string `json:"transport"`

	// ConnectionType defined the directionality of the network connection.  Options include 'client' and 'server'.  This parameter is optional and if not specified, 'client' will be used.
	ConnectionType string `json:"connection_type"`

	// Format specifies the codec to be used for the audio.  Options include 'slin16', 'ulaw' (and likely other codecs supported by Asterisk).  This parameter is MANDATORY and has not default.
	Format string `json:"format"`

	// Direction specifies the directionality of the audio stream.  Options include 'both'.  This parameter is optional and if not specified, 'both' will be used.
	Direction string `json:"direction"`

	// Data: when encapsulation=audiosocket this specifies the UUID to send
	Data string `json:"data,omitempty"`

	// Variables defines the set of channel variables which should be bound to this channel upon creation.  This parameter is optional.
	Variables map[string]string `json:"variables"`
}

// ChannelHandle provides a wrapper on the Channel interface for operations on a particular channel ID.
type ChannelHandle struct {
	key *Key
	c   Channel

	exec func(ch *ChannelHandle) error

	executed bool
}

// NewChannelHandle returns a handle to the given ARI channel
func NewChannelHandle(key *Key, c Channel, exec func(ch *ChannelHandle) error) *ChannelHandle {
	return &ChannelHandle{
		key:  key,
		c:    c,
		exec: exec,
	}
}

// ID returns the identifier for the channel handle
func (ch *ChannelHandle) ID() string {
	return ch.key.ID
}

// Key returns the key for the channel handle
func (ch *ChannelHandle) Key() *Key {
	return ch.key
}

// Exec executes any staged channel operations attached to this handle.
func (ch *ChannelHandle) Exec() (err error) {
	if !ch.executed {
		ch.executed = true
		if ch.exec != nil {
			err = ch.exec(ch)
			ch.exec = nil
		}
	}

	return err
}

// Data returns the channel's data
func (ch *ChannelHandle) Data() (*ChannelData, error) {
	return ch.c.Data(ch.key)
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension string, priority int) error {
	return ch.c.Continue(ch.key, context, extension, priority)
}

// Move moves the channel to a new Stasis app
func (ch *ChannelHandle) Move(app string, appArgs string) error {
	return ch.c.Move(ch.key, app, appArgs)
}

//---
// Play/Record operations
//---

// Play initiates playback of the specified media uri
// to the channel, returning the Playback handle
func (ch *ChannelHandle) Play(id string, mediaURI ...string) (ph *PlaybackHandle, err error) {
	return ch.c.Play(ch.key, id, mediaURI...)
}

// Record records the channel to the given filename
func (ch *ChannelHandle) Record(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
	return ch.c.Record(ch.key, name, opts)
}

// StagePlay stages a `Play` operation.
func (ch *ChannelHandle) StagePlay(id string, mediaURI ...string) (*PlaybackHandle, error) {
	return ch.c.StagePlay(ch.key, id, mediaURI...)
}

// StageRecord stages a `Record` operation
func (ch *ChannelHandle) StageRecord(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
	return ch.c.StageRecord(ch.key, name, opts)
}

//---
// Hangup Operations
//---

// Busy hangs up the channel with the "busy" cause code
func (ch *ChannelHandle) Busy() error {
	return ch.c.Busy(ch.key)
}

// Congestion hangs up the channel with the congestion cause code
func (ch *ChannelHandle) Congestion() error {
	return ch.c.Congestion(ch.key)
}

// Hangup hangs up the channel with the normal cause code
func (ch *ChannelHandle) Hangup() error {
	return ch.c.Hangup(ch.key, "normal")
}

//--

// --
// Answer operations
// --

// Answer answers the channel
func (ch *ChannelHandle) Answer() error {
	return ch.c.Answer(ch.key)
}

// IsAnswered checks the current state of the channel to see if it is "Up"
func (ch *ChannelHandle) IsAnswered() (bool, error) {
	updated, err := ch.Data()
	if err != nil {
		return false, eris.Wrap(err, "Failed to get updated channel")
	}

	return strings.ToLower(updated.State) == "up", nil
}

// ------

// --
// Ring Operations
// --

// Ring indicates ringing to the channel
func (ch *ChannelHandle) Ring() error {
	return ch.c.Ring(ch.key)
}

// StopRing stops ringing on the channel
func (ch *ChannelHandle) StopRing() error {
	return ch.c.StopRing(ch.key)
}

// ------

// --
// Mute operations
// --

// Mute mutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Mute(dir Direction) (err error) {
	if dir == "" {
		dir = DirectionIn
	}

	return ch.c.Mute(ch.key, dir)
}

// Unmute unmutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Unmute(dir Direction) (err error) {
	if dir == "" {
		dir = DirectionIn
	}

	return ch.c.Unmute(ch.key, dir)
}

// ----

// --
// Hold operations
// --

// Hold puts the channel on hold
func (ch *ChannelHandle) Hold() error {
	return ch.c.Hold(ch.key)
}

// StopHold retrieves the channel from hold
func (ch *ChannelHandle) StopHold() error {
	return ch.c.StopHold(ch.key)
}

// ----

// --
// Music on hold operations
// --

// MOH plays music on hold of the given class
// to the channel
func (ch *ChannelHandle) MOH(mohClass string) error {
	return ch.c.MOH(ch.key, mohClass)
}

// StopMOH stops playing of music on hold to the channel
func (ch *ChannelHandle) StopMOH() error {
	return ch.c.StopMOH(ch.key)
}

// ----

// GetVariable returns the value of a channel variable
func (ch *ChannelHandle) GetVariable(name string) (string, error) {
	return ch.c.GetVariable(ch.key, name)
}

// SetVariable sets the value of a channel variable
func (ch *ChannelHandle) SetVariable(name, value string) error {
	return ch.c.SetVariable(ch.key, name, value)
}

// --
// Misc
// --

// Originate creates (and dials) a new channel using the present channel as its Originator.
func (ch *ChannelHandle) Originate(req OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return ch.c.Originate(ch.key, req)
}

// StageOriginate stages an originate (channel creation and dial) to be Executed later.
func (ch *ChannelHandle) StageOriginate(req OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return ch.c.StageOriginate(ch.key, req)
}

// Create creates (but does not dial) a new channel, using the present channel as its Originator.
func (ch *ChannelHandle) Create(req ChannelCreateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}

	return ch.c.Create(ch.key, req)
}

// Dial dials a created channel.  `caller` is the optional
// channel ID of the calling party (if there is one).  Timeout
// is the length of time to wait before the dial is answered
// before aborting.
func (ch *ChannelHandle) Dial(caller string, timeout time.Duration) error {
	return ch.c.Dial(ch.key, caller, timeout)
}

// Snoop spies on a specific channel, creating a new snooping channel placed into the given app
func (ch *ChannelHandle) Snoop(snoopID string, opts *SnoopOptions) (*ChannelHandle, error) {
	return ch.c.Snoop(ch.key, snoopID, opts)
}

// StageSnoop stages a `Snoop` operation
func (ch *ChannelHandle) StageSnoop(snoopID string, opts *SnoopOptions) (*ChannelHandle, error) {
	return ch.c.StageSnoop(ch.key, snoopID, opts)
}

// StageExternalMedia creates a new non-telephony external media channel,
// when `Exec`ed, by which audio may be sent or received.  The stage version
// of this command will not actually communicate with Asterisk until Exec is
// called on the returned ExternalMedia channel.
func (ch *ChannelHandle) StageExternalMedia(opts ExternalMediaOptions) (*ChannelHandle, error) {
	return ch.c.StageExternalMedia(ch.key, opts)
}

// ExternalMedia creates a new non-telephony external media channel by which audio may be sent or received
func (ch *ChannelHandle) ExternalMedia(opts ExternalMediaOptions) (*ChannelHandle, error) {
	return ch.c.ExternalMedia(ch.key, opts)
}

// ----

// --
// Silence operations
// --

// Silence plays silence to the channel
func (ch *ChannelHandle) Silence() error {
	return ch.c.Silence(ch.key)
}

// StopSilence stops silence to the channel
func (ch *ChannelHandle) StopSilence() error {
	return ch.c.StopSilence(ch.key)
}

// ----

// --
// Subscription
// --

// Subscribe subscribes the list of channel events
func (ch *ChannelHandle) Subscribe(n ...string) Subscription {
	if ch == nil {
		return nil
	}

	return ch.c.Subscribe(ch.key, n...)
}

// TODO: rest of ChannelHandle

// --
// DTMF
// --

// SendDTMF sends the DTMF information to the server
func (ch *ChannelHandle) SendDTMF(dtmf string, opts *DTMFOptions) error {
	return ch.c.SendDTMF(ch.key, dtmf, opts)
}

// UserEvent sends user-event to AMI channel subscribers
func (ch *ChannelHandle) UserEvent(key *Key, ue *ChannelUserevent) error {
	return ch.c.UserEvent(ch.key, ue)
}

// Redirect tells Asterisk to redirect the channel to a different location
func (ch *ChannelHandle) Redirect(endpoint string) error {
	return ch.c.Redirect(ch.key, endpoint)
}
