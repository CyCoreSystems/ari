package native

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/satori/go.uuid"
)

type nativeChannel struct {
	conn *Conn
}

func (c *nativeChannel) Hangup(id, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return Delete(c.conn, "/channels/"+id, nil, req)
}

func (c *nativeChannel) Data(id string) (cd ari.ChannelData) {
	_ = Get(c.conn, "/channels/"+id, &cd)
	return
}

func (c *nativeChannel) Get(id string) *ari.ChannelHandle {
	//TODO: does Get need to do anything else??
	return ari.NewChannelHandle(id, c)
}

func (c *nativeChannel) Create() (*ari.ChannelHandle, error) {
	id := uuid.NewV1().String()
	h := ari.NewChannelHandle(id, c)

	var err error
	type request struct {
		//TODO: populate request
	}
	req := request{}

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
