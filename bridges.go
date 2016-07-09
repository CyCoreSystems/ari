package ari

import (
	"fmt"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Bridge describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type Bridge struct {
	ID         string   `json:"id"`          // Unique Id for this bridge
	Class      string   `json:"bridge"`      // Class of the bridge (TODO: huh?)
	Type       string   `json:"bridge_type"` // Type of bridge (mixing, holding, dtmf_events, proxy_media)
	ChannelIDs []string `json:"channels"`    // List of pariticipating channel ids
	Creator    string   `json:"creator"`     // Creating entity of the bridge
	Name       string   `json:"name"`        // The name of the bridge
	Technology string   `json:"technology"`  // Name of the bridging technology

	client *Client // Reference to the client which created or returned this bridge
}

// Add adds a channel to the bridge
func (b *Bridge) Add(channelID string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.AddChannel(b.ID, AddChannelRequest{ChannelID: channelID})
}

// Remove removes a channel from the bridge
func (b *Bridge) Remove(channelID string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.RemoveChannel(b.ID, channelID)
}

// Delete destroys the bridge
func (b *Bridge) Delete() error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.BridgeDelete(b.ID)
}

// Play plays an audio uri to a bridge, returning its playback ID
func (b *Bridge) Play(mediaURI string) (string, error) {
	id := uuid.NewV1().String()
	err := b.PlayWithID(id, mediaURI)
	return id, err
}

// PlayWithID plays an audio uri to the bridge with the provided playback ID
func (b *Bridge) PlayWithID(id, mediaURI string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}

	_, err := b.client.PlayToBridgeByID(b.ID, id, PlayMediaRequest{Media: mediaURI})
	return err
}

// Record starts recording the audio from a bridge to the given recording name.
func (b *Bridge) Record(name string, opts *RecordingOptions) (*LiveRecording, error) {
	if opts == nil {
		opts = &RecordingOptions{}
	}
	return b.GetClient().RecordBridge(b.ID, opts.ToRequest(name))
}

// AttachClient attaches the provided ARI client to the bridge
func (b *Bridge) AttachClient(a *Client) {
	b.client = a
}

// GetClient returns the ARI client which created the bridge
func (b *Bridge) GetClient() *Client {
	return b.client
}

// GetID returns the ID of this bridge
func (b *Bridge) GetID() string {
	return b.ID
}

// CreateBridgeRequest is the structure for creating a bridge. No properies are required, meaning an empty struct may be passed to 'CreateBridge'
type CreateBridgeRequest struct {
	ID   string `json:"bridgeId,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// AddChannelRequest is the structure to add a channel to a bridge. Only Channel is required.
// ChannelID field allows for comma-separated-values to add multiple channels.
type AddChannelRequest struct {
	ChannelID string `json:"channel"`
	Role      string `json:"role,omitempty"`
}

// ListBridges returns all active bridges in Asterisk
// Equivalent to GET /bridges
func (c *Client) ListBridges() ([]Bridge, error) {
	var m []Bridge
	err := c.Get("/bridges", &m)
	if err != nil {
		return m, err
	}

	// Attach the client to each bridge
	for _, b := range m {
		b.client = c
	}

	return m, nil
}

// NewBridge is a simple wrapper to create a new,
// unique bridge, with the default options
func (c *Client) NewBridge() (Bridge, error) {
	id := uuid.NewV1().String()
	return c.UpsertBridge(id, CreateBridgeRequest{ID: id})
}

// NewBridgeWithID is a simple wrapper to create a new,
// unique bridge, with the default options
func (c *Client) NewBridgeWithID(id string) (Bridge, error) {
	return c.UpsertBridge(id, CreateBridgeRequest{ID: id})
}

// CreateBridge creates a new bridge
//Equivalent to POST /bridges
func (c *Client) CreateBridge(req CreateBridgeRequest) (Bridge, error) {
	var m Bridge

	//send request
	err := c.Post("/bridges", &m, &req)
	if err != nil {
		return m, err
	}

	// Attach the client to the bridge
	m.client = c

	return m, nil
}

// UpsertBridge adds or updates a bridge
// Equivalent to POST /bridges/{bridgeId}
func (c *Client) UpsertBridge(id string, req CreateBridgeRequest) (Bridge, error) {
	var m Bridge

	//send request
	err := c.Post("/bridges/"+id, &m, &req)
	if err != nil {
		return m, err
	}

	// Attach the client to the bridge
	m.client = c

	return m, nil
}

// GetBridge returns the details of a bridge
// Equivalent to Get /bridges/{bridgeId}
func (c *Client) GetBridge(id string) (Bridge, error) {
	var m Bridge
	err := c.Get("/bridges/"+id, &m)
	if err != nil {
		return m, err
	}

	// Attach the client to the bridge
	m.client = c

	return m, nil
}

// AddChannel adds a channel to a bridge
// Equivalent to Post /bridges/{id}/addChannel
func (c *Client) AddChannel(id string, req AddChannelRequest) error {
	//No return, so no model to create

	//send request, no model so pass nil
	err := c.Post("/bridges/"+id+"/addChannel", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

// RemoveChannel removes the specified channel from a bridge
// Equivalent to Post /bridges/{id}/removeChannel
func (c *Client) RemoveChannel(id string, channelID string) error {
	req := struct {
		ChannelID string `json:"channel"`
	}{
		ChannelID: channelID,
	}

	//pass request
	err := c.Post("/bridges/"+id+"/removeChannel", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

// PlayMusicOnHold plays a music on hold class to a bridge or changes the existing MOH class
// Equivalent to  Post /bridges/{bridgeId}/moh (music on hold)
func (c *Client) PlayMusicOnHold(id string, class string) error {

	req := struct {
		Class string `json:"mohClass,omitempty"`
	}{
		Class: class,
	}

	//send request
	err := c.Post("/bridges/"+id+"/moh", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

// PlayToBridge starts playback of media on specified bridge
//Equivalent to  Post /bridges/{id}/play
func (c *Client) PlayToBridge(id string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	//send request
	err := c.Post("/bridges/"+id+"/play", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

// PlayToBridgeByID starts playback of specific media on specified bridge
//Equivalent to  Post /bridges/{id}/play/{playbackID}
func (c *Client) PlayToBridgeByID(id string, playbackID string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	err := c.Post("/bridges/"+id+"/play/"+playbackID, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

// RecordBridge starts a recording on specified bridge
//Equivalent to  Post /bridges/{id}/record
func (c *Client) RecordBridge(id string, req *RecordRequest) (*LiveRecording, error) {
	var m LiveRecording

	//send request
	err := c.Post("/bridges/"+id+"/record", &m, &req)
	return &m, err
}

// BridgeDelete shuts down a bridge. If any channels are in this bridge, they will be removed and resume whatever they were doing beforehand.
//This means that the channels themselves are not deleted.
//Equivalent to DELETE /bridges/{id}
func (c *Client) BridgeDelete(id string) error {
	err := c.Delete("/bridges/"+id, nil, "")
	return err
}

// BridgeStopMoh stops playing music on hold to a bridge. This will only stop music on hold being played via POST bridges/{id}/moh.
// Equivalent to DELETE /bridges/{id}/moh
func (c *Client) BridgeStopMoh(id string) error {
	err := c.Delete("/bridges/"+id+"/moh", nil, "")
	return err
}

//
//  Context-related items
//

// bridgeKey is the key type for contexts
type bridgeKey string

// NewBridgeContext returns a context with the bridge attached
func NewBridgeContext(ctx context.Context, b *Bridge) context.Context {
	return NewBridgeContextWithKey(ctx, b, "_default")
}

// NewBridgeContextWithKey returns a context with the bridge attached
// as the given key
func NewBridgeContextWithKey(ctx context.Context, b *Bridge, name string) context.Context {
	return context.WithValue(ctx, bridgeKey(name), b)
}

// BridgeFromContext returns the default Bridge stored in the context
func BridgeFromContext(ctx context.Context) (*Bridge, bool) {
	return BridgeFromContextWithKey(ctx, "_default")
}

// BridgeFromContextWithKey returns the Bridge stored in the context
func BridgeFromContextWithKey(ctx context.Context, name string) (*Bridge, bool) {
	c, ok := ctx.Value(bridgeKey(name)).(*Bridge)
	return c, ok
}
