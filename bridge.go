package ari

import "github.com/satori/go.uuid"

// Bridge represents a communication path to an
// Asterisk server for working with bridge resources
type Bridge interface {

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
// to the bridge, returning the Playback's ID
func (bh *BridgeHandle) Play(mediaURI string) (ph *PlaybackHandle, err error) {
	id := uuid.NewV1().String()
	ph, err = bh.b.Play(bh.id, id, mediaURI)
	return
}
