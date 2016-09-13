package nc

import (
	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsBridge struct {
	conn     *nats.Conn
	playback ari.Playback
}

func (b *natsBridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges []string
	err = request(b.conn, "ari.bridges.all", nil, &bridges)
	for _, bridge := range bridges {
		bx = append(bx, ari.NewBridgeHandle(bridge, b)) // TODO: replace NewBridgeHandle with b.Get
	}
	return
}

func (b *natsBridge) Data(id string) (d ari.BridgeData, err error) {
	err = request(b.conn, "ari.bridges.data."+id, nil, &d)
	return
}

func (b *natsBridge) AddChannel(bridgeID string, channelID string) (err error) {
	err = request(b.conn, "ari.bridges.addChannel."+bridgeID, channelID, nil)
	return
}

func (b *natsBridge) RemoveChannel(bridgeID string, channelID string) (err error) {
	err = request(b.conn, "ari.bridges.removeChannel."+bridgeID, channelID, nil)
	return
}

func (b *natsBridge) Delete(id string) (err error) {
	err = request(b.conn, "ari.bridges.delete."+id, nil, nil)
	return
}

// PlayRequest is the request for playback
type PlayRequest struct {
	PlaybackID string `json:"playback_id"`
	MediaURI   string `json:"media_uri"`
}

func (b *natsBridge) Play(id string, playbackID string, mediaURI string) (h *ari.PlaybackHandle, err error) {
	err = request(b.conn, "ari.bridges.play."+id, &PlayRequest{PlaybackID: playbackID, MediaURI: mediaURI}, nil)
	if err == nil {
		h = b.playback.Get(playbackID)
	}
	return
}
