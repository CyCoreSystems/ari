package generic

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/satori/go.uuid"
)

type Channel struct {
	Conn     Conn
	Playback ari.Playback
}

func (c *Channel) List() (cx []*ari.ChannelHandle, err error) {
	var channels = []struct {
		ID string `json:"id"`
	}{}

	err = c.Conn.Get("/channels", nil, &channels)
	for _, i := range channels {
		cx = append(cx, c.Get(i.ID))
	}

	return
}

func (c *Channel) Hangup(id, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return c.Conn.Delete("/channels/%s", []interface{}{id}, nil, req)
}

func (c *Channel) Data(id string) (cd ari.ChannelData, err error) {
	err = c.Conn.Get("/channels/%s", []interface{}{id}, &cd)
	return
}

func (c *Channel) Get(id string) *ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return ari.NewChannelHandle(id, c)
}

func (c *Channel) Create(req ari.OriginateRequest) (*ari.ChannelHandle, error) {
	id := uuid.NewV1().String()
	h := ari.NewChannelHandle(id, c)

	var err error
	err = c.Conn.Post("/channels/%s", []interface{}{id}, nil, &req)
	if err != nil {
		return nil, err
	}

	return h, err
}

func (c *Channel) Continue(id string, context, extension, priority string) (err error) {
	type request struct {
		//TODO: populate request
	}
	req := request{}
	err = c.Conn.Post("/channels/%s/continue", []interface{}{id}, nil, &req)
	return
}

func (c *Channel) Busy(id string) (err error) {
	err = c.Hangup(id, "busy")
	return
}

func (c *Channel) Congestion(id string) (err error) {
	err = c.Hangup(id, "congestion")
	return
}

func (c *Channel) Answer(id string) (err error) {
	err = c.Conn.Post("/channels/%s/answer", []interface{}{id}, nil, nil)
	return
}

func (c *Channel) Ring(id string) (err error) {
	err = c.Conn.Post("/channels/%s/ring", []interface{}{id}, nil, nil)
	return
}

func (c *Channel) StopRing(id string) (err error) {
	err = c.Conn.Delete("/channels/%s/ring", []interface{}{id}, nil, "")
	return
}

func (c *Channel) Hold(id string) (err error) {
	err = c.Conn.Post("/channels/%s/hold", []interface{}{id}, nil, nil)
	return
}

func (c *Channel) StopHold(id string) (err error) {
	err = c.Conn.Delete("/channels/%s/hold", []interface{}{id}, nil, "")
	return
}

func (c *Channel) Mute(id string, dir string) (err error) {
	type request struct {
		Direction string `json:"direction,omitempty"`
	}

	req := request{dir}
	err = c.Conn.Post("/channels/%s/mute", []interface{}{id}, nil, &req)
	return
}

func (c *Channel) Unmute(id string, dir string) (err error) {
	var req string
	if dir != "" {
		req = fmt.Sprintf("direction=%s", dir)
	}

	err = c.Conn.Delete("/channels/%s/mute", []interface{}{id}, nil, req)
	return
}

func (c *Channel) SendDTMF(id string, dtmf string) (err error) {
	type request struct {
		//TODO: populate request
	}
	req := request{}
	err = c.Conn.Post("/channels/%s/dtmf", []interface{}{id}, nil, &req)
	return
}

func (c *Channel) MOH(id string, mohClass string) (err error) {
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}
	req := request{mohClass}
	err = c.Conn.Post("/channels/%s/moh", []interface{}{id}, nil, &req)
	return
}

func (c *Channel) StopMOH(id string) (err error) {
	err = c.Conn.Delete("/channels/%s/moh", []interface{}{id}, nil, "")
	return
}

func (c *Channel) Silence(id string) (err error) {
	err = c.Conn.Post("/channels/%s/silence", []interface{}{id}, nil, nil)
	return
}

func (c *Channel) StopSilence(id string) (err error) {
	err = c.Conn.Delete("/channels/%s/silence", []interface{}{id}, nil, "")
	return
}

func (c *Channel) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = c.Conn.Post("/channels/%s/play/%s", []interface{}{id, playbackID}, resp, &req)
	ph = c.Playback.Get(playbackID)
	return
}
