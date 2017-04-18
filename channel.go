package ari

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Channel represents a communication path interacting with an Asterisk server.
type Channel interface {
	// Get returns a handle to a channel for further interaction
	Get(key *Key) *ChannelHandle

	// List lists the channels in asterisk, optionally using the key for filtering
	List(*Key) ([]*Key, error)

	// Originate creates a new channel, returning a handle to it or an
	// error, if the creation failed
	Originate(OriginateRequest) (*ChannelHandle, error)

	// StageOriginate creates a new Originate, created when the `Exec` method
	// on `ChannelHandle` is invoked
	StageOriginate(OriginateRequest) *ChannelHandle

	// Create creates a new channel, returning a handle to it or an
	// error, if the creation failed. Create is already Staged via `Dial`.
	Create(ChannelCreateRequest) (*ChannelHandle, error)

	// Data returns the channel data for a given channel
	Data(key *Key) (*ChannelData, error)

	// Continue tells Asterisk to return a channel to the dialplan
	Continue(key *Key, context, extension string, priority int) error

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

	// StopMOH stops music on hold
	StopMOH(key *Key) error

	// Silence plays silence to the channel
	Silence(key *Key) error

	// StopSilence stops the silence on the channel
	StopSilence(key *Key) error

	// Play plays the media URI to the channel
	Play(key *Key, playbackID string, mediaURI string) (PlaybackHandle, error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	StagePlay(key *Key, playbackID string, mediaURI string) (ph PlaybackHandle)

	// Record records the channel
	Record(key *Key, name string, opts *RecordingOptions) (LiveRecordingHandle, error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	StageRecord(key *Key, name string, opts *RecordingOptions) (rh LiveRecordingHandle)

	// Dial dials a created channel
	Dial(key *Key, caller string, timeout time.Duration) error

	// Snoop spies on a specific channel, creating a new snooping channel
	Snoop(key *Key, snoopID string, opts *SnoopOptions) (*ChannelHandle, error)

	// StageSnoop creates a new `ChannelHandle`, when `Exec`ed, snoops on the given channel ID and
	// creates a new snooping channel.
	StageSnoop(key *Key, snoopID string, opts *SnoopOptions) *ChannelHandle

	// Subscribe subscribes on the channel events
	Subscribe(key *Key, n ...string) Subscription

	// Variables gets the channel Variables
	Variables(key *Key) Variables
}

// ChannelData is the data for a specific channel
type ChannelData struct {
	ID           string      `json:"id"`    // Unique id for this channel (same as for AMI)
	Name         string      `json:"name"`  // Name of this channel (tech/name-id format)
	State        string      `json:"state"` // State of the channel
	Accountcode  string      `json:"accountcode"`
	Caller       CallerID    `json:"caller"`    // CallerId of the calling endpoint
	Connected    CallerID    `json:"connected"` // CallerId of the connected line
	Creationtime DateTime    `json:"creationtime"`
	Dialplan     DialplanCEP `json:"dialplan"` // Current location in the dialplan
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

//---
// Play/Record operations
//---

// Play initiates playback of the specified media uri
// to the channel, returning the Playback handle
func (ch *ChannelHandle) Play(id string, mediaURI string) (ph PlaybackHandle, err error) {
	ph, err = ch.c.Play(ch.key, id, mediaURI)
	return
}

// Record records the channel to the given filename
func (ch *ChannelHandle) Record(name string, opts *RecordingOptions) (rh LiveRecordingHandle, err error) {
	rh, err = ch.c.Record(ch.key, name, opts)
	return
}

// StagePlay stages a `Play` operation.
func (ch *ChannelHandle) StagePlay(id string, mediaURI string) (ph PlaybackHandle) {
	ph = ch.c.StagePlay(ch.key, id, mediaURI)
	return
}

// StageRecord stages a `Record` operation
func (ch *ChannelHandle) StageRecord(name string, opts *RecordingOptions) (rh LiveRecordingHandle) {
	rh = ch.c.StageRecord(ch.key, name, opts)
	return
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
		return false, errors.Wrap(err, "Failed to get updated channel")
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

// Variables returns the channel variables
func (ch *ChannelHandle) Variables() Variables {
	return ch.c.Variables(ch.key)
}

// --
// Misc
// --

// Originate creates (and dials) a new channel using the present channel as its Originator.
func (ch *ChannelHandle) Originate(req OriginateRequest) (*ChannelHandle, error) {
	if req.Originator == "" {
		req.Originator = ch.ID()
	}
	return ch.c.Originate(req)
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
func (ch *ChannelHandle) StageSnoop(snoopID string, opts *SnoopOptions) *ChannelHandle {
	return ch.c.StageSnoop(ch.key, snoopID, opts)
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

// Match returns true if the event matches the channel
func (ch *ChannelHandle) Match(e Event) bool {
	channelEvent, ok := e.(ChannelEvent)
	if !ok {
		return false
	}

	//channel ID comparisons
	//	do we compare based on id;N, where id == id and the N's are different
	//		 -> this happens in Local channels

	// NOTE: this code considers local channels equal
	//leftChannel := strings.Split(ch.key, ";")[0]
	channelIDs := channelEvent.GetChannelIDs()
	for _, i := range channelIDs {
		//rightChannel := strings.Split(i, ";")[0]
		if ch.key.ID == i {
			return true
		}
	}
	return false
}
