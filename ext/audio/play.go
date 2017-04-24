package audio

import (
	"sync"
	"time"

	"github.com/CyCoreSystems/ari"
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
func Play(ctx context.Context, p ari.Player, mediaURI string) (Status, error) {
	c := PlayAsync(ctx, p, mediaURI)

	<-c.Stopped()

	return c.Status(), c.Err()
}

// PlayAsync plays the audio asynchronously and returns a playback object
func PlayAsync(ctx context.Context, p ari.Player, mediaURI string) *Control {
	var err error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := Control{
		id:      uuid.NewV1().String(),
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
		status:  InProgress,
		cancel:  cancel,
	}

	c.pb, err = p.StagePlay(c.id, mediaURI)
	if err != nil {
		c.err = errors.Wrap(err, "failed to create playback")
		c.status = Failed
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go c.watchEvents(ctx, wg)
	wg.Wait()

	err = c.pb.Exec()
	if err != nil {
		c.err = errors.Wrap(err, "failed to start playback")
		c.status = Failed
	}

	return &c
}

// Control provides a mechanism for interacting with an audio playback
type Control struct {
	id string // playback ID

	started bool
	startCh chan struct{}

	stopped bool
	stopCh  chan struct{}

	p  ari.Player
	pb *ari.PlaybackHandle

	status Status
	err    error

	startedSub  ari.Subscription
	finishedSub ari.Subscription
	hangupSub   ari.Subscription

	cancel context.CancelFunc
}

// Handle returns the ARI reference to the playback object
func (c *Control) Handle() *ari.PlaybackHandle {
	return c.pb
}

func (c *Control) onStarted() {
	if !c.started {
		c.started = true
		close(c.startCh)
	}
}

// Started returns the channel that is closed when the playback has started
func (c *Control) Started() <-chan struct{} {
	return c.startCh
}

func (c *Control) onStopped() {
	if !c.stopped {
		c.stopped = true
		close(c.stopCh)
	}
}

// Stopped returns the channel that is closed when the playback has stopped
func (c *Control) Stopped() <-chan struct{} {
	return c.stopCh
}

// Status returns the current status of the playback
func (c *Control) Status() Status {
	return c.status
}

// Err returns any accumulated errors during playback
func (c *Control) Err() error {
	return c.err
}

// Cancel stops the playback
func (c *Control) Cancel() {
	if !c.stopped && c.pb != nil {
		c.pb.Stop()
	}

	c.onStopped()

	c.cancel()
}

type stateFn func(context.Context) stateFn

func (c *Control) watchEvents(ctx context.Context, wg *sync.WaitGroup) {
	defer c.Cancel()

	c.startedSub = c.pb.Subscribe(ari.Events.PlaybackStarted)
	defer c.startedSub.Cancel()
	c.finishedSub = c.pb.Subscribe(ari.Events.PlaybackFinished)
	defer c.finishedSub.Cancel()

	//TODO: confirm whether we need to listen on bridge events if p Player is a bridge
	c.hangupSub = c.pb.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)
	defer c.hangupSub.Cancel()

	wg.Done()

	for f := c.waitStart; f != nil; {
		f = f(ctx)
	}
}

func (c *Control) waitStart(ctx context.Context) stateFn {
	select {
	case <-ctx.Done():
		c.status = Canceled
		c.err = ctx.Err()
		return nil
	case <-time.After(PlaybackStartTimeout):
		c.status = Timeout
		c.err = errors.New("Timeout waiting for start of playback")
		return nil
	case <-c.hangupSub.Events():
		c.status = Hangup
		return nil
	case <-c.finishedSub.Events():
		Logger.Warn("Got playback finished before start")
		c.status = Finished
		c.onStopped()
		return nil
	case <-c.startedSub.Events():
		c.onStarted()
		return c.waitStop
	}
}

func (c *Control) waitStop(ctx context.Context) stateFn {
	select {
	case <-ctx.Done():
		c.status = Canceled
		c.err = ctx.Err()
	case <-time.After(MaxPlaybackTime):
		c.status = Timeout
		c.err = errors.New("Timeout waiting for stop of playback")
	case <-c.hangupSub.Events():
		c.status = Hangup
	case <-c.finishedSub.Events():
		c.status = Finished
	}
	return nil
}
