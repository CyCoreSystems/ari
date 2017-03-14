package audio

import (
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
)

// Options describes various options which
// are available to playback operations.
type Options struct {
	// ID is an optional ID to use for the playback's ID. If one
	// is not supplied, an ID will be randomly generated internally.
	// NOTE that this ID will only be used for the FIRST playback
	// in a queue.  All subsequent playback IDs will be randomly generated.
	ID string

	// ExitOnDTMF defines a list of DTMF digits on receipt of which will
	// terminate the playback of the queue.  You may set this to AllDTMF
	// in order to match any DTMF digit.
	ExitOnDTMF string

	// Done is an optional channel for receiving notification when the playback
	// is complete.  This is useful if the playback is to be executed asynchronously.
	// This channel will be closed by the playback when playback the is complete.
	Done chan<- struct{}
}

// Queue represents a sequence of audio playbacks
// which are to be played on the associated Player
type Queue struct {
	mu    sync.Mutex
	queue []string // List of mediaURI to be played

	receivedDTMF string // Storage for received DTMF, if we are listening for them

	cancel context.CancelFunc
}

// NewQueue creates (but does not start) a new playback queue.
func NewQueue() *Queue {
	return &Queue{}
}

// Add appends one or more mediaURIs to the playback queue
func (q *Queue) Add(mediaURIs ...string) {
	q.mu.Lock()

	// Add each media URI to the queue
	for _, u := range mediaURIs {
		if u == "" {
			continue
		}

		q.queue = append(q.queue, u)
	}
	q.mu.Unlock()
}

// Flush empties a playback queue.
// NOTE that this does NOT stop the current playback.
func (q *Queue) Flush() {
	q.queue = []string{}
}

// ReceivedDTMF returns any DTMF which has been received
// by the PlaybackQueue.
func (q *Queue) ReceivedDTMF() string {
	return q.receivedDTMF
}

// Play starts the playback of the queue to the Player.
func (q *Queue) Play(ctx context.Context, p ari.Player, opts *Options) (Status, error) {
	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel
	defer cancel()

	if opts == nil {
		opts = &Options{}
	}

	if opts.Done != nil {
		defer close(opts.Done)
	}

	if opts.ExitOnDTMF != "" {
		go q.monitorDTMF(ctx, p, opts.ExitOnDTMF)
	}

	// Start the playback
	for i := 0; i < len(q.queue); i++ {
		if ctx.Err() != nil {
			return Canceled, nil
		}
		status, err := Play(ctx, p, q.queue[i])
		if err != nil {
			return status, err
		}
	}

	return Finished, nil
}

func (q *Queue) monitorDTMF(ctx context.Context, p ari.Player, exitList string) {
	defer q.cancel()

	sub := p.Subscribe(ari.Events.ChannelDtmfReceived)
	defer sub.Cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case e := <-sub.Events():
			v, ok := e.(*ari.ChannelDtmfReceived)
			if !ok {
				return
			}
			q.receivedDTMF += v.Digit
			if strings.Contains(exitList, v.Digit) {
				return
			}
		}
	}
}
