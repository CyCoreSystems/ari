package play

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/CyCoreSystems/ari/v5"
)

// Session describes a structured Play session.
type Session interface {
	// Add appends a set of AudioURIs to a Play session.  Note that if the Play session
	// has already been completed, this will NOT make it start again.
	Add(list ...string)

	// Done returns a channel which is closed when the Play session completes
	Done() <-chan struct{}

	// Err waits for a session to end and returns its error
	Err() error

	// StopAudio stops the playback of the audio sequence (if there is one), but
	// unlike `Stop()`, this does _not_ necessarily terminate the session.  If
	// the Play session is configured to wait for DTMF following the playback,
	// it will still wait after StopAudio() is called.
	StopAudio()

	// Result waits for a session to end and returns its result
	Result() (*Result, error)

	// Stop stops a Play session immediately
	Stop()
}

type playSession struct {
	o *Options

	// cancel is the playback context's cancel function
	cancel context.CancelFunc

	// currentSequence is a pointer to the currently-playing sequence, if there is one
	currentSequence *sequence

	// digitChan is the channel on which any received DTMF digits will be sent.  The received DTMF will also be stored separately, so this channel is primarily for signaling purposes.
	digitChan chan string

	// closed is a wrapper for done which indicates that done has been closed
	closed bool

	// done is a channel which is closed when the playback completes execution
	done chan struct{}

	// mu provides locking for concurrency-related datastructures within the options
	mu sync.Mutex

	// result is the final result of the playback
	result *Result
}

type nilSession struct {
	res *Result
}

func (n *nilSession) Add(list ...string) {
}

func (n *nilSession) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)

	return ch
}

func (n *nilSession) Err() error {
	return n.res.Error
}

func (n *nilSession) StopAudio() {
}

func (n *nilSession) Result() (*Result, error) {
	return n.res, n.res.Error
}

func (n *nilSession) Stop() {
}

func errorSession(err error) *nilSession {
	s := &nilSession{
		res: new(Result),
	}

	s.res.Error = err
	s.res.Status = Failed

	return s
}

func newPlaySession(o *Options) *playSession {
	return &playSession{
		o:         o,
		result:    new(Result),
		digitChan: make(chan string, DigitBufferSize),
		done:      make(chan struct{}),
	}
}

func (s *playSession) play(ctx context.Context, p ari.Player) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	defer s.Stop()

	if s.result == nil {
		s.result = new(Result)
	}

	if s.o.uriList == nil || s.o.uriList.Empty() {
		s.result.Error = errors.New("empty playback URI list")
		return
	}

	// cancel if we go over the maximum time
	go s.watchMaxTime(ctx)

	// Listen for DTMF
	go s.listenDTMF(ctx, p)

	for i := 0; i < s.o.maxReplays+1; i++ {
		if ctx.Err() != nil {
			break
		}

		// Reset the digit cache
		s.result.mu.Lock()
		s.result.DTMF = ""
		s.result.mu.Unlock()

		// Play the sequence of audio URIs
		s.playSequence(ctx, p, i)

		if s.result.Error != nil {
			return
		}

		// Wait for digits in the silence after the playback sequence completes
		s.waitDigits(ctx)
	}
}

// playSequence plays the complete audio sequence
func (s *playSession) playSequence(ctx context.Context, p ari.Player, playbackCounter int) {
	seq := newSequence(s)

	s.mu.Lock()
	s.currentSequence = seq
	s.mu.Unlock()

	go seq.Play(ctx, p, playbackCounter)

	// Wait for sequence playback to complete (or context closure to be caught)
	select {
	case <-ctx.Done():
	case <-seq.Done():
		if s.result.Status == InProgress {
			s.result.Status = Finished
		}
	}

	// Stop audio playback if it is still running
	seq.Stop()

	// wait for cleanup of sequence so we can get the proper error result
	<-seq.Done()
}

// nolint: gocyclo
func (s *playSession) waitDigits(ctx context.Context) {
	overallTimer := time.NewTimer(s.o.overallDigitTimeout)
	defer overallTimer.Stop()

	digitTimeout := s.o.firstDigitTimeout

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(digitTimeout):
			return
		case <-overallTimer.C:
			return
		case <-s.digitChan:
			if len(s.result.DTMF) > 0 {
				digitTimeout = s.o.interDigitTimeout
			}

			// Determine if a match was found
			if s.o.matchFunc != nil {
				s.result.mu.Lock()
				s.result.DTMF, s.result.MatchResult = s.o.matchFunc(s.result.DTMF)
				s.result.mu.Unlock()

				switch s.result.MatchResult {
				case Complete:
					// If we have a complete response, close the entire playback
					// and return
					s.Stop()
					return
				case Invalid:
					// If invalid, return without waiting
					// for any more digits
					return
				default: // Incomplete means we should wait for more
				}
			}
		}
	}
}

// Stop terminates the execution of a playback
func (s *playSession) Stop() {
	if s.result == nil {
		s.result = new(Result)
	}

	// Stop any audio which is still playing
	if s.currentSequence != nil {
		s.currentSequence.Stop()
		<-s.currentSequence.Done()
	}

	// If we have no other status set, set it to Cancelled
	if s.result.Status == InProgress {
		s.result.Status = Cancelled
	}

	// Close out anything else
	if s.cancel != nil {
		s.cancel()
	}

	s.mu.Lock()
	
	if !s.closed {
		s.closed = true
		close(s.done)
	}

	s.mu.Unlock()
}

func (s *playSession) watchMaxTime(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(s.o.maxPlaybackTime):
		s.Stop()
	}
}

func (s *playSession) listenDTMF(ctx context.Context, p ari.Player) {
	sub := p.Subscribe(ari.Events.ChannelDtmfReceived)
	defer sub.Cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case e := <-sub.Events():
			if e == nil {
				return
			}

			v, ok := e.(*ari.ChannelDtmfReceived)
			if !ok {
				continue
			}

			s.result.mu.Lock()
			s.result.DTMF += v.Digit
			s.result.mu.Unlock()

			// Signal receipt of digit, but never block in doing so
			select {
			case s.digitChan <- v.Digit:
			default:
			}

			// If we have a MatchFunc, stop any playing audio
			if s.o.matchFunc != nil && s.currentSequence != nil {
				s.currentSequence.Stop()
			}
		}
	}
}

func (s *playSession) Add(list ...string) {
	for _, i := range list {
		s.o.uriList.Add(i)
	}
}

func (s *playSession) Done() <-chan struct{} {
	return s.done
}

func (s *playSession) Err() error {
	<-s.Done()
	return s.result.Error
}

func (s *playSession) StopAudio() {
	if s.currentSequence != nil {
		s.currentSequence.Stop()
		<-s.currentSequence.Done()
	}
}

func (s *playSession) Result() (*Result, error) {
	<-s.Done()
	return s.result, s.result.Error
}
