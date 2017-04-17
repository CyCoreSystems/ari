package native

import (
	"fmt"
	"strings"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"

	"github.com/satori/go.uuid"
)

// Channel provides the ARI Channel accessors for the native client
type Channel struct {
	client *Client
}

// List lists the current channels and returns the list of channel handles
func (c *Channel) List() (cx []ari.ChannelHandle, err error) {
	var channels = []struct {
		ID string `json:"id"`
	}{}

	err = c.client.get("/channels", &channels)
	for _, i := range channels {
		cx = append(cx, c.Get(i.ID))
	}

	return
}

// Hangup hangs up the given channel using the (optional) reason
func (c *Channel) Hangup(id, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return c.client.del("/channels/"+id, nil, req)
}

// Data retrieves the current state of the channel
func (c *Channel) Data(id string) (cd *ari.ChannelData, err error) {
	cd = &ari.ChannelData{}
	err = c.client.get("/channels/"+id, cd)
	if err != nil {
		cd = nil
		err = dataGetError(err, "channel", "%v", id)
	}
	return
}

// Get gets the lazy handle for the given channel
func (c *Channel) Get(id string) ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return NewChannelHandle(id, c, nil)
}

// Originate originates a channel and returns the handle TODO: expand
// differences between originate and create
func (c *Channel) Originate(req ari.OriginateRequest) (ari.ChannelHandle, error) {
	h := c.StageOriginate(req)
	err := h.Exec()
	return h, err
}

// StageOriginate creates a new channel handle with a channel originate request
// staged.
func (c *Channel) StageOriginate(req ari.OriginateRequest) ari.ChannelHandle {

	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	return NewChannelHandle(req.ChannelID, c, func(ch *ChannelHandle) error {
		type response struct {
			ID string `json:"id"`
		}

		var resp response

		err := c.client.post("/channels", &resp, &req)
		if err != nil {
			return nil
		}

		return err
	})
}

// Create creates a channel and returns the handle. TODO: expand
// differences between originate and create.
func (c *Channel) Create(req ari.ChannelCreateRequest) (ari.ChannelHandle, error) {
	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	err := c.client.post("/channels/create", nil, &req)
	if err != nil {
		return nil, err
	}

	h := NewChannelHandle(req.ChannelID, c, nil)
	return h, err
}

// Continue tells a channel to process to the given ARI context and extension
func (c *Channel) Continue(id string, context, extension string, priority int) (err error) {
	type request struct {
		Context   string `json:"context"`
		Extension string `json:"extension"`
		Priority  int    `json:"priority"`
	}
	req := request{Context: context, Extension: extension, Priority: priority}
	err = c.client.post("/channels/"+id+"/continue", nil, &req)
	return
}

// Busy sends the busy status code to the channel (TODO: does this play a busy signal too)
func (c *Channel) Busy(id string) (err error) {
	err = c.Hangup(id, "busy")
	return
}

// Congestion sends the congestion status code to the channel (TODO: does this play a tone?)
func (c *Channel) Congestion(id string) (err error) {
	err = c.Hangup(id, "congestion")
	return
}

// Answer answers a channel, if ringing (TODO: does this return an error if already answered?)
func (c *Channel) Answer(id string) (err error) {
	err = c.client.post("/channels/"+id+"/answer", nil, nil)
	return
}

// Ring causes a channel to start ringing (TODO: does this return an error if already ringing?)
func (c *Channel) Ring(id string) (err error) {
	err = c.client.post("/channels/"+id+"/ring", nil, nil)
	return
}

// StopRing causes a channel to stop ringing (TODO: does this return an error if not ringing?)
func (c *Channel) StopRing(id string) (err error) {
	err = c.client.del("/channels/"+id+"/ring", nil, "")
	return
}

// Hold puts a channel on hold (TODO: does this return an error if already on hold?)
func (c *Channel) Hold(id string) (err error) {
	err = c.client.post("/channels/"+id+"/hold", nil, nil)
	return
}

// StopHold removes a channel from hold (TODO: does this return an error if not on hold)
func (c *Channel) StopHold(id string) (err error) {
	err = c.client.del("/channels/"+id+"/hold", nil, "")
	return
}

// Mute mutes a channel in the given direction (TODO: does this return an error if already muted)
// TODO: enumerate direction
func (c *Channel) Mute(id string, dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
	}

	req := struct {
		Direction ari.Direction `json:"direction,omitempty"`
	}{
		Direction: dir,
	}
	return c.client.post("/channels/"+id+"/mute", nil, &req)
}

// Unmute unmutes a channel in the given direction (TODO: does this return an error if unmuted)
// TODO: enumerate direction
func (c *Channel) Unmute(id string, dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
	}
	req := fmt.Sprintf("direction=%s", dir)

	return c.client.del("/channels/"+id+"/mute", nil, req)
}

// SendDTMF sends a string of digits and symbols to the channel
func (c *Channel) SendDTMF(id string, dtmf string, opts *ari.DTMFOptions) (err error) {

	type request struct {
		Dtmf     string `json:"dtmf,omitempty"`
		Before   *int   `json:"before,omitempty"`
		Between  *int   `json:"between,omitempty"`
		Duration *int   `json:"duration,omitempty"`
		After    *int   `json:"after,omitempty"`
	}
	req := request{}

	if opts != nil {
		if opts.Before != 0 {
			req.Before = new(int)
			*req.Before = int(opts.Before / time.Millisecond)
		}
		if opts.After != 0 {
			req.After = new(int)
			*req.After = int(opts.After / time.Millisecond)
		}
		if opts.Duration != 0 {
			req.Duration = new(int)
			*req.Duration = int(opts.Duration / time.Millisecond)
		}
		if opts.Between != 0 {
			req.Between = new(int)
			*req.Between = int(opts.Between / time.Millisecond)
		}
	}

	req.Dtmf = dtmf

	err = c.client.post("/channels/"+id+"/dtmf", nil, &req)
	return
}

// MOH plays the given music on hold class to the channel TODO: does this error when already playing MOH?
func (c *Channel) MOH(id string, mohClass string) (err error) {
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}
	req := request{mohClass}
	err = c.client.post("/channels/"+id+"/moh", nil, &req)
	return
}

// StopMOH stops any music on hold playing on the channel (TODO: does this error when no MOH is playing?)
func (c *Channel) StopMOH(id string) (err error) {
	err = c.client.del("/channels/"+id+"/moh", nil, "")
	return
}

// Silence silences a channel (TODO: does this error when already silences)
func (c *Channel) Silence(id string) (err error) {
	err = c.client.post("/channels/"+id+"/silence", nil, nil)
	return
}

// StopSilence stops the silence on a channel (TODO: does this error when not silenced)
func (c *Channel) StopSilence(id string) (err error) {
	err = c.client.del("/channels/"+id+"/silence", nil, "")
	return
}

// Play plays the given media URI on the channel, using the playbackID as
// the identifier of the created ARI Playback entity
func (c *Channel) Play(id string, playbackID string, mediaURI string) (ph ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = c.client.post("/channels/"+id+"/play/"+playbackID, &resp, &req)
	ph = c.client.Playback().Get(playbackID)
	return
}

// Record records audio on the channel, using the name parameter as the name of the
// created LiveRecording entity.
func (c *Channel) Record(id string, name string, opts *ari.RecordingOptions) (rh ari.LiveRecordingHandle, err error) {

	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	resp := make(map[string]interface{})
	type request struct {
		Name        string `json:"name"`
		Format      string `json:"format"`
		MaxDuration int    `json:"maxDurationSeconds"`
		MaxSilence  int    `json:"maxSilenceSeconds"`
		IfExists    string `json:"ifExists,omitempty"`
		Beep        bool   `json:"beep"`
		TerminateOn string `json:"terminateOn,omitempty"`
	}
	req := request{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}
	err = c.client.post("/channels/"+id+"/record", &resp, &req)
	if err != nil {
		rh = c.client.LiveRecording().Get(name)
	}
	return
}

// Snoop snoops on a channel, using the the given snoopID as the new channel handle ID (TODO: confirm and expand description)
func (c *Channel) Snoop(id string, snoopID string, opts *ari.SnoopOptions) (ch ari.ChannelHandle, err error) {
	ch = c.StageSnoop(id, snoopID, opts)
	err = ch.Exec()
	return
}

// StageSnoop creates a new `ChannelHandle` with a `Snoop` operation staged.
func (c *Channel) StageSnoop(id string, snoopID string, opts *ari.SnoopOptions) ari.ChannelHandle {
	if opts == nil {
		opts = &ari.SnoopOptions{App: c.client.ApplicationName()}
	}
	if snoopID == "" {
		snoopID = uuid.NewV1().String()
	}
	return NewChannelHandle(id, c, func(ch *ChannelHandle) (err error) {
		err = c.client.post("/channels/"+id+"/snoop/"+snoopID, nil, &opts)
		return
	})
}

// Dial dials the given calling channel identifier
func (c *Channel) Dial(id string, callingChannelID string, timeout time.Duration) (err error) {
	req := struct {
		// Caller is the (optional) channel ID of another channel to which media negotiations for the newly-dialed channel will be associated.
		Caller string `json:"caller,omitempty"`

		// Timeout is the maximum amount of time to allow for the dial to complete.
		Timeout int `json:"timeout"`
	}{
		Caller:  callingChannelID,
		Timeout: int(timeout.Seconds()),
	}
	err = c.client.post("/channels/"+id+"/dial", nil, &req)
	return
}

// Subscribe creates a new subscription for ARI events related to this channel
func (c *Channel) Subscribe(id string, n ...string) ari.Subscription {
	inSub := c.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		ch := c.Get(id)

		for {
			select {
			case <-outSub.closedChan:
				return
			case e, ok := <-inSub.Events():
				if !ok {
					return
				}
				if ch.Match(e) {
					outSub.events <- e
				}
			}
		}
	}()

	return outSub
}

// ChannelVariables provides the ARI Variables accessor scoped to a channel identifier for the native client
type ChannelVariables struct {
	client    *Client
	channelID string
}

// Variables returns the variables interface for channel
func (c *Channel) Variables(id string) ari.Variables {
	return &ChannelVariables{c.client, id}
}

// Get gets the value of the given variable
func (v *ChannelVariables) Get(key string) (string, error) {
	type variable struct {
		Value string `json:"value"`
	}

	var m variable

	path := "/channels/" + v.channelID + "/variable?variable=" + key
	err := v.client.get(path, &m)
	if err != nil {
		return "", err
	}
	return m.Value, nil
}

// Set sets the value of the given variable
func (v *ChannelVariables) Set(key string, value string) error {
	path := "/channels/" + v.channelID + "/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{key, value}

	err := v.client.post(path, nil, &req)
	return err
}

// ChannelHandle provides a wrapper on the Channel interface for operations on a particular channel ID.
type ChannelHandle struct {
	id string
	c  *Channel

	exec func(ch *ChannelHandle) error
}

// NewChannelHandle returns a handle to the given ARI channel
func NewChannelHandle(id string, c *Channel, exec func(ch *ChannelHandle) error) *ChannelHandle {
	return &ChannelHandle{
		id:   id,
		c:    c,
		exec: exec,
	}
}

// ID returns the identifier for the channel handle
func (ch *ChannelHandle) ID() string {
	return ch.id
}

// Exec executes any staged channel operations attached to this handle.
func (ch *ChannelHandle) Exec() (err error) {
	if ch.exec != nil {
		err = ch.exec(ch)
		ch.exec = nil
	}
	return err
}

// Data returns the channel's data
func (ch *ChannelHandle) Data() (*ari.ChannelData, error) {
	return ch.c.Data(ch.id)
}

// Continue tells Asterisk to return the channel to the dialplan
func (ch *ChannelHandle) Continue(context, extension string, priority int) error {
	return ch.c.Continue(ch.id, context, extension, priority)
}

//---
// Play/Record operations
//---

// Play initiates playback of the specified media uri
// to the channel, returning the Playback handle
func (ch *ChannelHandle) Play(id string, mediaURI string) (ph ari.PlaybackHandle, err error) {
	ph, err = ch.c.Play(ch.id, id, mediaURI)
	return
}

// Record records the channel to the given filename
func (ch *ChannelHandle) Record(name string, opts *ari.RecordingOptions) (rh ari.LiveRecordingHandle, err error) {
	rh, err = ch.c.Record(ch.id, name, opts)
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
func (ch *ChannelHandle) Mute(dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
	}

	return ch.c.Mute(ch.id, dir)
}

// Unmute unmutes the channel in the given direction (in, out, both)
func (ch *ChannelHandle) Unmute(dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
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

// Variables returns the channel variables
func (ch *ChannelHandle) Variables() ari.Variables {
	return ch.c.Variables(ch.id)
}

// --
// Misc
// --

// Originate creates (and dials) a new channel using the present channel as its Originator.
func (ch *ChannelHandle) Originate(req ari.OriginateRequest) (ari.ChannelHandle, error) {
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
	return ch.c.Dial(ch.id, caller, timeout)
}

// Snoop spies on a specific channel, creating a new snooping channel placed into the given app
func (ch *ChannelHandle) Snoop(snoopID string, opts *ari.SnoopOptions) (ari.ChannelHandle, error) {
	return ch.c.Snoop(ch.id, snoopID, opts)
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
func (ch *ChannelHandle) Subscribe(n ...string) ari.Subscription {
	if ch == nil {
		return nil
	}
	return ch.c.Subscribe(ch.id, n...)
}

// TODO: rest of ChannelHandle

// --
// DTMF
// --

// SendDTMF sends the DTMF information to the server
func (ch *ChannelHandle) SendDTMF(dtmf string, opts *ari.DTMFOptions) error {
	return ch.c.SendDTMF(ch.id, dtmf, opts)
}

// Match returns true if the event matches the channel
func (ch *ChannelHandle) Match(e ari.Event) bool {
	channelEvent, ok := e.(ari.ChannelEvent)
	if !ok {
		return false
	}

	//channel ID comparisons
	//	do we compare based on id;N, where id == id and the N's are different
	//		 -> this happens in Local channels

	// NOTE: this code considers local channels equal
	//leftChannel := strings.Split(ch.id, ";")[0]
	channelIDs := channelEvent.GetChannelIDs()
	for _, i := range channelIDs {
		//rightChannel := strings.Split(i, ";")[0]
		if ch.id == i {
			return true
		}
	}
	return false
}
