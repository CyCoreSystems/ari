package native

import (
	"time"

	"github.com/CyCoreSystems/ari"
)

// Bridge provides the ARI Bridge accessors for the native client
type Bridge struct {
	client *Client
}

// Create creates a bridge and returns the lazy handle for the bridge
func (b *Bridge) Create(id string, t string, name string) (bh *ari.BridgeHandle, err error) {

	req := struct {
		ID   string `json:"bridgeId,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	}{
		ID:   id,
		Type: t,
		Name: name,
	}

	err = b.client.conn.Post("/bridges/"+id, &req, nil)
	if err != nil {
		return
	}

	bh = b.Get(id)
	return
}

// Get gets the lazy handle for the given bridge id
func (b *Bridge) Get(id string) *ari.BridgeHandle {
	return ari.NewBridgeHandle(id, b)
}

// List lists the current bridges and returns a list of lazy handles
func (b *Bridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges = []struct {
		ID string `json:"id"`
	}{}

	err = b.client.conn.Get("/bridges", &bridges)
	for _, i := range bridges {
		bx = append(bx, b.Get(i.ID))
	}
	return
}

// Data returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (b *Bridge) Data(id string) (bd ari.BridgeData, err error) {
	err = b.client.conn.Get("/bridges/"+id, &bd)
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
	err = b.client.conn.Post("/bridges/"+bridgeID+"/addChannel", nil, &req)
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
	err = b.client.conn.Post("/bridges/"+id+"/removeChannel", nil, &req)
	return
}

// Delete shuts down a bridge. If any channels are in this bridge,
// they will be removed and resume whatever they were doing beforehand.
// This means that the channels themselves are not deleted.
// Equivalent to DELETE /bridges/{id}
func (b *Bridge) Delete(id string) (err error) {
	err = b.client.conn.Delete("/bridges/"+id, nil, "")
	return
}

// Play attempts to play the given mediaURI on the bridge, using the playbackID
// as the identifier to the created playback handle
func (b *Bridge) Play(id string, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	err = b.client.conn.Post("/bridges/"+id+"/play/"+playbackID, &resp, &req)
	ph = b.client.Playback().Get(playbackID)
	return
}

// Record attempts to record audio on the bridge, using name as the identifier for
// the created live recording handle
func (b *Bridge) Record(id string, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle, err error) {

	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	resp := make(map[string]interface{})
	type request struct {
		Name        string `json:"name"`
		Format      string `json:"format"`
		MaxDuration int    `json:"maxDurationSeconds"`
		MaxSilence  int    `json:"maxSilenceSeconds"`
		IfExists    string `json:"ifExists,omitempty"`
		Beep        bool   `json:"beep"`
		TerminateOn string `json:"terminateOn,omitempty"`
	}
	req := request{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}
	err = b.client.conn.Post("/bridges/"+id+"/record", &resp, &req)
	if err != nil {
		rh = b.client.LiveRecording().Get(name)
	}
	return
}

// Subscribe creates an event subscription for events related to the given
// bridge␃entity
func (b *Bridge) Subscribe(id string, n ...string) ari.Subscription {
	var ns nativeSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	bridgeHandle := b.Get(id)

	go func() {
		sub := b.client.Bus().Subscribe(n...)
		defer sub.Cancel()
		for {

			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt := <-sub.Events():
				//TODO: do we want to send in events on the bridge for a specific channel?
				if bridgeHandle.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()

	return &ns
}
