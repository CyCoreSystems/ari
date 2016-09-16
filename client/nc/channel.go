package nc

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"

	"github.com/nats-io/nats"
)

// ContinueRequest is the request body for continuing over the message queue
type ContinueRequest struct {
	Context   string `json:"context"`
	Extension string `json:"extension"`
	Priority  string `json:"priority"`
}

type natsChannel struct {
	conn     *Conn
	playback ari.Playback
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

func (c *natsChannel) Continue(id string, context string, extension string, priority string) (err error) {
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

func (c *natsChannel) SendDTMF(id string, dtmf string) (err error) {
	err = c.conn.standardRequest("ari.channels.senddtmf."+id, &dtmf, nil)
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

func (c *natsChannel) Play(id string, playbackID string, mediaURI string) (p *ari.PlaybackHandle, err error) {
	err = c.conn.standardRequest("ari.channels.play."+id, &PlayRequest{
		PlaybackID: playbackID,
		MediaURI:   mediaURI,
	}, nil)

	if err == nil {
		p = c.playback.Get(playbackID)
	}

	return
}

func (c *natsChannel) Subscribe(id string, nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan v2.Eventer, 10)
	ns.closeChan = make(chan struct{})

	go func() {
		for _, n := range nx {
			subj := fmt.Sprintf("ari.channels.events.%s.%s", n, id)
			sub, err := c.conn.conn.Subscribe(subj, func(msg *nats.Msg) {
				//TODO: convert from nats messages to ARI events
				//TODO: send ari event to ns.events
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

type natsSubscription struct {
	closeChan chan struct{}
	events    chan v2.Eventer
}

func (ns *natsSubscription) Events() chan v2.Eventer {
	return ns.events
}

func (ns *natsSubscription) Cancel() {
	if ns.closeChan != nil {
		close(ns.closeChan)
	}
}
