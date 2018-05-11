package native

import (
	"errors"
	"time"

	"github.com/CyCoreSystems/ari"
	uuid "github.com/satori/go.uuid"
)

// Bridge provides the ARI Bridge accessors for the native client
type Bridge struct {
	client *Client
}

// Create creates a bridge and returns the lazy handle for the bridge
func (b *Bridge) Create(key *ari.Key, t string, name string) (bh *ari.BridgeHandle, err error) {
	bh, err = b.StageCreate(key, t, name)
	if err != nil {
		return nil, err
	}
	return bh, bh.Exec()
}

// StageCreate creates a new bridge handle, staged with a bridge `Create` operation.
func (b *Bridge) StageCreate(key *ari.Key, btype, name string) (*ari.BridgeHandle, error) {
	if key.ID == "" {
		u, err := uuid.NewV1()
		if err != nil {
			return nil, err
		}
		key.ID = u.String()
	}

	req := struct {
		ID   string `json:"bridgeId,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	}{
		ID:   key.ID,
		Type: btype,
		Name: name,
	}

	return ari.NewBridgeHandle(b.client.stamp(key), b, func(bh *ari.BridgeHandle) (err error) {
		return b.client.post("/bridges/"+key.ID, nil, &req)
	}), nil
}

// Get gets the lazy handle for the given bridge id
func (b *Bridge) Get(key *ari.Key) *ari.BridgeHandle {
	return ari.NewBridgeHandle(b.client.stamp(key), b, nil)
}

// List lists the current bridges and returns a list of lazy handles
func (b *Bridge) List(filter *ari.Key) (bx []*ari.Key, err error) {
	// native client ignores filter

	var bridges = []struct {
		ID string `json:"id"`
	}{}

	err = b.client.get("/bridges", &bridges)
	for _, i := range bridges {
		k := b.client.stamp(ari.NewKey(ari.BridgeKey, i.ID))
		if filter.Match(k) {
			bx = append(bx, k)
		}
	}
	return
}

// Data returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (b *Bridge) Data(key *ari.Key) (*ari.BridgeData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("bridge key not supplied")
	}

	var data = new(ari.BridgeData)
	if err := b.client.get("/bridges/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "bridge", "%v", key.ID)
	}

	data.Key = b.client.stamp(key)

	return data, nil
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
func (b *Bridge) Play(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	h, err := b.StagePlay(key, playbackID, mediaURI)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StagePlay stages a `Play` operation on the bridge
func (b *Bridge) StagePlay(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	if playbackID == "" {
		u, err := uuid.NewV1()
		if err != nil {
			return nil, err
		}
		playbackID = u.String()
	}

	resp := make(map[string]interface{})
	req := struct {
		Media string `json:"media"`
	}{
		Media: mediaURI,
	}
	playbackKey := b.client.stamp(ari.NewKey(ari.PlaybackKey, playbackID))

	return ari.NewPlaybackHandle(playbackKey, b.client.Playback(), func(h *ari.PlaybackHandle) error {
		return b.client.post("/bridges/"+key.ID+"/play/"+playbackID, &resp, &req)
	}), nil
}

// Record attempts to record audio on the bridge, using name as the identifier for
// the created live recording handle
func (b *Bridge) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	h, err := b.StageRecord(key, name, opts)
	if err != nil {
		return nil, err
	}
	return h, h.Exec()
}

// StageRecord stages a `Record` opreation
func (b *Bridge) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	resp := make(map[string]interface{})
	req := struct {
		Name        string `json:"name"`
		Format      string `json:"format"`
		MaxDuration int    `json:"maxDurationSeconds"`
		MaxSilence  int    `json:"maxSilenceSeconds"`
		IfExists    string `json:"ifExists,omitempty"`
		Beep        bool   `json:"beep"`
		TerminateOn string `json:"terminateOn,omitempty"`
	}{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}

	recordingKey := b.client.stamp(ari.NewKey(ari.LiveRecordingKey, name))

	return ari.NewLiveRecordingHandle(recordingKey, b.client.LiveRecording(), func(h *ari.LiveRecordingHandle) error {
		return b.client.post("/bridges/"+key.ID+"/record", &resp, &req)
	}), nil
}

// Subscribe creates an event subscription for events related to the given
// bridge␃entity
func (b *Bridge) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	return b.client.Bus().Subscribe(key, n...)
}
