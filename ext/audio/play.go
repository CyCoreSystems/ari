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
func Play(ctx context.Context, playback ari.Playback, p Player, mediaURI string) (st Status, err error) {
	pb := PlayAsync(ctx, playback, p, mediaURI)

	select {
	case <-pb.Stopped():
	}

	st, err = pb.Status(), pb.Err()
	return
}

// PlayAsync plays the audio asynchronously and returns a playback object
func PlayAsync(ctx context.Context, playback ari.Playback, p Player, mediaURI string) *Playback {

	var pb Playback

	id := uuid.NewV1().String()
	handle := playback.Get(id)

	pb.handle = handle
	pb.startCh = make(chan struct{})
	pb.stopCh = make(chan struct{})
	pb.status = InProgress
	pb.err = nil
	pb.ctx, pb.cancel = context.WithCancel(ctx)

	// register for events on the playback handle
	playbackStarted := handle.Subscribe(ari.Events.PlaybackStarted)
	playbackFinished := handle.Subscribe(ari.Events.PlaybackFinished)

	//TODO: confirm whether we need to listen on bridge events if p Player is a bridge
	hangup := p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)

	go func() {
		defer func() {
			playbackStarted.Cancel()
			playbackFinished.Cancel()
			hangup.Cancel()
			close(pb.stopCh)
		}()

		pb.handle, pb.err = p.Play(id, mediaURI)
		if pb.err != nil {
			close(pb.startCh)
			pb.status = Failed
			return
		}

		select {
		case <-time.After(PlaybackStartTimeout):
			pb.status = Timeout
			pb.err = errors.New("Timeout waiting for start of playback")
			close(pb.startCh)
			return
		case <-hangup.Events():
			pb.status = Hangup
			close(pb.startCh)
			return
		case <-pb.ctx.Done():
			pb.status = Canceled
			pb.err = pb.ctx.Err()
			close(pb.startCh)
			return
		case <-playbackStarted.Events():
			close(pb.startCh)
		}

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
		}

		pb.status = Finished
		return
	}()

	return &pb
}
