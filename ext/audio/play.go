package audio

import (
	"time"

	"sync"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext"
	"github.com/pkg/errors"
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
func Play(ctx context.Context, p ari.Player, mediaURI string) (ext.Status, error) {
	pb := p.StagePlay(uuid.NewV1().String(), mediaURI)
	return watchEvents(ctx, p, pb, nil)
}

func watchEvents(ctx context.Context, p ari.Player, pb ari.PlaybackHandle, wg *sync.WaitGroup) (st ext.Status, err error) {

	startedSub := pb.Subscribe(ari.Events.PlaybackStarted)
	defer startedSub.Cancel()
	finishedSub := pb.Subscribe(ari.Events.PlaybackFinished)
	defer finishedSub.Cancel()

	//TODO: confirm whether we need to listen on bridge events if p Player is a bridge
	hangupSub := p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)
	defer hangupSub.Cancel()

	if wg != nil {
		wg.Done()
	}

	if err = pb.Exec(); err != nil {
		return
	}

	select {
	case <-ctx.Done():
		st = ext.Canceled
		err = ctx.Err()
	case <-time.After(PlaybackStartTimeout):
		st = ext.Timeout
		err = errors.New("Timeout waiting for start of playback")
	case <-hangupSub.Events():
		st = ext.Hangup
	case <-finishedSub.Events():
		Logger.Warn("Got playback finished before start")
		st = ext.Complete
	case <-startedSub.Events():
		st = ext.Incomplete
	}

	if err == nil || st != ext.Incomplete {
		return
	}

	select {
	case <-ctx.Done():
		st = ext.Canceled
		err = ctx.Err()
	case <-time.After(MaxPlaybackTime):
		st = ext.Timeout
		err = errors.New("Timeout waiting for stop of playback")
	case <-hangupSub.Events():
		st = ext.Hangup
	case <-finishedSub.Events():
		st = ext.Complete
	}

	return
}
