package ari

// Bridge represents a communication path to an
// Asterisk server for working with bridge resources
type Bridge interface {

	// Create creates a bridge
	Create(id string, btype string, name string) (BridgeHandle, error)

	// StageCreate creates a new bridge handle, staged with a bridge `Create` operation.
	StageCreate(id string, btype string, name string) BridgeHandle

	// Get gets the BridgeHandle
	Get(id string) BridgeHandle

	// Lists returns the lists of bridges in asterisk
	List() ([]BridgeHandle, error)

	// Data gets the bridge data
	Data(id string) (*BridgeData, error)

	// AddChannel adds a channel to the bridge
	AddChannel(bridgeID string, channelID string) error

	// RemoveChannel removes a channel from the bridge
	RemoveChannel(bridgeID string, channelID string) error

	// Delete deletes the bridge
	Delete(id string) error

	// Play plays the media URI to the bridge
	Play(id string, playbackID string, mediaURI string) (PlaybackHandle, error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	StagePlay(id string, playbackID string, mediaURI string) (ph PlaybackHandle)

	// Record records the bridge
	Record(id string, name string, opts *RecordingOptions) (LiveRecordingHandle, error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	StageRecord(id string, name string, opts *RecordingOptions) (rh LiveRecordingHandle)

	// Subscribe subscribes the given bridge events events
	Subscribe(id string, n ...string) Subscription
}

// BridgeData describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type BridgeData struct {
	ID         string   `json:"id"`           // Unique Id for this bridge
	Class      string   `json:"bridge_class"` // Class of the bridge
	Type       string   `json:"bridge_type"`  // Type of bridge (mixing, holding, dtmf_events, proxy_media)
	ChannelIDs []string `json:"channels"`     // List of pariticipating channel ids
	Creator    string   `json:"creator"`      // Creating entity of the bridge
	Name       string   `json:"name"`         // The name of the bridge
	Technology string   `json:"technology"`   // Name of the bridging technology
}

// BridgeHandle is the handle to a bridge for performing operations on a specific bridge
type BridgeHandle interface {
	// ID returns the identifier for the bridge
	ID() string

	// AddChannel adds a channel to the bridge
	AddChannel(channelID string) (err error)

	// RemoveChannel removes a channel from the bridge
	RemoveChannel(channelID string) (err error)

	// Delete deletes the bridge
	Delete() (err error)

	// Data gets the bridge data
	Data() (bd *BridgeData, err error)

	// Play initiates playback of the specified media uri
	// to the bridge, returning the Playback handle
	Play(id string, mediaURI string) (ph PlaybackHandle, err error)

	// StagePlay stages a `Play` operation and returns the `PlaybackHandle`
	// for invoking it.
	StagePlay(id string, mediaURI string) (ph PlaybackHandle)

	// Record records the bridge to the given filename
	Record(name string, opts *RecordingOptions) (rh LiveRecordingHandle, err error)

	// StageRecord stages a `Record` operation and returns the `PlaybackHandle`
	// for invoking it.
	StageRecord(name string, opts *RecordingOptions) (rh LiveRecordingHandle)

	// Subscribe creates a subscription to the list of events
	Subscribe(n ...string) Subscription

	// Match returns true if the event matches the bridge
	Match(e Event) bool

	// Exec executes any staged operations attached on the bridge handle
	Exec() error
}
