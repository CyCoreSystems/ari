package native

import (
	"errors"
	"time"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/rid"
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
		key.ID = rid.New(rid.Bridge)
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
	bridges := []struct {
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

	data := new(ari.BridgeData)
	if err := b.client.get("/bridges/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "bridge", "%v", key.ID)
	}

	data.Key = b.client.stamp(key)

	return data, nil
}

// AddChannel adds a channel to a bridge
// Equivalent to Post /bridges/{id}/addChannel
func (b *Bridge) AddChannel(key *ari.Key, channelID string) (err error) {
	return b.AddChannelWithOptions(key, channelID, nil)
}

// AddChannelWithOptions adds a channel to a bridge, specifying additional options to be applied to that channel
func (b *Bridge) AddChannelWithOptions(key *ari.Key, channelID string, options *ari.BridgeAddChannelOptions) error {
	if options == nil {
		options = new(ari.BridgeAddChannelOptions)
	}

	req := struct {
		AbsorbDTMF bool   `json:"absorbDTMF,omitempty"`
		ChannelID  string `json:"channel"`
		Mute       bool   `json:"mute,omitempty"`
		Role       string `json:"role,omitempty"`
	}{
		AbsorbDTMF: options.AbsorbDTMF,
		ChannelID:  channelID,
		Mute:       options.Mute,
		Role:       options.Role,
	}

	return b.client.post("/bridges/"+key.ID+"/addChannel", nil, &req)
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

	// pass request
	err = b.client.post("/bridges/"+id+"/removeChannel", nil, &req)

	return
}

// Delete shuts down a bridge. If any channels are in this bridge,
// they will be removed and resume whatever they were doing beforehand.
// This means that the channels themselves are not deleted.
// Equivalent to DELETE /bridges/{id}
func (b *Bridge) Delete(key *ari.Key) (err error) {
	return b.client.del("/bridges/"+key.ID, nil, "")
}

// MOH requests that the given musiconhold class be played to the bridge
func (b *Bridge) MOH(key *ari.Key, class string) error {
	req := struct {
		Class string `json:"mohClass"`
	}{
		Class: class,
	}

	return b.client.post("/bridges/"+key.ID+"/moh", nil, &req)
}

// StopMOH requests that any MusicOnHold which is playing to the bridge be stopped.
func (b *Bridge) StopMOH(key *ari.Key) error {
	return b.client.del("/bridges/"+key.ID+"/moh", nil, "")
}

// Play attempts to play the given mediaURI on the bridge, using the playbackID
// as the identifier to the created playback handle
func (b *Bridge) Play(key *ari.Key, playbackID string, opts interface{}) (*ari.PlaybackHandle, error) {
	if playbackID == "" {
		playbackID = rid.New(rid.Playback)
	}

	h, err := b.StagePlay(key, playbackID, opts)
	if err != nil {
		return nil, err
	}

	return h, h.Exec()
}

// StagePlay stages a `Play` operation on the bridge
func (b *Bridge) StagePlay(key *ari.Key, playbackID string, opts interface{}) (*ari.PlaybackHandle, error) {
	if playbackID == "" {
		playbackID = rid.New(rid.Playback)
	}

	resp := make(map[string]interface{})

	var req interface{}
	switch v := opts.(type) {
	case string:
		req = struct {
			Media string `json:"media"`
		}{
			Media: v,
		}
	case ari.PlaybackOptions:
		req = v
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

// VideoSource sets a channel as the video source in a multi-party mixing bridge.
// This operation has no effect on bridges with two or fewer participants.
// Equivalent to POST /bridges/{bridgeId}/videoSource/{channelId}
func (b *Bridge) VideoSource(key *ari.Key, channelID string) error {
	return b.client.post("/bridges/"+key.ID+"/videoSource/"+channelID, nil, nil)
}

// VideoSourceDelete removes any explicit video source in a multi-party mixing bridge.
// This operation has no effect on bridges with two or fewer participants.
// When no explicit video source is set, talk detection will be used to determine the active video stream.
// Equivalent to DELETE /bridges/{bridgeId}/videoSource
func (b *Bridge) VideoSourceDelete(key *ari.Key) error {
	return b.client.del("/bridges/"+key.ID+"/videoSource", nil, "")
}
