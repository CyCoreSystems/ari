package native

import (
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"

	"github.com/satori/go.uuid"
)

type nativeChannel struct {
	conn          *Conn
	subscriber    ari.Subscriber
	playback      ari.Playback
	liveRecording ari.LiveRecording
}

func (c *nativeChannel) Playback() ari.Playback {
	return c.playback
}

func (c *nativeChannel) List() (cx []*ari.ChannelHandle, err error) {
	var channels = []struct {
		ID string `json:"id"`
	}{}

	err = Get(c.conn, "/channels", &channels)
	for _, i := range channels {
		cx = append(cx, c.Get(i.ID))
	}

	return
}

func (c *nativeChannel) Hangup(id, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return Delete(c.conn, "/channels/"+id, nil, req)
}

func (c *nativeChannel) Data(id string) (cd ari.ChannelData, err error) {
	err = Get(c.conn, "/channels/"+id, &cd)
	return
}

func (c *nativeChannel) Get(id string) *ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return ari.NewChannelHandle(id, c)
}

func (c *nativeChannel) Create(req ari.OriginateRequest) (*ari.ChannelHandle, error) {
	id := uuid.NewV1().String()
	h := ari.NewChannelHandle(id, c)

	var err error
	err = Post(c.conn, "/channels/"+id, nil, &req)
	if err != nil {
		return nil, err
	}

	return h, err
}

func (c *nativeChannel) Continue(id string, context, extension string, priority int) (err error) {
	type request struct {
		Context   string `json:"context"`
		Extension string `json:"extension"`
		Priority  int    `json:"priority"`
	}
	req := request{Context: context, Extension: extension, Priority: priority}
	err = Post(c.conn, "/channels/"+id+"/continue", nil, &req)
	return
}

func (c *nativeChannel) Busy(id string) (err error) {
	err = c.Hangup(id, "busy")
	return
}

func (c *nativeChannel) Congestion(id string) (err error) {
	err = c.Hangup(id, "congestion")
	return
}

func (c *nativeChannel) Answer(id string) (err error) {
	err = Post(c.conn, "/channels/"+id+"/answer", nil, nil)
	return
}

func (c *nativeChannel) Ring(id string) (err error) {
	err = Post(c.conn, "/channels/"+id+"/ring", nil, nil)
	return
}

func (c *nativeChannel) StopRing(id string) (err error) {
	err = Delete(c.conn, "/channels/"+id+"/ring", nil, "")
	return
}

func (c *nativeChannel) Hold(id string) (err error) {
	err = Post(c.conn, "/channels/"+id+"/hold", nil, nil)
	return
}

func (c *nativeChannel) StopHold(id string) (err error) {
	err = Delete(c.conn, "/channels/"+id+"/hold", nil, "")
	return
}

func (c *nativeChannel) Mute(id string, dir string) (err error) {
	type request struct {
		Direction string `json:"direction,omitempty"`
	}

	req := request{dir}
	err = Post(c.conn, "/channels/"+id+"/mute", nil, &req)
	return
}

func (c *nativeChannel) Unmute(id string, dir string) (err error) {
	var req string
	if dir != "" {
		req = fmt.Sprintf("direction=%s", dir)
	}

	err = Delete(c.conn, "/channels/"+id+"/mute", nil, req)
	return
}

func (c *nativeChannel) SendDTMF(id string, dtmf string, opts *ari.DTMFOptions) (err error) {

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

	err = Post(c.conn, "/channels/"+id+"/dtmf", nil, &req)
	return
}

func (c *nativeChannel) MOH(id string, mohClass string) (err error) {
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}
	req := request{mohClass}
	err = Post(c.conn, "/channels/"+id+"/moh", nil, &req)
	return
}

func (c *nativeChannel) StopMOH(id string) (err error) {
	err = Delete(c.conn, "/channels/"+id+"/moh", nil, "")
	return
}

func (c *nativeChannel) Silence(id string) (err error) {
	err = Post(c.conn, "/channels/"+id+"/silence", nil, nil)
	return
}

func (c *nativeChannel) StopSilence(id string) (err error) {
	err = Delete(c.conn, "/channels/"+id+"/silence", nil, "")
	return
}

func (c *nativeChannel) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = Post(c.conn, "/channels/"+id+"/play/"+playbackID, &resp, &req)
	ph = c.playback.Get(playbackID)
	return
}

func (c *nativeChannel) Record(id string, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle, err error) {

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
	err = Post(c.conn, "/channels/"+id+"/record", &resp, &req)
	if err != nil {
		rh = c.liveRecording.Get(name)
	}
	return
}

func (c *nativeChannel) Snoop(id string, snoopID string, app string, opts *ari.SnoopOptions) (ch *ari.ChannelHandle, err error) {
	if opts == nil {
		opts = &ari.SnoopOptions{}
	}

	resp := make(map[string]interface{})
	type request struct {
		Direction string `json:"spy,omitempty"`
		Whisper   string `json:"whisper,omitempty"`
		App       string `json:"app"`
		AppArgs   string `json:"appArgs"`
	}
	req := request{opts.Direction, opts.Whisper, app, opts.AppArgs}
	err = Post(c.conn, "/channels/"+id+"/snoop/"+snoopID, &resp, &req)
	if err == nil {
		ch = c.Get(snoopID)
	}
	return
}

func (c *nativeChannel) Dial(id string, caller string, timeout time.Duration) (err error) {
	type request struct {
		Caller  string `json:"caller"`
		Timeout int    `json:"timeout"`
	}
	//TODO: the dial documentation does not reference the unit of timeout,
	// second is assumed from similar parameters
	req := request{caller, int(timeout / time.Second)}
	err = Post(c.conn, "/channels/"+id+"/dial", nil, &req)
	return
}

func (c *nativeChannel) Subscribe(id string, n ...string) ari.Subscription {
	var ns nativeSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	channelHandle := c.Get(id)

	go func() {
		sub := c.subscriber.Subscribe(n...)
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

type nativeSubscription struct {
	closeChan chan struct{}
	events    chan ari.Event
}

func (ns *nativeSubscription) Events() chan ari.Event {
	return ns.events
}

func (ns *nativeSubscription) Cancel() {
	if ns.closeChan != nil {
		close(ns.closeChan)
	}
}
