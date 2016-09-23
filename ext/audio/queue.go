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
	queue        []string // List of mediaURI to be played
	mu           sync.Mutex
	receivedDTMF string // Storage for received DTMF, if we are listening for them
}

// NewQueue creates (but does not start) a new playback queue.
func NewQueue() *Queue {
	return &Queue{}
}

// Add appends one or more mediaURIs to the playback queue
func (pq *Queue) Add(mediaURIs ...string) {
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
func (pq *Queue) Flush() {
	pq.mu.Lock()
	pq.queue = []string{}
	pq.mu.Unlock()
}

// ReceivedDTMF returns any DTMF which has been received
// by the PlaybackQueue.
func (pq *Queue) ReceivedDTMF() string {
	return pq.receivedDTMF
}

// Play starts the playback of the queue to the Player.
func (pq *Queue) Play(ctx context.Context, playback ari.Playback, p Player, opts *Options) (Status, error) {

	if opts == nil {
		opts = &Options{}
	}

	if opts.Done != nil {
		defer close(opts.Done)
	}

	queue := make(chan string)

	dtmfSub := p.Subscribe(ari.Events.ChannelDtmfReceived)
	defer dtmfSub.Cancel()

	pq.queue = append(pq.queue, "")

	dtmfExit := make(chan struct{})

	go func() {
		defer close(queue)
		for i := 0; i != len(pq.queue); {
			select {
			case e := <-dtmfSub.Events(): // read dtmf input
				d := e.(*ari.ChannelDtmfReceived)
				pq.receivedDTMF += d.Digit

				if strings.Contains(opts.ExitOnDTMF, d.Digit) {
					close(dtmfExit)
					return
				}

			case queue <- pq.queue[i]: // send next item
				i++
				continue

			case <-ctx.Done(): // wait for cancellation
				return

			}
		}
	}()

	// Start the playback
	for q := range queue {
		if q == "" {
			break
		}
		pb := PlayAsync(ctx, playback, p, q)

		select {
		case <-dtmfExit:
			pb.Cancel()
			return Finished, nil
		case <-pb.Stopped():
			if pb.Status() > Finished {
				return pb.Status(), pb.Err()
			}
		}
	}

	return Finished, nil
}
