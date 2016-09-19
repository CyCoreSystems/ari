package nc

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"

	"github.com/nats-io/nats"
)

type natsBridge struct {
	conn     *Conn
	playback ari.Playback
}

// CreateBridgeRequest is the request for creating bridges
type CreateBridgeRequest struct {
	ID   string `json:"bridgeId,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func (b *natsBridge) Create(id string, t string, name string) (h *ari.BridgeHandle, err error) {
	var bridgeID string
	req := CreateBridgeRequest{id, t, name}
	err = b.conn.standardRequest("ari.bridges.create", &req, &bridgeID)
	if err != nil {
		return
	}
	h = b.Get(bridgeID)
	return
}

func (b *natsBridge) Get(id string) *ari.BridgeHandle {
	return ari.NewBridgeHandle(id, b)
}

func (b *natsBridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges []string
	err = b.conn.readRequest("ari.bridges.all", nil, &bridges)
	for _, bridge := range bridges {
		bx = append(bx, b.Get(bridge))
	}
	return
}

func (b *natsBridge) Data(id string) (d ari.BridgeData, err error) {
	err = b.conn.readRequest("ari.bridges.data."+id, nil, &d)
	return
}

func (b *natsBridge) AddChannel(bridgeID string, channelID string) (err error) {
	err = b.conn.standardRequest("ari.bridges.addChannel."+bridgeID, channelID, nil)
	return
}

func (b *natsBridge) RemoveChannel(bridgeID string, channelID string) (err error) {
	err = b.conn.standardRequest("ari.bridges.removeChannel."+bridgeID, channelID, nil)
	return
}

func (b *natsBridge) Delete(id string) (err error) {
	err = b.conn.standardRequest("ari.bridges.delete."+id, nil, nil)
	return
}

// PlayRequest is the request for playback
type PlayRequest struct {
	PlaybackID string `json:"playback_id"`
	MediaURI   string `json:"media_uri"`
}

func (b *natsBridge) Play(id string, playbackID string, mediaURI string) (h *ari.PlaybackHandle, err error) {
	err = b.conn.standardRequest("ari.bridges.play."+id, &PlayRequest{PlaybackID: playbackID, MediaURI: mediaURI}, nil)
	if err == nil {
		h = b.playback.Get(playbackID)
	}
	return
}

func (b *natsBridge) Subscribe(id string, nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan v2.Eventer, 10)
	ns.closeChan = make(chan struct{})

	go func() {
		for _, n := range nx {
			subj := fmt.Sprintf("ari.events.%s", n)
			sub, err := b.conn.conn.Subscribe(subj, func(msg *nats.Msg) {
				eventType := msg.Subject[len("ari.events."):]

				var ariMessage v2.Message
				ariMessage.SetRaw(&msg.Data)
				ariMessage.Type = eventType

				evt := v2.Parse(&ariMessage)

				be, ok := evt.(ari.BridgeEvent)
				if !ok {
					// ignore non-channel events
					return
				}

				Logger.Debug("Got bridge event", "bridgeid", be.BridgeID(), "eventtype", evt.GetType())

				if be.BridgeID() != id {
					// ignore unrelated channel events
					return
				}

				ns.events <- evt
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
