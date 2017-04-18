package native

import (
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
		filter = ari.NodeKey(c.client.ApplicationName(), c.client.node)
	}

	err = c.client.get("/channels", &channels)
	for _, i := range channels {
		k := ari.NewKey(ari.ChannelKey, i.ID, ari.WithApp(c.client.ApplicationName()), ari.WithNode(c.client.node))
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
func (c *Channel) Data(key *ari.Key) (cd *ari.ChannelData, err error) {
	id := key.ID
	cd = &ari.ChannelData{}
	err = c.client.get("/channels/"+id, cd)
	if err != nil {
		cd = nil
		err = dataGetError(err, "channel", "%v", id)
	}
	return
}

// Get gets the lazy handle for the given channel
func (c *Channel) Get(key *ari.Key) *ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return ari.NewChannelHandle(key, c, nil)
}

// Originate originates a channel and returns the handle TODO: expand
// differences between originate and create
func (c *Channel) Originate(req ari.OriginateRequest) (*ari.ChannelHandle, error) {
	h := c.StageOriginate(req)
	err := h.Exec()
	return h, err
}

// StageOriginate creates a new channel handle with a channel originate request
// staged.
func (c *Channel) StageOriginate(req ari.OriginateRequest) *ari.ChannelHandle {

	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	return ari.NewChannelHandle(ari.NewKey(ari.ChannelKey, req.ChannelID, ari.WithApp(c.client.ApplicationName()), ari.WithNode(c.client.node)), c, func(ch *ari.ChannelHandle) error {
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
func (c *Channel) Create(req ari.ChannelCreateRequest) (*ari.ChannelHandle, error) {
	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	err := c.client.post("/channels/create", nil, &req)
	if err != nil {
		return nil, err
	}

	h := ari.NewChannelHandle(ari.NewKey(ari.ChannelKey, req.ChannelID, ari.WithApp(c.client.ApplicationName()), ari.WithNode(c.client.node)), c, nil)
	return h, err
}

// Continue tells a channel to process to the given ARI context and extension
func (c *Channel) Continue(key *ari.Key, context, extension string, priority int) (err error) {
	id := key.ID
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
func (c *Channel) Busy(key *ari.Key) (err error) {
	err = c.Hangup(key, "busy")
	return
}

// Congestion sends the congestion status code to the channel (TODO: does this play a tone?)
func (c *Channel) Congestion(key *ari.Key) (err error) {
	err = c.Hangup(key, "congestion")
	return
}

// Answer answers a channel, if ringing (TODO: does this return an error if already answered?)
func (c *Channel) Answer(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.post("/channels/"+id+"/answer", nil, nil)
	return
}

// Ring causes a channel to start ringing (TODO: does this return an error if already ringing?)
func (c *Channel) Ring(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.post("/channels/"+id+"/ring", nil, nil)
	return
}

// StopRing causes a channel to stop ringing (TODO: does this return an error if not ringing?)
func (c *Channel) StopRing(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.del("/channels/"+id+"/ring", nil, "")
	return
}

// Hold puts a channel on hold (TODO: does this return an error if already on hold?)
func (c *Channel) Hold(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.post("/channels/"+id+"/hold", nil, nil)
	return
}

// StopHold removes a channel from hold (TODO: does this return an error if not on hold)
func (c *Channel) StopHold(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.del("/channels/"+id+"/hold", nil, "")
	return
}

// Mute mutes a channel in the given direction (TODO: does this return an error if already muted)
// TODO: enumerate direction
func (c *Channel) Mute(key *ari.Key, dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
	}
	id := key.ID
	req := struct {
		Direction ari.Direction `json:"direction,omitempty"`
	}{
		Direction: dir,
	}
	return c.client.post("/channels/"+id+"/mute", nil, &req)
}

// Unmute unmutes a channel in the given direction (TODO: does this return an error if unmuted)
// TODO: enumerate direction
func (c *Channel) Unmute(key *ari.Key, dir ari.Direction) (err error) {
	if dir == "" {
		dir = ari.DirectionIn
	}
	req := fmt.Sprintf("direction=%s", dir)
	id := key.ID
	return c.client.del("/channels/"+id+"/mute", nil, req)
}

// SendDTMF sends a string of digits and symbols to the channel
func (c *Channel) SendDTMF(key *ari.Key, dtmf string, opts *ari.DTMFOptions) (err error) {

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
	id := key.ID
	err = c.client.post("/channels/"+id+"/dtmf", nil, &req)
	return
}

// MOH plays the given music on hold class to the channel TODO: does this error when already playing MOH?
func (c *Channel) MOH(key *ari.Key, mohClass string) (err error) {
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}
	id := key.ID
	req := request{mohClass}
	err = c.client.post("/channels/"+id+"/moh", nil, &req)
	return
}

// StopMOH stops any music on hold playing on the channel (TODO: does this error when no MOH is playing?)
func (c *Channel) StopMOH(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.del("/channels/"+id+"/moh", nil, "")
	return
}

// Silence silences a channel (TODO: does this error when already silences)
func (c *Channel) Silence(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.post("/channels/"+id+"/silence", nil, nil)
	return
}

// StopSilence stops the silence on a channel (TODO: does this error when not silenced)
func (c *Channel) StopSilence(key *ari.Key) (err error) {
	id := key.ID
	err = c.client.del("/channels/"+id+"/silence", nil, "")
	return
}

// Play plays the given media URI on the channel, using the playbackID as
// the identifier of the created ARI Playback entity
func (c *Channel) Play(key *ari.Key, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	ph = c.StagePlay(key, playbackID, mediaURI)
	err = ph.Exec()
	return
}

// StagePlay stages a `Play` operation on the bridge
func (c *Channel) StagePlay(key *ari.Key, playbackID string, mediaURI string) (ph *ari.PlaybackHandle) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	id := key.ID

	playbackKey := ari.NewKey(ari.PlaybackKey, playbackID, ari.WithApp(c.client.ApplicationName()), ari.WithNode(c.client.node))
	ph = c.client.Playback().Get(playbackKey)
	return ari.NewPlaybackHandle(playbackKey, c.client.Playback(), func(pb *ari.PlaybackHandle) (err error) {
		err = c.client.post("/channels/"+id+"/play/"+playbackID, &resp, &req)
		return
	})
}

// Record records audio on the channel, using the name parameter as the name of the
// created LiveRecording entity.
func (c *Channel) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle, err error) {
	rh = c.StageRecord(key, name, opts)
	err = rh.Exec()
	return
}

// StageRecord stages a `Record` opreation
func (c *Channel) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle) {

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

	recordingKey := ari.NewKey(ari.LiveRecordingKey, name, ari.WithApp(c.client.ApplicationName()), ari.WithNode(c.client.node))
	id := key.ID
	return ari.NewLiveRecordingHandle(recordingKey, c.client.LiveRecording(), func() error {
		return c.client.post("/channels/"+id+"/record", &resp, &req)
	})
}

// Snoop snoops on a channel, using the the given snoopID as the new channel handle ID (TODO: confirm and expand description)
func (c *Channel) Snoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) (ch *ari.ChannelHandle, err error) {
	ch = c.StageSnoop(key, snoopID, opts)
	err = ch.Exec()
	return
}

// StageSnoop creates a new `ChannelHandle` with a `Snoop` operation staged.
func (c *Channel) StageSnoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) *ari.ChannelHandle {
	if opts == nil {
		opts = &ari.SnoopOptions{App: c.client.ApplicationName()}
	}
	if snoopID == "" {
		snoopID = uuid.NewV1().String()
	}
	id := key.ID
	return ari.NewChannelHandle(key, c, func(ch *ari.ChannelHandle) (err error) {
		err = c.client.post("/channels/"+id+"/snoop/"+snoopID, nil, &opts)
		return
	})
}

// Dial dials the given calling channel identifier
func (c *Channel) Dial(key *ari.Key, callingChannelID string, timeout time.Duration) (err error) {
	req := struct {
		// Caller is the (optional) channel ID of another channel to which media negotiations for the newly-dialed channel will be associated.
		Caller string `json:"caller,omitempty"`

		// Timeout is the maximum amount of time to allow for the dial to complete.
		Timeout int `json:"timeout"`
	}{
		Caller:  callingChannelID,
		Timeout: int(timeout.Seconds()),
	}
	id := key.ID
	err = c.client.post("/channels/"+id+"/dial", nil, &req)
	return
}

// Subscribe creates a new subscription for ARI events related to this channel
func (c *Channel) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	inSub := c.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		ch := c.Get(key)

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
	client *Client
	key    *ari.Key
}

// Variables returns the variables interface for channel
func (c *Channel) Variables(key *ari.Key) ari.Variables {
	return &ChannelVariables{c.client, key}
}

// Get gets the value of the given variable
func (v *ChannelVariables) Get(key string) (string, error) {
	type variable struct {
		Value string `json:"value"`
	}

	var m variable

	path := "/channels/" + v.key.ID + "/variable?variable=" + key
	err := v.client.get(path, &m)
	if err != nil {
		return "", err
	}
	return m.Value, nil
}

// Set sets the value of the given variable
func (v *ChannelVariables) Set(key string, value string) error {
	path := "/channels/" + v.key.ID + "/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{key, value}

	err := v.client.post(path, nil, &req)
	return err
}
