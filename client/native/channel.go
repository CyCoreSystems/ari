package native

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/satori/go.uuid"
)

type nativeChannel struct {
	conn     *Conn
	playback ari.Playback
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

func (c *nativeChannel) Continue(id string, context, extension, priority string) (err error) {
	type request struct {
		//TODO: populate request
	}
	req := request{}
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

func (c *nativeChannel) SendDTMF(id string, dtmf string) (err error) {
	type request struct {
		//TODO: populate request
	}
	req := request{}
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
	err = Post(c.conn, "/channels/"+id+"/play/"+playbackID, resp, &req)
	ph = c.playback.Get(playbackID)
	return
}
