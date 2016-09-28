package nc

import (
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"

	"github.com/nats-io/nats"
)

// ContinueRequest is the request body for continuing over the message queue
type ContinueRequest struct {
	Context   string `json:"context"`
	Extension string `json:"extension"`
	Priority  int    `json:"priority"`
}

type natsChannel struct {
	conn          *Conn
	playback      ari.Playback
	liveRecording ari.LiveRecording
}

func (c *natsChannel) Playback() ari.Playback {
	return c.playback
}

func (c *natsChannel) Get(id string) *ari.ChannelHandle {
	return ari.NewChannelHandle(id, c)
}

func (c *natsChannel) List() (cx []*ari.ChannelHandle, err error) {
	var channels []string
	err = c.conn.readRequest("ari.channels.all", nil, &channels)
	for _, ch := range channels {
		cx = append(cx, c.Get(ch))
	}
	return
}

func (c *natsChannel) Create(req ari.OriginateRequest) (h *ari.ChannelHandle, err error) {
	var channelID string
	err = c.conn.standardRequest("ari.channels.create", &req, &channelID)
	if err != nil {
		return
	}
	h = c.Get(channelID)
	return
}

func (c *natsChannel) Data(id string) (cd ari.ChannelData, err error) {
	err = c.conn.readRequest("ari.channels.data."+id, nil, &cd)
	return
}

func (c *natsChannel) Continue(id string, context string, extension string, priority int) (err error) {
	err = c.conn.standardRequest("ari.channels.continue."+id, &ContinueRequest{
		Context:   context,
		Extension: extension,
		Priority:  priority,
	}, nil)
	return
}

func (c *natsChannel) Busy(id string) (err error) {
	err = c.Hangup(id, "busy")
	return
}

func (c *natsChannel) Congestion(id string) (err error) {
	err = c.Hangup(id, "congestion")
	return
}

func (c *natsChannel) Hangup(id string, reason string) (err error) {
	err = c.conn.standardRequest("ari.channels.hangup."+id, &reason, nil)
	return
}

func (c *natsChannel) Answer(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.answer."+id, nil, nil)
	return
}

func (c *natsChannel) Ring(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.ring."+id, nil, nil)
	return
}

func (c *natsChannel) StopRing(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.stopring."+id, nil, nil)
	return
}

func (c *natsChannel) SendDTMF(id string, dtmf string, opts *ari.DTMFOptions) (err error) {
	if opts == nil {
		opts = &ari.DTMFOptions{}
	}

	type request struct {
		Dtmf string           `json:"dtmf,omitempty"`
		Opts *ari.DTMFOptions `json:"options,omitempty"`
	}

	req := request{dtmf, opts}

	err = c.conn.standardRequest("ari.channels.dtmf."+id, &req, nil)
	return
}

func (c *natsChannel) Hold(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.hold."+id, nil, nil)
	return
}

func (c *natsChannel) StopHold(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.stophold."+id, nil, nil)
	return
}

func (c *natsChannel) Mute(id string, dir string) (err error) {
	err = c.conn.standardRequest("ari.channels.mute."+id, &dir, nil)
	return
}

func (c *natsChannel) Unmute(id string, dir string) (err error) {
	err = c.conn.standardRequest("ari.channels.unmute."+id, &dir, nil)
	return
}

func (c *natsChannel) MOH(id string, moh string) (err error) {
	err = c.conn.standardRequest("ari.channels.moh."+id, &moh, nil)
	return
}

func (c *natsChannel) StopMOH(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.stopmoh."+id, nil, nil)
	return
}

func (c *natsChannel) Silence(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.silence."+id, nil, nil)
	return
}

func (c *natsChannel) StopSilence(id string) (err error) {
	err = c.conn.standardRequest("ari.channels.stopsilence."+id, nil, nil)
	return
}

// SnoopRequest is the NATs snoop request
type SnoopRequest struct {
	SnoopID string
	App     string
	Options *ari.SnoopOptions
}

func (c *natsChannel) Snoop(id string, snoopID string, app string, opts *ari.SnoopOptions) (ch *ari.ChannelHandle, err error) {
	req := &SnoopRequest{
		SnoopID: snoopID,
		App:     app,
		Options: opts,
	}
	err = c.conn.standardRequest("ari.channels.snoop."+id, &req, nil)
	if err == nil {
		ch = c.Get(snoopID)
	}
	return
}

// DialRequest is the request for the channel dial operation
type DialRequest struct {
	Caller  string `json:"caller"`
	Timeout int    `json:"timeout"`
}

func (c *natsChannel) Dial(id string, caller string, timeout time.Duration) (err error) {
	//TODO: the dial documentation does not reference the unit of timeout,
	// second is assumed from similar parameters
	req := DialRequest{caller, int(timeout / time.Second)}
	err = c.conn.standardRequest("ari.channels.dial."+id, &req, nil)
	return
}

func (c *natsChannel) Play(id string, playbackID string, mediaURI string) (p *ari.PlaybackHandle, err error) {
	err = c.conn.standardRequest("ari.channels.play."+id, &PlayRequest{
		PlaybackID: playbackID,
		MediaURI:   mediaURI,
	}, nil)

	p = c.playback.Get(playbackID)

	return
}

func (c *natsChannel) Record(id string, name string, opts *ari.RecordingOptions) (h *ari.LiveRecordingHandle, err error) {

	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	req := RecordRequest{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}
	err = c.conn.standardRequest("ari.channels.record."+id, req, nil)
	if err == nil {
		h = c.liveRecording.Get(name)
	}
	return
}

func (c *natsChannel) Subscribe(id string, nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	channelHandle := c.Get(id)

	go func() {
		for _, n := range nx {
			subj := fmt.Sprintf("ari.events.%s", n)
			sub, err := c.conn.conn.Subscribe(subj, func(msg *nats.Msg) {
				eventType := msg.Subject[len("ari.events."):]

				var ariMessage ari.Message
				ariMessage.SetRaw(&msg.Data)
				ariMessage.Type = eventType

				evt := ari.Events.Parse(&ariMessage)

				if channelHandle.Match(evt) {
					ns.events <- evt
				}
			})
			if err != nil {
				//TODO: handle error
				panic(err)
			}
			defer sub.Unsubscribe()
		}

		<-ns.closeChan
		ns.closeChan = nil
	}()

	return &ns
}
