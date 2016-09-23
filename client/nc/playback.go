package nc

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsPlayback struct {
	conn *Conn
}

func (p *natsPlayback) Get(id string) *ari.PlaybackHandle {
	return ari.NewPlaybackHandle(id, p)
}

func (p *natsPlayback) Data(id string) (d ari.PlaybackData, err error) {
	err = p.conn.readRequest("ari.playback.data."+id, nil, &d)
	return
}

func (p *natsPlayback) Control(id string, op string) (err error) {
	err = p.conn.standardRequest("ari.playback.control."+id, &op, nil)
	return
}

func (p *natsPlayback) Stop(id string) (err error) {
	err = p.conn.standardRequest("ari.playback.stop."+id, nil, nil)
	return
}

func (p *natsPlayback) Subscribe(id string, nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	playbackHandle := p.Get(id)

	go func() {
		for _, n := range nx {
			subj := fmt.Sprintf("ari.events.%s", n)
			sub, err := p.conn.conn.Subscribe(subj, func(msg *nats.Msg) {
				eventType := msg.Subject[len("ari.events."):]

				var ariMessage ari.Message
				ariMessage.SetRaw(&msg.Data)
				ariMessage.Type = eventType

				evt := ari.Events.Parse(&ariMessage)

				if playbackHandle.Match(evt) {
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
