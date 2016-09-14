package nc

import (
	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsPlayback struct {
	conn *nats.Conn
}

func (p *natsPlayback) Get(id string) *ari.PlaybackHandle {
	return ari.NewPlaybackHandle(id, p)
}

func (p *natsPlayback) Data(id string) (d ari.PlaybackData, err error) {
	err = request(p.conn, "ari.playback.data."+id, nil, &d)
	return
}

func (p *natsPlayback) Control(id string, op string) (err error) {
	err = request(p.conn, "ari.playback.control."+id, &op, nil)
	return
}

func (p *natsPlayback) Stop(id string) (err error) {
	err = request(p.conn, "ari.playback.stop."+id, nil, nil)
	return
}
