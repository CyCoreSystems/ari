package native

import (
	"errors"
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"

	"github.com/satori/go.uuid"
)

// Channel provides the ARI Channel accessors for the native client
type Channel struct {
	client *Client
}

// List lists the current channels and returns the list of channel handles
func (c *Channel) List(filter *ari.Key) (cx []*ari.Key, err error) {
	var channels = []struct {
		ID string `json:"id"`
	}{}

	if filter == nil {
		filter = ari.NewKey(ari.ChannelKey, "")
	}

	err = c.client.get("/channels", &channels)
	for _, i := range channels {
		k := c.client.stamp(ari.NewKey(ari.ChannelKey, i.ID))
		if filter.Match(k) {
			cx = append(cx, k)
		}
	}

	return
}

// Hangup hangs up the given channel using the (optional) reason
func (c *Channel) Hangup(key *ari.Key, reason string) error {
	id := key.ID
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return c.client.del("/channels/"+id, nil, req)
}

// Data retrieves the current state of the channel
func (c *Channel) Data(key *ari.Key) (*ari.ChannelData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("channel key not supplied")
	}

	var data = new(ari.ChannelData)
	if err := c.client.get("/channels/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "channel", "%v", key.ID)
	}
	data.Key = c.client.stamp(key)

	return data, nil
}

// Get gets the lazy handle for the given channel
func (c *Channel) Get(key *ari.Key) *ari.ChannelHandle {
	return ari.NewChannelHandle(c.client.stamp(key), c, nil)
}

// Originate originates a channel and returns the handle
func (c *Channel) Originate(key *ari.Key, req ari.OriginateRequest) (*ari.ChannelHandle, error) {
	h, err := c.StageOriginate(key, req)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StageOriginate creates a new channel handle with a channel originate request
// staged.
func (c *Channel) StageOriginate(key *ari.Key, req ari.OriginateRequest) (*ari.ChannelHandle, error) {
	if key != nil && req.Originator == "" && key.Kind == ari.ChannelKey {
		req.Originator = key.ID
	}

	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	return ari.NewChannelHandle(c.client.stamp(ari.NewKey(ari.ChannelKey, req.ChannelID)), c,
		func(ch *ari.ChannelHandle) error {
			type response struct {
				ID string `json:"id"`
			}

			var resp response

			err := c.client.post("/channels", &resp, &req)
			if err != nil {
				return nil
			}

			return err
		},
	), nil
}

// Create creates a channel and returns the handle. TODO: expand
// differences between originate and create.
func (c *Channel) Create(key *ari.Key, req ari.ChannelCreateRequest) (*ari.ChannelHandle, error) {
	if key != nil && req.Originator == "" && key.Kind == ari.ChannelKey {
		req.Originator = key.ID
	}

	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	err := c.client.post("/channels/create", nil, &req)
	if err != nil {
		return nil, err
	}

	return ari.NewChannelHandle(c.client.stamp(ari.NewKey(ari.ChannelKey, req.ChannelID)), c, nil), nil
}

// Continue tells a channel to process to the given ARI context and extension
func (c *Channel) Continue(key *ari.Key, context, extension string, priority int) (err error) {
	req := struct {
		Context   string `json:"context"`
		Extension string `json:"extension"`
		Priority  int    `json:"priority"`
	}{
		Context:   context,
		Extension: extension,
		Priority:  priority,
	}
	return c.client.post("/channels/"+key.ID+"/continue", nil, &req)
}

// Busy sends the busy status code to the channel (TODO: does this play a busy signal too)
func (c *Channel) Busy(key *ari.Key) error {
	return c.Hangup(key, "busy")
}

// Congestion sends the congestion status code to the channel (TODO: does this play a tone?)
func (c *Channel) Congestion(key *ari.Key) error {
	return c.Hangup(key, "congestion")
}

// Answer answers a channel, if ringing (TODO: does this return an error if already answered?)
func (c *Channel) Answer(key *ari.Key) error {
	return c.client.post("/channels/"+key.ID+"/answer", nil, nil)
}

// Ring causes a channel to start ringing (TODO: does this return an error if already ringing?)
func (c *Channel) Ring(key *ari.Key) error {
	return c.client.post("/channels/"+key.ID+"/ring", nil, nil)
}

// StopRing causes a channel to stop ringing (TODO: does this return an error if not ringing?)
func (c *Channel) StopRing(key *ari.Key) error {
	return c.client.del("/channels/"+key.ID+"/ring", nil, "")
}

// Hold puts a channel on hold (TODO: does this return an error if already on hold?)
func (c *Channel) Hold(key *ari.Key) error {
	return c.client.post("/channels/"+key.ID+"/hold", nil, nil)
}

// StopHold removes a channel from hold (TODO: does this return an error if not on hold)
func (c *Channel) StopHold(key *ari.Key) (err error) {
	return c.client.del("/channels/"+key.ID+"/hold", nil, "")
}

// Mute mutes a channel in the given direction (TODO: does this return an error if already muted)
func (c *Channel) Mute(key *ari.Key, dir ari.Direction) error {
	if dir == "" {
		dir = ari.DirectionBoth
	}

	req := struct {
		Direction ari.Direction `json:"direction,omitempty"`
	}{
		Direction: dir,
	}
	return c.client.post("/channels/"+key.ID+"/mute", nil, &req)
}

// Unmute unmutes a channel in the given direction (TODO: does this return an error if unmuted)
func (c *Channel) Unmute(key *ari.Key, dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionBoth
	}
	req := fmt.Sprintf("direction=%s", dir)
	return c.client.del("/channels/"+key.ID+"/mute", nil, req)
}

// SendDTMF sends a string of digits and symbols to the channel
func (c *Channel) SendDTMF(key *ari.Key, dtmf string, opts *ari.DTMFOptions) error {

	if opts == nil {
		opts = &ari.DTMFOptions{}
	}

	if opts.Duration < 1 {
		opts.Duration = 100 // ARI default, for documenation
	}
	if opts.Between < 1 {
		opts.Between = 100 // ARI default, for documentation
	}

	req := struct {
		Dtmf     string `json:"dtmf,omitempty"`
		Before   int    `json:"before,omitempty"`
		Between  int    `json:"between,omitempty"`
		Duration int    `json:"duration,omitempty"`
		After    int    `json:"after,omitempty"`
	}{
		Dtmf:     dtmf,
		Before:   int(opts.Before / time.Millisecond),
		After:    int(opts.After / time.Millisecond),
		Duration: int(opts.Duration / time.Millisecond),
		Between:  int(opts.Between / time.Millisecond),
	}

	return c.client.post("/channels/"+key.ID+"/dtmf", nil, &req)
}

// MOH plays the given music on hold class to the channel TODO: does this error when already playing MOH?
func (c *Channel) MOH(key *ari.Key, class string) error {
	req := struct {
		Class string `json:"mohClass,omitempty"`
	}{
		Class: class,
	}
	return c.client.post("/channels/"+key.ID+"/moh", nil, &req)
}

// StopMOH stops any music on hold playing on the channel (TODO: does this error when no MOH is playing?)
func (c *Channel) StopMOH(key *ari.Key) error {
	return c.client.del("/channels/"+key.ID+"/moh", nil, "")
}

// Silence silences a channel (TODO: does this error when already silences)
func (c *Channel) Silence(key *ari.Key) error {
	return c.client.post("/channels/"+key.ID+"/silence", nil, nil)
}

// StopSilence stops the silence on a channel (TODO: does this error when not silenced)
func (c *Channel) StopSilence(key *ari.Key) error {
	return c.client.del("/channels/"+key.ID+"/silence", nil, "")
}

// Play plays the given media URI on the channel, using the playbackID as
// the identifier of the created ARI Playback entity
func (c *Channel) Play(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	h, err := c.StagePlay(key, playbackID, mediaURI)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StagePlay stages a `Play` operation on the bridge
func (c *Channel) StagePlay(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	resp := make(map[string]interface{})
	req := struct {
		Media string `json:"media"`
	}{
		Media: mediaURI,
	}

	playbackKey := c.client.stamp(ari.NewKey(ari.PlaybackKey, playbackID))
	return ari.NewPlaybackHandle(playbackKey, c.client.Playback(), func(pb *ari.PlaybackHandle) error {
		return c.client.post("/channels/"+key.ID+"/play/"+playbackID, &resp, &req)
	}), nil
}

// Record records audio on the channel, using the name parameter as the name of the
// created LiveRecording entity.
func (c *Channel) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	h, err := c.StageRecord(key, name, opts)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StageRecord stages a `Record` opreation
func (c *Channel) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {

	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	resp := make(map[string]interface{})
	req := struct {
		Name        string `json:"name"`
		Format      string `json:"format"`
		MaxDuration int    `json:"maxDurationSeconds"`
		MaxSilence  int    `json:"maxSilenceSeconds"`
		IfExists    string `json:"ifExists,omitempty"`
		Beep        bool   `json:"beep"`
		TerminateOn string `json:"terminateOn,omitempty"`
	}{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}

	recordingKey := c.client.stamp(ari.NewKey(ari.LiveRecordingKey, name))

	return ari.NewLiveRecordingHandle(recordingKey, c.client.LiveRecording(), func(h *ari.LiveRecordingHandle) error {
		return c.client.post("/channels/"+key.ID+"/record", &resp, &req)
	}), nil
}

// Snoop snoops on a channel, using the the given snoopID as the new channel handle ID (TODO: confirm and expand description)
func (c *Channel) Snoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) (*ari.ChannelHandle, error) {
	h, err := c.StageSnoop(key, snoopID, opts)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StageSnoop creates a new `ChannelHandle` with a `Snoop` operation staged.
func (c *Channel) StageSnoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) (*ari.ChannelHandle, error) {
	if opts == nil {
		opts = &ari.SnoopOptions{App: c.client.ApplicationName()}
	}
	if snoopID == "" {
		snoopID = uuid.NewV1().String()
	}

	// Create the snooping channel's key
	k := c.client.stamp(ari.NewKey(ari.ChannelKey, snoopID))

	return ari.NewChannelHandle(k, c, func(ch *ari.ChannelHandle) error {
		return c.client.post("/channels/"+key.ID+"/snoop/"+snoopID, nil, &opts)
	}), nil
}

// Dial dials the given calling channel identifier
func (c *Channel) Dial(key *ari.Key, callingChannelID string, timeout time.Duration) error {
	req := struct {
		// Caller is the (optional) channel ID of another channel to which media negotiations for the newly-dialed channel will be associated.
		Caller string `json:"caller,omitempty"`

		// Timeout is the maximum amount of time to allow for the dial to complete.
		Timeout int `json:"timeout"`
	}{
		Caller:  callingChannelID,
		Timeout: int(timeout.Seconds()),
	}

	return c.client.post("/channels/"+key.ID+"/dial", nil, &req)
}

// Subscribe creates a new subscription for ARI events related to this channel
func (c *Channel) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	return c.client.Bus().Subscribe(key, n...)
}

// GetVariable gets the value of the given variable
func (c *Channel) GetVariable(key *ari.Key, name string) (string, error) {
	var m struct {
		Value string `json:"value"`
	}

	err := c.client.get(fmt.Sprintf("/channels/%s/variable?variable=%s", key.ID, name), &m)
	return m.Value, err
}

// SetVariable sets the value of the given channel variable
func (c *Channel) SetVariable(key *ari.Key, name, value string) error {
	req := struct {
		Name  string `json:"variable"`
		Value string `json:"value,omitempty"`
	}{
		Name:  name,
		Value: value,
	}

	return c.client.post(fmt.Sprintf("/channels/%s/variable", key.ID), nil, &req)
}
