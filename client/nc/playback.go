package nc

import "github.com/CyCoreSystems/ari"

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
