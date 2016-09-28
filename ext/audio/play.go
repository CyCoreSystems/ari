package audio

import (
	"errors"
	"time"

	"github.com/CyCoreSystems/ari"
	uuid "github.com/satori/go.uuid"
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

// Play plays the given media URI
func Play(ctx context.Context, p Player, mediaURI string) (st Status, err error) {
	pb := PlayAsync(ctx, p, mediaURI)

	<-pb.Stopped()

	st, err = pb.Status(), pb.Err()
	return
}

// PlayAsync plays the audio asynchronously and returns a playback object
func PlayAsync(ctx context.Context, p Player, mediaURI string) *Playback {

	var pb Playback

	pb.startCh = make(chan struct{})
	pb.stopCh = make(chan struct{})
	pb.status = InProgress
	pb.err = nil
	pb.ctx, pb.cancel = context.WithCancel(ctx)

	//TODO: confirm whether we need to listen on bridge events if p Player is a bridge
	hangup := p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)

	id := uuid.NewV1().String()
	pb.handle, pb.err = p.Play(id, mediaURI)

	// register for events on the ~~playback~~ player handle. This means
	// we have to filter the events using the evnentual playback handle.
	playbackStarted := pb.handle.Subscribe(ari.Events.PlaybackStarted)
	playbackFinished := pb.handle.Subscribe(ari.Events.PlaybackFinished)

	go func() {
		defer func() {
			playbackStarted.Cancel()
			playbackFinished.Cancel()
			hangup.Cancel()
			close(pb.stopCh)
		}()

		// wait to check error here so
		// subscriptions are cleaned up
		if pb.err != nil {
			close(pb.startCh)
			pb.status = Failed
			return
		}

		go func() {
			defer close(pb.startCh)

			for {
				select {
				case <-time.After(PlaybackStartTimeout):
					pb.status = Timeout
					pb.err = errors.New("Timeout waiting for start of playback")
					return
				case <-hangup.Events():
					pb.status = Hangup
					return
				case <-pb.ctx.Done():
					pb.status = Canceled
					pb.err = pb.ctx.Err()
					return
				case <-playbackFinished.Events():
					Logger.Debug("Got playback finished before start")
					pb.status = Finished
					return
				case <-playbackStarted.Events():
					return
				}
			}
		}()

		<-pb.startCh

		if pb.status != InProgress {
			return
		}

		for {
			select {
			case <-time.After(MaxPlaybackTime):
				pb.status = Timeout
				pb.err = errors.New("Timeout waiting for stop of playback")
				return
			case <-hangup.Events():
				pb.status = Hangup
				return
			case <-pb.ctx.Done():
				pb.status = Canceled
				pb.err = pb.ctx.Err()
				return
			case <-playbackFinished.Events():
				pb.status = Finished
				return
			}
		}
	}()

	return &pb
}
