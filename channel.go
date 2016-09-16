package ari

import (
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Channel represents a communication path interacting with an Asterisk server.
type Channel interface {
	// Get returns a handle to a channel for further interaction
	Get(id string) *ChannelHandle

	// List lists the channels in asterisk
	List() ([]*ChannelHandle, error)

	// Create creates a new channel, returning a handle to it or an
	// error, if the creation failed
	Create(OriginateRequest) (*ChannelHandle, error)

	// Data returns the channel data for a given channel
	Data(id string) (ChannelData, error)

	// Continue tells Asterisk to return a channel to the dialplan
	Continue(id, context, extension, priority string) error

	// Busy hangs up the channel with the "busy" cause code
	Busy(id string) error

	// Congestion hangs up the channel with the "congestion" cause code
	Congestion(id string) error

	// Answer answers the channel
	Answer(id string) error

	// Hangup hangs up the given channel
	Hangup(id string, reason string) error

	// Ring indicates ringing to the channel
	Ring(id string) error

	// StopRing stops ringing on the channel
	StopRing(id string) error

	// SendDTMF sends DTMF to the channel
	SendDTMF(id string, dtmf string) error

	// Hold puts the channel on hold
	Hold(id string) error

	// StopHold retrieves the channel from hold
	StopHold(id string) error

	// Mute mutes a channel in the given direction (in,out,both)
	Mute(id string, dir string) error

	// Unmute unmutes a channel in the given direction (in,out,both)
	Unmute(id string, dir string) error

	// MOH plays music on hold
	MOH(id string, moh string) error

	// StopMOH stops music on hold
	StopMOH(id string) error

	// Silence plays silence to the channel
	Silence(id string) error

	// StopSilence stops the silence on the channel
	StopSilence(id string) error

	// Play plays the media URI to the channel
	Play(id string, playbackID string, mediaURI string) (*PlaybackHandle, error)

	// Subscribe subscribes on the channel events
	Subscribe(id string, n ...string) Subscription
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

// NewChannelHandle returns a handle to the given ARI channel
func NewChannelHandle(id string, c Channel) *ChannelHandle {
	return &ChannelHandle{
		id: id,
		c:  c,
	}
}

// ChannelHandle provides a wrapper to a Channel interface for
// operations on a particular channel ID
type ChannelHandle struct {
	id string  // id of the channel on which we are operating
	c  Channel // the Channel interface with which we are operating
}

// ID returns the identifier for the channel handle
func (ch *ChannelHandle) ID() string {
	return ch.id
}

// Data returns the channel's data
func (ch *ChannelHandle) Data() (ChannelData, error) {
	return ch.c.Data(ch.id)
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension, priority string) error {
	return ch.c.Continue(ch.id, context, extension, priority)
}

//---
// Play operations
//---

// Play initiates playback of the specified media uri
// to the channel, returning the Playback's ID
func (ch *ChannelHandle) Play(mediaURI string) (ph *PlaybackHandle, err error) {
	id := uuid.NewV1().String()
	ph, err = ch.c.Play(ch.id, id, mediaURI)
	return
}

//---
// Hangup Operations
//---

// Busy hangs up the channel with the "busy" cause code
func (ch *ChannelHandle) Busy() error {
	return ch.c.Busy(ch.id)
}

// Congestion hangs up the channel with the congestion cause code
func (ch *ChannelHandle) Congestion() error {
	return ch.c.Congestion(ch.id)
}

// Hangup hangs up the channel with the normal cause code
func (ch *ChannelHandle) Hangup() error {
	return ch.c.Hangup(ch.id, "normal")
}

//--

// --
// Answer operations
// --

// Answer answers the channel
func (ch *ChannelHandle) Answer() error {
	return ch.c.Answer(ch.id)
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
	return ch.c.Ring(ch.id)
}

// StopRing stops ringing on the channel
func (ch *ChannelHandle) StopRing() error {
	return ch.c.StopRing(ch.id)
}

// ------

// --
// Mute operations
// --

// Mute mutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Mute(dir string) (err error) {
	if err = normalizeDirection(&dir); err != nil {
		return
	}

	return ch.c.Mute(ch.id, dir)
}

// Unmute unmutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Unmute(dir string) (err error) {
	if err = normalizeDirection(&dir); err != nil {
		return
	}

	return ch.c.Unmute(ch.id, dir)
}

// ----

// --
// Hold operations
// --

// Hold puts the channel on hold
func (ch *ChannelHandle) Hold() error {
	return ch.c.Hold(ch.id)
}

// StopHold retrieves the channel from hold
func (ch *ChannelHandle) StopHold() error {
	return ch.c.StopHold(ch.id)
}

// ----

// --
// Music on hold operations
// --

// MOH plays music on hold of the given class
// to the channel
func (ch *ChannelHandle) MOH(mohClass string) error {
	return ch.c.MOH(ch.id, mohClass)
}

// StopMOH stops playing of music on hold to the channel
func (ch *ChannelHandle) StopMOH() error {
	return ch.c.StopMOH(ch.id)
}

// ----

// --
// Silence operations
// --

// Silence plays silence to the channel
func (ch *ChannelHandle) Silence() error {
	return ch.c.Silence(ch.id)
}

// StopSilence stops silence to the channel
func (ch *ChannelHandle) StopSilence() error {
	return ch.c.StopSilence(ch.id)
}

// ----

// --
// Subscription
// --

// Subscribe subscribes the list of channel events
func (ch *ChannelHandle) Subscribe(n ...string) Subscription {
	return ch.c.Subscribe(ch.id, n...)
}

// TODO: rest of ChannelHandle
