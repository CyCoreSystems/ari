package ari

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// PlaybackStartTimeout is the time to allow for Asterisk to
// send the PlaybackStarted before giving up.
var PlaybackStartTimeout = 1 * time.Second

// MaxPlaybackTime is the maximum amount of time to allow for
// a playback to complete.
var MaxPlaybackTime = 10 * time.Minute

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
func (p *Playback) Control(operation string) error {
	if p.client == nil {
		return fmt.Errorf("No client found in Playback")
	}
	return p.client.ControlPlayback(p.Id, operation)
}

// Stop the current Playback.
func (p *Playback) Stop() error {
	if p.client == nil {
		return fmt.Errorf("No client found in Playback")
	}
	return p.client.StopPlayback(p.Id)
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

// A Player is anyhing which can "Play" an mediaUri
type Player interface {
	Play(string) (string, error)
	GetClient() *Client
}

// Play plays audio to the given Player, returning a channel
// which is closed on completion.  If an error occurs, the
// error is sent on the channel first.
func Play(ctx context.Context, p Player, mediaUri string) error {
	c := p.GetClient()
	if c == nil {
		return fmt.Errorf("Failed to find *ari.Client in Player")
	}

	s := c.Bus.Subscribe("PlaybackStarted", "PlaybackFinished")
	defer s.Cancel()

	id, err := p.Play(mediaUri)
	if err != nil {
		return err
	}
	defer c.StopPlayback(id)

	// Wait for the playback to start
	startTimer := time.After(PlaybackStartTimeout)
PlaybackStartLoop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case v := <-s.C:
			if v == nil {
				Logger.Debug("Nil event received")
				continue PlaybackStartLoop
			}
			switch v.GetType() {
			case "PlaybackStarted":
				e := v.(*PlaybackStarted)
				if e.Playback.Id != id {
					Logger.Debug("Ignoring unrelated playback")
					continue PlaybackStartLoop
				}
				Logger.Debug("Playback started")
				break PlaybackStartLoop
			case "PlaybackFinished":
				e := v.(*PlaybackFinished)
				if e.Playback.Id != id {
					Logger.Debug("Ignoring unrelated playback")
					continue PlaybackStartLoop
				}
				Logger.Debug("Playback stopped (before PlaybackStated received)")
				return nil
			default:
				Logger.Debug("Unhandled e.Type", v.GetType())
				continue PlaybackStartLoop
			}
		case <-startTimer:
			Logger.Error("Playback timed out")
			return fmt.Errorf("Timeout waiting for start of playback")
		}
	}

	// Playback has started.  Wait for it to finish
	stopTimer := time.After(MaxPlaybackTime)
PlaybackStopLoop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case v := <-s.C:
			if v == nil {
				Logger.Debug("Nil event received")
				continue PlaybackStopLoop
			}
			switch v.GetType() {
			case "PlaybackFinished":
				e := v.(*PlaybackFinished)
				if e.Playback.Id != id {
					Logger.Debug("Ignoring unrelated playback")
					continue PlaybackStopLoop
				}
				Logger.Debug("Playback stopped")
				return nil
			default:
				Logger.Debug("Unhandled e.Type", v.GetType())
				continue PlaybackStopLoop
			}
		case <-stopTimer:
			Logger.Error("Playback timed out")
			return fmt.Errorf("Timeout waiting for start of playback")
		}
	}
}

// PlaybackQueue represents a sequence of audio playbacks
// which are to be played on the associated Player
type PlaybackQueue struct {
	player Player             // the player (channel or bridge on which the audio is to be played)
	cancel context.CancelFunc // the cancel function for this playback context

	presentPlaybackId string // Id of the currently-playing playback

	queue []string // List of mediaUris to be played
	mu    sync.Mutex
}

func NewPlaybackQueue(p Player) *PlaybackQueue {
	return &PlaybackQueue{player: p}
}
