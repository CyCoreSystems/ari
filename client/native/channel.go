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
func (c *Channel) List() (cx []*ari.ChannelHandle, err error) {
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
func (c *Channel) Get(id string) *ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return ari.NewChannelHandle(id, c)
}

// Originate originates a channel and returns the handle TODO: expand
// differences between originate and create
func (c *Channel) Originate(req ari.OriginateRequest) (*ari.ChannelHandle, error) {

	type response struct {
		ID string `json:"id"`
	}

	var resp response

	var err error
	err = c.client.post("/channels", &resp, &req)
	if err != nil {
		return nil, err
	}

	h := ari.NewChannelHandle(resp.ID, c)
	return h, err
}

// Create creates a channel and returns the handle. TODO: expand
// differences between originate and create.
func (c *Channel) Create(req ari.ChannelCreateRequest) (*ari.ChannelHandle, error) {
	if req.ChannelID == "" {
		req.ChannelID = uuid.NewV1().String()
	}

	var err error
	err = c.client.post("/channels/create", nil, &req)
	if err != nil {
		return nil, err
	}

	h := ari.NewChannelHandle(req.ChannelID, c)
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
func (c *Channel) Mute(id string, dir string) (err error) {
	type request struct {
		Direction string `json:"direction,omitempty"`
	}

	req := request{dir}
	err = c.client.post("/channels/"+id+"/mute", nil, &req)
	return
}

// Unmute unmutes a channel in the given direction (TODO: does this return an error if unmuted)
// TODO: enumerate direction
func (c *Channel) Unmute(id string, dir string) (err error) {
	var req string
	if dir != "" {
		req = fmt.Sprintf("direction=%s", dir)
	}

	err = c.client.del("/channels/"+id+"/mute", nil, req)
	return
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
func (c *Channel) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
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
func (c *Channel) Record(id string, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle, err error) {

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
func (c *Channel) Snoop(id string, snoopID string, app string, opts *ari.SnoopOptions) (ch *ari.ChannelHandle, err error) {
	if opts == nil {
		opts = &ari.SnoopOptions{}
	}

	resp := make(map[string]interface{})
	req := struct {
		Direction string `json:"spy,omitempty"`
		Whisper   string `json:"whisper,omitempty"`
		App       string `json:"app"`
		AppArgs   string `json:"appArgs"`
	}{
		Direction: opts.Direction,
		Whisper:   opts.Whisper,
		App:       app,
		AppArgs:   opts.AppArgs,
	}
	err = c.client.post("/channels/"+id+"/snoop/"+snoopID, &resp, &req)
	if err == nil {
		ch = c.Get(snoopID)
	}
	return
}

// Dial dials the given calling channel identifier
func (c *Channel) Dial(id string, callingChannelID string, timeout time.Duration) (err error) {
	req := struct {
		Caller  string `json:"caller,omitempty"` // the CHANNEL ID (not CallerID) of the channel for whom this dial is being generated
		Timeout int    `json:"timeout"`
	}{
		Caller:  callingChannelID,
		Timeout: int(timeout.Seconds()),
	}
	err = c.client.post("/channels/"+id+"/dial", nil, &req)
	return
}

// Subscribe creates a new subscription for ARI events related to this channel
func (c *Channel) Subscribe(id string, n ...string) ari.Subscription {
	var ns nativeSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	channelHandle := c.Get(id)

	go func() {
		sub := c.client.Bus().Subscribe(n...)
		defer sub.Cancel()
		for {

			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt := <-sub.Events():
				if channelHandle.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()

	return &ns
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
