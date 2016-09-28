package ari

// Bridge represents a communication path to an
// Asterisk server for working with bridge resources
type Bridge interface {

	// Create creates a bridge
	Create(id string, btype string, name string) (*BridgeHandle, error)

	// Get gets the BridgeHandle
	Get(id string) *BridgeHandle

	// Lists returns the lists of bridges in asterisk
	List() ([]*BridgeHandle, error)

	// Data gets the bridge data
	Data(id string) (BridgeData, error)

	// AddChannel adds a channel to the bridge
	AddChannel(bridgeID string, channelID string) error

	// RemoveChannel removes a channel from the bridge
	RemoveChannel(bridgeID string, channelID string) error

	// Delete deletes the bridge
	Delete(id string) error

	// Play plays the media URI to the bridge
	Play(id string, playbackID string, mediaURI string) (*PlaybackHandle, error)

	// Record records the bridge
	Record(id string, name string, opts *RecordingOptions) (*LiveRecordingHandle, error)

	// Subscribe subscribes the given bridge events events
	Subscribe(id string, n ...string) Subscription
}

// BridgeData describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type BridgeData struct {
	ID         string   `json:"id"`          // Unique Id for this bridge
	Class      string   `json:"bridge"`      // Class of the bridge (TODO: huh?)
	Type       string   `json:"bridge_type"` // Type of bridge (mixing, holding, dtmf_events, proxy_media)
	ChannelIDs []string `json:"channels"`    // List of pariticipating channel ids
	Creator    string   `json:"creator"`     // Creating entity of the bridge
	Name       string   `json:"name"`        // The name of the bridge
	Technology string   `json:"technology"`  // Name of the bridging technology
}

// NewBridgeHandle creates a new bridge handle
func NewBridgeHandle(id string, b Bridge) *BridgeHandle {
	return &BridgeHandle{
		id: id,
		b:  b,
	}
}

// BridgeHandle is the handle to a bridge for performing operations
type BridgeHandle struct {
	id string
	b  Bridge
}

// ID returns the identifier for the bridge
func (bh *BridgeHandle) ID() string {
	return bh.id
}

// AddChannel adds a channel to the bridge
func (bh *BridgeHandle) AddChannel(channelID string) (err error) {
	err = bh.b.AddChannel(bh.id, channelID)
	return
}

// RemoveChannel removes a channel from the bridge
func (bh *BridgeHandle) RemoveChannel(channelID string) (err error) {
	err = bh.b.RemoveChannel(bh.id, channelID)
	return
}

// Delete deletes the bridge
func (bh *BridgeHandle) Delete(channelID string) (err error) {
	err = bh.b.Delete(bh.id)
	return
}

// Data gets the bridge data
func (bh *BridgeHandle) Data() (bd BridgeData, err error) {
	bd, err = bh.b.Data(bh.id)
	return
}

// Play initiates playback of the specified media uri
// to the bridge, returning the Playback handle
func (bh *BridgeHandle) Play(id string, mediaURI string) (ph *PlaybackHandle, err error) {
	ph, err = bh.b.Play(bh.id, id, mediaURI)
	return
}

// Record records the bridge to the given filename
func (bh *BridgeHandle) Record(name string, opts *RecordingOptions) (rh *LiveRecordingHandle, err error) {
	rh, err = bh.b.Record(bh.id, name, opts)
	return
}

// Playback returns the playback transport
func (bh *BridgeHandle) Playback() Playback {
	if pb, ok := bh.b.(Playbacker); ok {
		return pb.Playback()
	}
	return nil
}

// Subscribe creates a subscription to the list of events
func (bh *BridgeHandle) Subscribe(n ...string) Subscription {
	if bh == nil {
		return nil
	}
	return bh.b.Subscribe(bh.id, n...)
}

// Match returns true if the event matches the bridge
func (bh *BridgeHandle) Match(e Event) bool {
	bridgeEvent, ok := e.(BridgeEvent)
	if !ok {
		return false
	}
	bridgeIDs := bridgeEvent.GetBridgeIDs()
	for _, i := range bridgeIDs {
		if i == bh.id {
			return true
		}
	}
	return false
}
