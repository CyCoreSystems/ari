package audio

import (
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"

	"golang.org/x/net/context"
)

// AllDTMF is a string which contains all possible
// DTMF digits.
const AllDTMF = "0123456789ABCD*#"

// PlaybackStartTimeout is the time to allow for Asterisk to
// send the PlaybackStarted before giving up.
var PlaybackStartTimeout = 1 * time.Second

// MaxPlaybackTime is the maximum amount of time to allow for
// a playback to complete.
var MaxPlaybackTime = 10 * time.Minute

// Play plays audio to the given Player, waiting for completion
// and returning any error encountered during playback.
func Play(ctx context.Context, bus ari.Bus, p Player, mediaURI string) error {

	s := bus.Subscribe("PlaybackStarted", "PlaybackFinished")
	defer s.Cancel()

	h, err := p.Play(mediaURI)
	if err != nil {
		return err
	}
	defer h.Stop()

	// NOTE: this is where we may want to be able to access handle.ID directly?
	data, err := h.Data()
	if err != nil {
		return err
	}

	id := data.ID

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
				e := v.(*v2.PlaybackStarted)
				if e.Playback.ID != id {
					Logger.Debug("Ignoring unrelated playback")
					continue PlaybackStartLoop
				}
				Logger.Debug("Playback started")
				break PlaybackStartLoop
			case "PlaybackFinished":
				e := v.(*v2.PlaybackFinished)
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
				e := v.(*v2.PlaybackFinished)
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
