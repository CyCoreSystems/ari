package ari

import (
	"fmt"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Bridge describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type Bridge struct {
	Id           string   `json:"id"`          // Unique Id for this bridge
	Bridge_class string   `json:"bridge"`      // Class of the bridge (TODO: huh?)
	Bridge_type  string   `json:"bridge_type"` // Type of bridge (mixing, holding, dtmf_events, proxy_media)
	Channels     []string `json:"channels"`    // List of pariticipating channel ids
	Creator      string   `json:"creator"`     // Creating entity of the bridge
	Name         string   `json:"name"`        // The name of the bridge
	Technology   string   `json:"technology"`  // Name of the bridging technology

	client *Client // Reference to the client which created or returned this bridge
}

// Add adds a channel to the bridge
func (b *Bridge) Add(channelId string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.AddChannel(b.Id, AddChannelRequest{ChannelId: channelId})
}

// Remove removes a channel from the bridge
func (b *Bridge) Remove(channelId string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.RemoveChannel(b.Id, channelId)
}

// Delete destroys the bridge
func (b *Bridge) Delete() error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}
	return b.client.BridgeDelete(b.Id)
}

// Play plays an audio uri to a bridge, returning its playback ID
func (b *Bridge) Play(mediaUri string) (string, error) {
	id := uuid.NewV1().String()
	err := b.PlayWithID(id, mediaUri)
	return id, err
}

// PlayWithID plays an audio uri to the bridge with the provided playback ID
func (b *Bridge) PlayWithID(id, mediaUri string) error {
	if b.client == nil {
		return fmt.Errorf("No client found in Bridge")
	}

	_, err := b.client.PlayToBridgeById(b.Id, id, PlayMediaRequest{Media: mediaUri})
	return err
}

// Record starts recording the audio from a bridge to the given recording name.
func (b *Bridge) Record(name string, opts *RecordingOptions) (*LiveRecording, error) {
	if opts == nil {
		opts = &RecordingOptions{}
	}
	return b.GetClient().RecordBridge(b.Id, opts.ToRequest(name))
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
	return b.Id
}

//Request structure for creating a bridge. No properies are required, meaning an empty struct may be passed to 'CreateBridge'
type CreateBridgeRequest struct {
	Id   string `json:"bridgeId,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

//Request structure to add a channel to a bridge. Only Channel is required.
//Channel field allows for comma-separated-values to add multiple channels.
type AddChannelRequest struct {
	ChannelId string `json:"channel"`
	Role      string `json:"role,omitempty"`
}

//List all active bridges in Asterisk
//Equivalent to GET /bridges
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
	return c.UpsertBridge(id, CreateBridgeRequest{Id: id})
}

// NewBridgeWithId is a simple wrapper to create a new,
// unique bridge, with the default options
func (c *Client) NewBridgeWithId(id string) (Bridge, error) {
	return c.UpsertBridge(id, CreateBridgeRequest{Id: id})
}

//Create a new bridge
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

//Update a bridge or create a new one (upsert)
//Equivalent to POST /bridges/{bridgeId}
func (c *Client) UpsertBridge(bridgeId string, req CreateBridgeRequest) (Bridge, error) {
	var m Bridge

	//send request
	err := c.Post("/bridges/"+bridgeId, &m, &req)
	if err != nil {
		return m, err
	}

	// Attach the client to the bridge
	m.client = c

	return m, nil
}

//Get bridge details
//Equivalent to Get /bridges/{bridgeId}
func (c *Client) GetBridge(bridgeId string) (Bridge, error) {
	var m Bridge
	err := c.Get("/bridges/"+bridgeId, &m)
	if err != nil {
		return m, err
	}

	// Attach the client to the bridge
	m.client = c

	return m, nil
}

//Add a channel to a bridge
//Equivalent to Post /bridges/{bridgeId}/addChannel
func (c *Client) AddChannel(bridgeId string, req AddChannelRequest) error {
	//No return, so no model to create

	//send request, no model so pass nil
	err := c.Post("/bridges/"+bridgeId+"/addChannel", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Remove a specific channel from a bridge
//Equivalent to Post /bridges/{bridgeId}/removeChannel
func (c *Client) RemoveChannel(bridgeId string, channelId string) error {
	//Request structure to remove a channel from a bridge. Channel is required.
	type request struct {
		ChannelId string `json:"channel"`
	}

	req := request{channelId}

	//pass request
	err := c.Post("/bridges/"+bridgeId+"/removeChannel", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Play music on hold to a bridge or change the MOH class that's playing
//Equivalent to  Post /bridges/{bridgeId}/moh (music on hold)
func (c *Client) PlayMusicOnHold(bridgeId string, mohClass string) error {

	//Request structure for playing music on hold to a bridge. MohClass is _not_ required.
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}

	req := request{mohClass}

	//send request
	err := c.Post("/bridges/"+bridgeId+"/moh", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Start playback of media on specified bridge
//Equivalent to  Post /bridges/{bridgeId}/play
func (c *Client) PlayToBridge(bridgeId string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	//send request
	err := c.Post("/bridges/"+bridgeId+"/play", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Start playback of specific media on specified bridge
//Equivalent to  Post /bridges/{bridgeId}/play/{playbackId}
func (c *Client) PlayToBridgeById(bridgeId string, playbackId string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	err := c.Post("/bridges/"+bridgeId+"/play/"+playbackId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//start a recording on specified bridge
//Equivalent to  Post /bridges/{bridgeId}/record
func (c *Client) RecordBridge(bridgeId string, req *RecordRequest) (*LiveRecording, error) {
	var m LiveRecording

	//send request
	err := c.Post("/bridges/"+bridgeId+"/record", &m, &req)
	return &m, err
}

//Shut down a bridge. If any channels are in this bridge, they will be removed and resume whatever they were doing beforehand.
//This means that the channels themselves are not deleted.
//Equivalent to DELETE /bridges/{bridgeId}
func (c *Client) BridgeDelete(bridgeId string) error {
	err := c.Delete("/bridges/"+bridgeId, nil, "")
	return err
}

//Stop playing music on hold to a bridge. This will only stop music on hold being played via POST bridges/{bridgeId}/moh.
//Equivalent to DELETE /bridges/{bridgeId}/moh
func (c *Client) BridgeStopMoh(bridgeId string) error {
	err := c.Delete("/bridges/"+bridgeId+"/moh", nil, "")
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
