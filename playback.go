package ari

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// AllDTMF is a string which contains all possible
// DTMF digits.
const AllDTMF = "0123456789ABCD*#"

// PlaybackStartTimeout is the time to allow for Asterisk to
// send the PlaybackStarted before giving up.
var PlaybackStartTimeout = 10 * time.Second

// MaxPlaybackTime is the maximum amount of time to allow for
// a playback to complete.
var MaxPlaybackTime = 10 * time.Minute

// Playback describes a session of playing media to a channel
// MediaURI is of the form 'type:name', where type can be one of:
//  - sound : a Sound on the Asterisk system
//  - recording : a StoredRecording on the Asterisk system
//  - number : a number, to be spoken (integers, positive or negative)
//  - digits : a set of digits, to be spoken (includes -*#0123456789)
//  - characters : a set of characters, to be spoken
//  - tone : a tone sequence, which may optionally take a tonezone parameter (e.g, tone:ring:tonezone=fr)
//
// TargetURI is of the form 'type:id', and looks like the following two options:
//  - bridge:bridgeID
//  - channel:channelID

// Playback describes an ARI playback handle
type Playback struct {
	ID        string `json:"id"` // Unique ID for this playback session
	Language  string `json:"language,omitempty"`
	MediaURI  string `json:"media_uri"`  // URI for the media which is to be played
	State     string `json:"state"`      // State of the playback operation
	TargetURI string `json:"target_uri"` // URI of the channel or bridge on which the media should be played (follows format of 'type':'name')

	client *Client // Reference to the client which created or returned this channel
}

// PlaybackOptions describes various options which
// are available to playback operations.
type PlaybackOptions struct {
	// ID is an optional ID to use for the playback's ID. If one
	// is not supplied, an ID will be randomly generated internally.
	// NOTE that this ID will only be used for the FIRST playback
	// in a queue.  All subsequent playback IDs will be randomly generated.
	ID string

	// DTMF is an optional channel for received DTMF tones received during the playback.
	// This channel will NOT be closed by the playback.
	DTMF chan<- *ChannelDtmfReceived

	// ExitOnDTMF defines a list of DTMF digits on receipt of which will
	// terminate the playback of the queue.  You may set this to AllDTMF
	// in order to match any DTMF digit.
	ExitOnDTMF string

	// Done is an optional channel for receiving notification when the playback
	// is complete.  This is useful if the playback is to be executed asynchronously.
	// This channel will be closed by the playback when playback the is complete.
	Done chan<- struct{}
}

// GetPlaybackDetails returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (c *Client) GetPlaybackDetails(playbackID string) (Playback, error) {
	var m Playback
	err := c.Get("/playbacks/"+playbackID, &m)
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
	return p.client.ControlPlayback(p.ID, operation)
}

// Stop the current Playback.
func (p *Playback) Stop() error {
	if p.client == nil {
		return fmt.Errorf("No client found in Playback")
	}
	return p.client.StopPlayback(p.ID)
}

// ControlPlayback allows the user to manipulate an in-process playback.
// TODO: list available operations.
// (Equivalent to POST /playbacks/{playbackID}/control)
func (c *Client) ControlPlayback(playbackID string, operation string) error {

	//Request structure for controlling playback. Operation is required.
	type request struct {
		Operation string `json:"operation"`
	}

	req := request{operation}

	//Make the request
	err := c.Post("/playbacks/"+playbackID+"/control", nil, &req)

	if err != nil {
		return err
	}
	return nil
}

// StopPlayback stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (c *Client) StopPlayback(playbackID string) error {
	err := c.Delete("/playbacks/"+playbackID, nil, "")
	return err
}

// A Player is anyhing which can "Play" an mediaURI
type Player interface {
	Play(string) (string, error)
	GetClient() *Client
}

// Play plays audio to the given Player, waiting for completion
// and returning any error encountered during playback.
func Play(ctx context.Context, p Player, mediaURI string) error {
	c := p.GetClient()
	if c == nil {
		return fmt.Errorf("Failed to find *ari.Client in Player")
	}

	s := c.Bus.Subscribe("PlaybackStarted", "PlaybackFinished")
	defer s.Cancel()

	id, err := p.Play(mediaURI)
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
				if e.Playback.ID != id {
					Logger.Debug("Ignoring unrelated playback")
					continue PlaybackStartLoop
				}
				Logger.Debug("Playback started")
				break PlaybackStartLoop
			case "PlaybackFinished":
				e := v.(*PlaybackFinished)
				if e.Playback.ID != id {
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
				if e.Playback.ID != id {
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
			return fmt.Errorf("Timeout waiting for stop of playback")
		}
	}
}

// PlaybackQueue represents a sequence of audio playbacks
// which are to be played on the associated Player
type PlaybackQueue struct {
	queue []string // List of mediaURI to be played
	mu    sync.Mutex

	receivedDTMF string // Storage for received DTMF, if we are listening for them
}

// NewPlaybackQueue creates (but does not start) a new playback queue.
func NewPlaybackQueue() *PlaybackQueue {
	return &PlaybackQueue{}
}

// Add appends one or more mediaURIs to the playback queue
func (pq *PlaybackQueue) Add(mediaURIs ...string) {
	// Make sure our queue exists
	pq.mu.Lock()
	if pq.queue == nil {
		pq.queue = []string{}
	}
	pq.mu.Unlock()

	// Add each media URI to the queue
	for _, u := range mediaURIs {
		if u == "" {
			continue
		}

		pq.mu.Lock()
		pq.queue = append(pq.queue, u)
		pq.mu.Unlock()
	}
}

// Flush empties a playback queue.
// NOTE that this does NOT stop the current playback.
func (pq *PlaybackQueue) Flush() {
	pq.mu.Lock()
	pq.queue = []string{}
	pq.mu.Unlock()
}

// ReceivedDTMF returns any DTMF which has been received
// by the PlaybackQueue.
func (pq *PlaybackQueue) ReceivedDTMF() string {
	return pq.receivedDTMF
}

// Play starts the playback of the queue to the Player.
func (pq *PlaybackQueue) Play(ctx context.Context, p Player, opts *PlaybackOptions) error {
	if opts == nil {
		opts = &PlaybackOptions{}
	}

	// Handle any options we were given
	if opts != nil {
		// Close the done channel when we finish,
		// if we were given one.
		if opts.Done != nil {
			defer close(opts.Done)
		}

		// Listen for DTMF, if we were asked to do so
		if opts.DTMF != nil {
			go func() {
				dtmfSub := p.GetClient().Bus.Subscribe("ChannelDtmfReceived")
				defer dtmfSub.Cancel()
				for {
					select {
					case <-ctx.Done():
						return
					case e := <-dtmfSub.C:
						opts.DTMF <- e.(*ChannelDtmfReceived)
					}
				}
			}()
		}
	}

	// Record any DTMF (this is separate from opts.DTMF) so that we can
	//  - Service ReceivedDTMF requests
	//  - Exit if we were given an ExitOnDTMF list
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		dtmfSub := p.GetClient().Bus.Subscribe("ChannelDtmfReceived")
		defer dtmfSub.Cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-dtmfSub.C:
				if e == nil {
					return
				}
				digit := e.(*ChannelDtmfReceived).Digit
				pq.receivedDTMF += digit
				if strings.Contains(opts.ExitOnDTMF, digit) {
					cancel()
				}
			}
		}
	}()

	// Start the playback
	for i := 0; len(pq.queue) > i; i++ {
		// Make sure our context isn't closed
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		// Get the next clip
		err := Play(ctx, p, pq.queue[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// PlayAsync plays the queue, returing immediately with an error channel,
// which will pass any errors and be closed on completion of the queue.
func (pq *PlaybackQueue) PlayAsync(ctx context.Context, p Player, opts *PlaybackOptions) chan error {
	errChan := make(chan error)
	go func() {
		err := pq.Play(ctx, p, opts)
		if err != nil {
			errChan <- err
		}
		close(errChan)
		return
	}()

	return errChan
}

// IsOpenPattern determines whether the regular expression is
// open-ended (allows for an indeterminite number of trailing
// parts) or not.
func IsOpenPattern(p string) bool {
	strings.TrimSuffix(p, "$")
	if strings.HasSuffix(p, ".") {
		return true
	}
	if strings.HasSuffix(p, "*") {
		return true
	}
	return false
}
