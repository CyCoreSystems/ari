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
func (b *Bridge) Create(key *ari.Key, t string, name string) (bh *ari.BridgeHandle, err error) {
	bh = b.StageCreate(key, t, name)
	err = bh.Exec()
	return
}

// StageCreate creates a new bridge handle, staged with a bridge `Create` operation.
func (b *Bridge) StageCreate(key *ari.Key, t string, name string) *ari.BridgeHandle {
	id := key.ID
	req := struct {
		ID   string `json:"bridgeId,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	}{
		ID:   id,
		Type: t,
		Name: name,
	}

	return ari.NewBridgeHandle(key, b, func(bh *ari.BridgeHandle) (err error) {
		err = b.client.post("/bridges/"+id, &req, nil)
		if err != nil {
			return
		}
		return
	})
}

// Get gets the lazy handle for the given bridge id
func (b *Bridge) Get(key *ari.Key) *ari.BridgeHandle {
	return ari.NewBridgeHandle(key, b, nil)
}

// List lists the current bridges and returns a list of lazy handles
func (b *Bridge) List(filter *ari.Key) (bx []*ari.Key, err error) {
	var bridges = []struct {
		ID string `json:"id"`
	}{}

	if filter == nil {
		filter = ari.NodeKey(b.client.ApplicationName(), b.client.node)
	}

	err = b.client.get("/bridges", &bridges)
	for _, i := range bridges {
		k := ari.NewKey(ari.BridgeKey, i.ID, ari.WithApp(b.client.ApplicationName()), ari.WithNode(b.client.node))
		if filter.Match(k) {
			bx = append(bx, k)
		}
	}
	return
}

// Data returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (b *Bridge) Data(key *ari.Key) (bd *ari.BridgeData, err error) {
	bd = &ari.BridgeData{}
	id := key.ID
	err = b.client.get("/bridges/"+id, bd)
	if err != nil {
		bd = nil
		err = dataGetError(err, "bridge", "%v", id)
	}
	return
}

// AddChannel adds a channel to a bridge
// Equivalent to Post /bridges/{id}/addChannel
func (b *Bridge) AddChannel(key *ari.Key, channelID string) (err error) {
	id := key.ID

	type request struct {
		ChannelID string `json:"channel"`
		Role      string `json:"role,omitempty"`
	}

	req := request{channelID, ""}
	err = b.client.post("/bridges/"+id+"/addChannel", nil, &req)
	return
}

// RemoveChannel removes the specified channel from a bridge
// Equivalent to Post /bridges/{id}/removeChannel
func (b *Bridge) RemoveChannel(key *ari.Key, channelID string) (err error) {
	id := key.ID

	req := struct {
		ChannelID string `json:"channel"`
	}{
		ChannelID: channelID,
	}

	//pass request
	err = b.client.post("/bridges/"+id+"/removeChannel", nil, &req)
	return
}

// Delete shuts down a bridge. If any channels are in this bridge,
// they will be removed and resume whatever they were doing beforehand.
// This means that the channels themselves are not deleted.
// Equivalent to DELETE /bridges/{id}
func (b *Bridge) Delete(key *ari.Key) (err error) {
	id := key.ID
	err = b.client.del("/bridges/"+id, nil, "")
	return
}

// Play attempts to play the given mediaURI on the bridge, using the playbackID
// as the identifier to the created playback handle
func (b *Bridge) Play(key *ari.Key, playbackID string, mediaURI string) (ph *ari.PlaybackHandle, err error) {
	ph = b.StagePlay(key, playbackID, mediaURI)
	err = ph.Exec()
	return
}

// StagePlay stages a `Play` operation on the bridge
func (b *Bridge) StagePlay(key *ari.Key, playbackID string, mediaURI string) (ph *ari.PlaybackHandle) {
	id := key.ID

	resp := make(map[string]interface{})
	type request struct {
		Media string `json:"media"`
	}
	req := request{mediaURI}
	playbackKey := ari.NewKey(ari.PlaybackKey, playbackID, ari.WithApp(b.client.ApplicationName()), ari.WithNode(b.client.node))
	return ari.NewPlaybackHandle(playbackKey, b.client.Playback(), func(pb *ari.PlaybackHandle) (err error) {
		err = b.client.post("/bridges/"+id+"/play/"+playbackID, &resp, &req)
		return
	})
}

// Record attempts to record audio on the bridge, using name as the identifier for
// the created live recording handle
func (b *Bridge) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle, err error) {
	rh = b.StageRecord(key, name, opts)
	err = rh.Exec()
	return
}

// StageRecord stages a `Record` opreation
func (b *Bridge) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (rh *ari.LiveRecordingHandle) {
	id := key.ID

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

	recordingKey := ari.NewKey(ari.LiveRecordingKey, name, ari.WithApp(b.client.ApplicationName()), ari.WithNode(b.client.node))

	return ari.NewLiveRecordingHandle(recordingKey, b.client.LiveRecording(), func() (err error) {
		err = b.client.post("/bridges/"+id+"/record", &resp, &req)
		return
	})
}

// Subscribe creates an event subscription for events related to the given
// bridge␃entity
func (b *Bridge) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	inSub := b.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		for {
			select {
			case <-outSub.closedChan:
				return
			case e, ok := <-inSub.Events():
				if !ok {
					return
				}
				for _, k := range e.Keys() {
					if k.Match(key) {
						outSub.events <- e
						break
					}
				}
			}
		}
	}()

	return outSub
}
