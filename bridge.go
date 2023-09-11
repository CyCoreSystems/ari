package ari

// Bridge represents a communication path to an
// Asterisk server for working with bridge resources
type Bridge interface {
	// Create creates a bridge
	Create(key *Key, btype string, name string) (*BridgeHandle, error)

	// StageCreate creates a new bridge handle, staged with a bridge `Create` operation.
	StageCreate(key *Key, btype string, name string) (*BridgeHandle, error)

	// Get gets the BridgeHandle
	Get(key *Key) *BridgeHandle

	// Lists returns the lists of bridges in asterisk, optionally using the key for filtering.
	List(*Key) ([]*Key, error)

	// Data gets the bridge data
	Data(key *Key) (*BridgeData, error)

	// AddChannel adds a channel to the bridge
	AddChannel(key *Key, channelID string) error

	// AddChannelWithOptions adds a channel to a bridge, specifying additional options to be applied to that channel
	AddChannelWithOptions(key *Key, channelID string, options *BridgeAddChannelOptions) error

	// RemoveChannel removes a channel from the bridge
	RemoveChannel(key *Key, channelID string) error

	// Delete deletes the bridge
	Delete(key *Key) error

	// MOH plays music on hold
	MOH(key *Key, moh string) error

	// StopMOH stops music on hold
	StopMOH(key *Key) error

	// Play plays the media URI to the bridge
	Play(key *Key, playbackID string, mediaURI ...string) (*PlaybackHandle, error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	StagePlay(key *Key, playbackID string, mediaURI ...string) (*PlaybackHandle, error)

	// Record records the bridge
	Record(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	StageRecord(key *Key, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// Subscribe subscribes the given bridge events events
	Subscribe(key *Key, n ...string) Subscription

	// VideoSource add Channel as Video-Source-ID at bridge
	VideoSource(key *Key, channelID string) error

	// VideoSourceDelete delete Video-Source-ID from bridge
	VideoSourceDelete(key *Key) error
}

// BridgeData describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type BridgeData struct {
	// Key is the cluster-unique identifier for this bridge
	Key *Key `json:"key"`

	ID         string   `json:"id"`           // Unique Id for this bridge
	Class      string   `json:"bridge_class"` // Class of the bridge
	Type       string   `json:"bridge_type"`  // Type of bridge (mixing, holding, dtmf_events, proxy_media)
	ChannelIDs []string `json:"channels"`     // List of pariticipating channel ids
	Creator    string   `json:"creator"`      // Creating entity of the bridge
	Name       string   `json:"name"`         // The name of the bridge
	Technology string   `json:"technology"`   // Name of the bridging technology
}

// BridgeAddChannelOptions describes additional options to be applied to a channel when it is joined to a bridge
type BridgeAddChannelOptions struct {
	// AbsorbDTMF indicates that DTMF coming from this channel will not be passed through to the bridge
	AbsorbDTMF bool

	// Mute indicates that the channel should be muted, preventing audio from it passing through to the bridge
	Mute bool

	// Role indicates the channel's role in the bridge
	Role string
}

// Channels returns the list of channels found in the bridge
func (b *BridgeData) Channels() (list []*Key) {
	for _, id := range b.ChannelIDs {
		list = append(list, b.Key.New(ChannelKey, id))
	}

	return
}

// NewBridgeHandle creates a new bridge handle
func NewBridgeHandle(key *Key, b Bridge, exec func(bh *BridgeHandle) error) *BridgeHandle {
	return &BridgeHandle{
		key:  key,
		b:    b,
		exec: exec,
	}
}

// BridgeHandle is the handle to a bridge for performing operations
type BridgeHandle struct {
	key      *Key
	b        Bridge
	exec     func(bh *BridgeHandle) error
	executed bool
}

// ID returns the identifier for the bridge
func (bh *BridgeHandle) ID() string {
	return bh.key.ID
}

// Key returns the Key of the bridge
func (bh *BridgeHandle) Key() *Key {
	return bh.key
}

// Exec executes any staged operations attached on the bridge handle
func (bh *BridgeHandle) Exec() error {
	if !bh.executed {
		bh.executed = true
		if bh.exec != nil {
			err := bh.exec(bh)
			bh.exec = nil

			return err
		}
	}

	return nil
}

// AddChannel adds a channel to the bridge
func (bh *BridgeHandle) AddChannel(channelID string) error {
	return bh.b.AddChannel(bh.key, channelID)
}

// AddChannelWithOptions adds a channel to the bridge, specifying additional options
func (bh *BridgeHandle) AddChannelWithOptions(channelID string, options *BridgeAddChannelOptions) error {
	return bh.b.AddChannelWithOptions(bh.key, channelID, options)
}

// RemoveChannel removes a channel from the bridge
func (bh *BridgeHandle) RemoveChannel(channelID string) error {
	return bh.b.RemoveChannel(bh.key, channelID)
}

// Delete deletes the bridge
func (bh *BridgeHandle) Delete() (err error) {
	err = bh.b.Delete(bh.key)
	return
}

// Data gets the bridge data
func (bh *BridgeHandle) Data() (*BridgeData, error) {
	return bh.b.Data(bh.key)
}

// MOH requests that the given MusicOnHold class being played to the bridge
func (bh *BridgeHandle) MOH(class string) error {
	return bh.b.MOH(bh.key, class)
}

// StopMOH requests that any MusicOnHold which is being played to the bridge is stopped.
func (bh *BridgeHandle) StopMOH() error {
	return bh.b.StopMOH(bh.key)
}

// Play initiates playback of the specified media uri
// to the bridge, returning the Playback handle
func (bh *BridgeHandle) Play(id string, mediaURI ...string) (*PlaybackHandle, error) {
	return bh.b.Play(bh.key, id, mediaURI...)
}

// StagePlay stages a `Play` operation.
func (bh *BridgeHandle) StagePlay(id string, mediaURI ...string) (*PlaybackHandle, error) {
	return bh.b.StagePlay(bh.key, id, mediaURI...)
}

// Record records the bridge to the given filename
func (bh *BridgeHandle) Record(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
	return bh.b.Record(bh.key, name, opts)
}

// StageRecord stages a `Record` operation
func (bh *BridgeHandle) StageRecord(name string, opts *RecordingOptions) (*LiveRecordingHandle, error) {
	return bh.b.StageRecord(bh.key, name, opts)
}

// Subscribe creates a subscription to the list of events
func (bh *BridgeHandle) Subscribe(n ...string) Subscription {
	if bh == nil {
		return nil
	}

	return bh.b.Subscribe(bh.key, n...)
}

// VideoSource sets channel as Video-Source-Id in a multi-party mixing bridge
func (bh *BridgeHandle) VideoSource(channelID string) error {
	return bh.b.VideoSource(bh.key, channelID)
}

// VideoSourceDelete deletes Video-Source-Id in a multi-party mixing bridge
func (bh *BridgeHandle) VideoSourceDelete() error {
	return bh.b.VideoSourceDelete(bh.key)
}
