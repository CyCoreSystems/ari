package ari

// Bridge describes an Asterisk Bridge, the entity which merges media from
// one or more channels into a common audio output
type Bridge struct {
	Bridge_class string   `json:"bridge"`      // Class of the bridge (TODO: huh?)
	Bridge_type  string   `json:"bridge_type"` // Type of bridge (TODO: enumerate types)
	Channels     []string `json:"channels"`    // List of pariticipating channel ids
	Creator      string   `json:"creator"`     // Creating entity of the bridge
	Id           string   `json:"id"`          // Unique Id for this bridge
	Name         string   `json:"name"`        // The name of the bridge
	Technology   string   `json:"technology"`  // Name of the bridging technology
}

//Request structure for creating a bridge. No properies are required, meaning an empty struct may be passed to 'CreateBridge'
type CreateBridgeRequest struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`

	//BridgeId only necessary for the origination of a bridge.
	BridgeId string `json:"bridgeId,omitempty"`
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
	err := c.AriGet("/bridges", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Create a new bridge
//Equivalent to POST /bridges
func (c *Client) CreateBridge(req CreateBridgeRequest) (Bridge, error) {
	var m Bridge

	//send request
	err := c.AriPost("/bridges", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Update a bridge or create a new one (upsert)
//Equivalent to POST /bridges/{bridgeId}
func (c *Client) UpsertBridge(bridgeId string, req CreateBridgeRequest) (Bridge, error) {
	var m Bridge

	//send request
	err := c.AriPost("/bridges/"+bridgeId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Get bridge details
//Equivalent to Get /bridges/{bridgeId}
func (c *Client) GetBridge(bridgeId string) (Bridge, error) {
	var m Bridge
	err := c.AriGet("/bridges/"+bridgeId, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Add a channel to a bridge
//Equivalent to Post /bridges/{bridgeId}/addChannel
func (c *Client) AddChannel(bridgeId string, req AddChannelRequest) error {
	//No return, so no model to create

	//send request, no model so pass nil
	err := c.AriPost("/bridges/"+bridgeId+"/addChannel", nil, &req)
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
	err := c.AriPost("/bridges/"+bridgeId+"/removeChannel", nil, &req)
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
	err := c.AriPost("/bridges/"+bridgeId+"/moh", nil, &req)
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
	err := c.AriPost("/bridges/"+bridgeId+"/play", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Start playback of specific media on specified bridge
//Equivalent to  Post /bridges/{bridgeId}/play/{playbackId}
func (c *Client) PlayToBridgeById(bridgeId string, playbackId string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	err := c.AriPost("/bridges/"+bridgeId+"/play/"+playbackId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//start a recording on specified bridge
//Equivalent to  Post /bridges/{bridgeId}/record
func (c *Client) RecordBridge(bridgeId string, req RecordRequest) (LiveRecording, error) {

	var m LiveRecording

	//send request
	err := c.AriPost("/bridges/"+bridgeId+"/record", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Shut down a bridge. If any channels are in this bridge, they will be removed and resume whatever they were doing beforehand.
//This means that the channels themselves are not deleted.
//Equivalent to DELETE /bridges/{bridgeId}
func (c *Client) BridgeDelete(bridgeId string) error {
	err := c.AriDelete("/bridges/"+bridgeId, nil, nil)
	return err
}

//Stop playing music on hold to a bridge. This will only stop music on hold being played via POST bridges/{bridgeId}/moh.
//Equivalent to DELETE /bridges/{bridgeId}/moh
func (c *Client) BridgeStopMoh(bridgeId string) error {
	err := c.AriDelete("/bridges/"+bridgeId+"/moh", nil, nil)
	return err
}
