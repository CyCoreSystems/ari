package native

import (
	"github.com/CyCoreSystems/ari"
)

type nativeBridge struct {
	conn       *Conn
	subscriber ari.Subscriber
	playback   ari.Playback
}

func (b *nativeBridge) Create(id string, t string, name string) (bh *ari.BridgeHandle, err error) {

	type request struct {
		ID   string `json:"bridgeId,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	}

	req := request{id, t, name}
	var bd ari.BridgeData

	err = Post(b.conn, "/bridges/"+id, &req, &bd)
	if err != nil {
		return
	}

	bh = b.Get(bd.ID)
	return
}

func (b *nativeBridge) Get(id string) *ari.BridgeHandle {
	return ari.NewBridgeHandle(id, b)
}

func (b *nativeBridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges = []struct {
		ID string `json:"id"`
	}{}

	err = Get(b.conn, "/bridges", &bridges)
	for _, i := range bridges {
		bx = append(bx, b.Get(i.ID))
	}
	return
}

// Data returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (b *nativeBridge) Data(id string) (bd ari.BridgeData, err error) {
	err = Get(b.conn, "/bridges/"+id, &bd)
	return
}

// AddChannel adds a channel to a bridge
// Equivalent to Post /bridges/{id}/addChannel
func (b *nativeBridge) AddChannel(bridgeID string, channelID string) (err error) {

	type request struct {
		ChannelID string `json:"channel"`
		Role      string `json:"role,omitempty"`
	}

	req := request{channelID, ""}
	err = Post(b.conn, "/bridges/"+bridgeID+"/addChannel", nil, &req)
	return
}

// RemoveChannel removes the specified channel from a bridge
// Equivalent to Post /bridges/{id}/removeChannel
func (b *nativeBridge) RemoveChannel(id string, channelID string) (err error) {
	req := struct {
		ChannelID string `json:"channel"`
	}{
		ChannelID: channelID,
	}

	//pass request
	err = Post(b.conn, "/bridges/"+id+"/removeChannel", nil, &req)
	return
}

// BridgeDelete shuts down a bridge. If any channels are in this bridge,
// they will be removed and resume whatever they were doing beforehand.
// This means that the channels themselves are not deleted.
// Equivalent to DELETE /bridges/{id}
func (b *nativeBridge) Delete(id string) (err error) {
	err = Delete(b.conn, "/bridges/"+id, nil, "")
	return
}

func (b *nativeBridge) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = Post(b.conn, "/bridges/"+id+"/play/"+playbackID, resp, &req)
	ph = b.playback.Get(playbackID)
	return
}

func (b *nativeBridge) Subscribe(id string, n ...string) ari.Subscription {
	var ns nativeSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	go func() {
		sub := b.subscriber.Subscribe(n...)
		defer sub.Cancel()
		for {

			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt := <-sub.Events():

				//TODO: do we want to send in events on the bridge
				// for a specific channel?

				be, ok := evt.(ari.BridgeEvent)
				if !ok {
					// ignore non-channel events
					continue
				}

				Logger.Debug("Got bridge event", "bridgeid", be.GetBridgeID(), "eventtype", evt.GetType())

				if be.GetBridgeID() != id {
					// ignore unrelated channel events
					continue
				}

				ns.events <- evt
			}
		}
	}()

	return &ns
}
