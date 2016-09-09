package generic

import "github.com/CyCoreSystems/ari"

type Bridge struct {
	Conn     Conn
	Playback ari.Playback
}

func (b *Bridge) Get(id string) *ari.BridgeHandle {
	return ari.NewBridgeHandle(id, b)
}

func (b *Bridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges = []struct {
		ID string `json:"id"`
	}{}

	err = b.Conn.Get("/bridges", nil, &bridges)
	for _, i := range bridges {
		bx = append(bx, b.Get(i.ID))
	}
	return
}

// Data returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (b *Bridge) Data(id string) (bd ari.BridgeData, err error) {
	err = b.Conn.Get("/bridges/%s", []interface{}{id}, &bd)
	return
}

// AddChannel adds a channel to a bridge
// Equivalent to Post /bridges/{id}/addChannel
func (b *Bridge) AddChannel(bridgeID string, channelID string) (err error) {

	type request struct {
		ChannelID string `json:"channel"`
		Role      string `json:"role,omitempty"`
	}

	req := request{channelID, ""}
	err = b.Conn.Post("/bridges/%s/addChannel", []interface{}{bridgeID}, nil, &req)
	return
}

// RemoveChannel removes the specified channel from a bridge
// Equivalent to Post /bridges/{id}/removeChannel
func (b *Bridge) RemoveChannel(id string, channelID string) (err error) {
	req := struct {
		ChannelID string `json:"channel"`
	}{
		ChannelID: channelID,
	}

	//pass request
	err = b.Conn.Post("/bridges/%s/removeChannel", []interface{}{id}, nil, &req)
	return
}

// BridgeDelete shuts down a bridge. If any channels are in this bridge,
// they will be removed and resume whatever they were doing beforehand.
// This means that the channels themselves are not deleted.
// Equivalent to DELETE /bridges/{id}
func (b *Bridge) Delete(id string) (err error) {
	err = b.Conn.Delete("/bridges/%s", []interface{}{id}, nil, "")
	return
}

func (b *Bridge) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = b.Conn.Post("/bridges/%s/play/%s", []interface{}{id, playbackID}, resp, &req)
	ph = b.Playback.Get(playbackID)
	return
}
