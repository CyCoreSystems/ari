package audio

import (
	"strings"
	"sync"

	"github.com/CyCoreSystems/ari"

	"golang.org/x/net/context"
)

// Options describes various options which
// are available to playback operations.
type Options struct {
	// ID is an optional ID to use for the playback's ID. If one
	// is not supplied, an ID will be randomly generated internally.
	// NOTE that this ID will only be used for the FIRST playback
	// in a queue.  All subsequent playback IDs will be randomly generated.
	ID string

	// DTMF is an optional channel for received DTMF tones received during the playback.
	// This channel will NOT be closed by the playback.
	DTMF chan<- *ari.ChannelDtmfReceived

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
	s            ari.Subscriber
	receivedDTMF string // Storage for received DTMF, if we are listening for them
}

// NewQueue creates (but does not start) a new playback queue.
func NewQueue(s ari.Subscriber) *Queue {
	return &Queue{
		s: s,
	}
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
func (pq *Queue) Play(ctx context.Context, p Player, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}

	ctrlCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// NOTE: this code used to call Subscribe("ChannelDtmfReceived") twice. This
	// was /fine/ until trying to unit test. We use one subscription
	// and send out to dtmfChan depending on whether opts.DTMF is not nil.
	// This simplifies the workflow and ensures we have a 1-1 correlation between
	// a subscription and a Queue

	dtmfChan := make(chan *ari.ChannelDtmfReceived)
	dtmfSub := pq.s.Subscribe(ari.Events.ChannelDtmfReceived)

	handlingDTMF := false

	// Handle any options we were given
	if opts != nil {
		// Close the done channel when we finish,
		// if we were given one.
		if opts.Done != nil {
			defer close(opts.Done)
		}

		// Listen for DTMF, if we were asked to do so
		if opts.DTMF != nil {
			handlingDTMF = true
			go func() {
				for {
					select {
					case <-ctrlCtx.Done():
						return
					case <-ctx.Done():
						return
					case e := <-dtmfChan:
						opts.DTMF <- e
					}
				}
			}()
		}
	}

	// if we aren't forwarding to opts.DTMF, then discard DTMF messages
	if !handlingDTMF {
		go func() {
			for {
				select {
				case <-ctrlCtx.Done():
					return
				case <-ctx.Done():
					return
				case <-dtmfChan:
				}
			}
		}()
	}

	// Record any DTMF (this is separate from opts.DTMF) so that we can
	//  - Service ReceivedDTMF requests
	//  - Exit if we were given an ExitOnDTMF list
	go func() {
		defer dtmfSub.Cancel()
		for {
			select {
			case <-ctrlCtx.Done():
				return
			case <-ctx.Done():
				return
			case e := <-dtmfSub.Events():
				if e == nil {
					return
				}
				dtmfChan <- e.(*ari.ChannelDtmfReceived)
				digit := e.(*ari.ChannelDtmfReceived).Digit
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
		case <-ctrlCtx.Done():
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get the next clip
		err := Play(ctx, pq.s, p, pq.queue[i])
		if err != nil {
			return err
		}
	}

	return nil
}
