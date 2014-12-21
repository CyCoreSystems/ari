package ari

import "fmt"

// Playback describes a session of playing media to a channel
// MediaUri is of the form 'type:name', where type can be one of:
//  - sound : a Sound on the Asterisk system
//  - recording : a StoredRecording on the Asterisk system
//  - number : a number, to be spoken
//  - digits : a set of digits, to be spoken
//  - characters : a set of characters, to be spoken
//  - tone : a tone sequence, which may optionally take a tonezone parameter (e.g, tone:ring:tonezone=fr)
//
// TargetUri is of the form 'type:id', and looks like the following two options:
//  - bridge:bridgeId
//  - channel:channelId

type Playback struct {
	Id        string `json:"id"` // Unique Id for this playback session
	Language  string `json:"language,omitempty"`
	MediaUri  string `json:"media_uri"`  // URI for the media which is to be played
	State     string `json:"state"`      // State of the playback operation
	TargetUri string `json:"target_uri"` // URI of the channel or bridge on which the media should be played (follows format of 'type':'name')

	client *Client // Reference to the client which created or returned this channel
}

//Get a playback's details
//Equivalent to GET /playbacks/{playbackId}
func (c *Client) GetPlaybackDetails(playbackId string) (Playback, error) {
	var m Playback
	err := c.AriGet("/playbacks/"+playbackId, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// Control the current Playback
func (p *Playback) Control(playbackId string, operation string) error {
	if p.client == nil {
		return fmt.Errorf("No client found in Playback")
	}
	return p.client.ControlPlayback(playbackId, operation)
}

// Stop the current Playback.
func (p *Playback) Stop(playbackId string) error {
	if p.client == nil {
		return fmt.Errorf("No client found in Playback")
	}
	return p.client.StopPlayback(playbackId)
}

//Equivalent to POST /playbacks/{playbackId}/control
func (c *Client) ControlPlayback(playbackId string, operation string) error {

	//Request structure for controlling playback. Operation is required.
	type request struct {
		Operation string `json:"operation"`
	}

	req := request{operation}

	//Make the request
	err := c.AriPost("/playbacks/"+playbackId+"/control", nil, &req)

	if err != nil {
		return err
	}
	return nil
}

//Stop a playback.
//Equivalent to DELETE /playbacks/{playbackId}
func (c *Client) StopPlayback(playbackId string) error {
	err := c.AriDelete("/playbacks/"+playbackId, nil, nil)
	return err
}
