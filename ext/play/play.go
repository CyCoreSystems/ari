package play

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// AllDTMF is a string which contains all possible
// DTMF digits.
const AllDTMF = "0123456789ABCD*#"

// playStages executes a staged playback, waiting for its completion
func playStaged(ctx context.Context, h *ari.PlaybackHandle, opts *Options) (Status, error) {
	if opts == nil {
		opts = NewDefaultOptions()
	}

	started := h.Subscribe(ari.Events.PlaybackStarted)
	defer started.Cancel()
	finished := h.Subscribe(ari.Events.PlaybackFinished)
	defer finished.Cancel()

	err := h.Exec()
	if err != nil {
		return Failed, errors.Wrap(err, "failed to start playback")
	}
	defer h.Stop() // nolint: errcheck

	select {
	case <-ctx.Done():
		return Cancelled, nil
	case <-time.After(opts.playbackStartTimeout):
		return Timeout, errors.New("timeout waiting for playback to start")
	case <-finished.Events():
		return Finished, nil
	case <-started.Events():
	}

	// Wait for playback to complete
	select {
	case <-ctx.Done():
		return Cancelled, nil
	case <-finished.Events():
		return Finished, nil
	}
}

// NewPlay creates a new audio Options suitable for general audio playback
func NewPlay(ctx context.Context, p ari.Player, opts ...OptionFunc) (*Options, error) {
	o := NewDefaultOptions()
	err := o.ApplyOptions(opts...)

	return o, err
}

// Play plays the given media URI
func Play(ctx context.Context, p ari.Player, opts ...OptionFunc) *Result {
	o, err := NewPlay(ctx, p, opts...)
	if err != nil && o.result.Error != nil {
		o.result.Error = err
		return o.result
	}

	o.result.Error = o.Play(ctx, p)
	return o.result
}

// PlayAsync executes a playback, returning a channel which will be closed when the playback is ended
func (o *Options) PlayAsync(ctx context.Context, p ari.Player) <-chan error {
	ch := make(chan error)

	go func() {
		ch <- o.Play(ctx, p)
		close(ch)
	}()

	return ch
}

// Play executes a playback with the existing options.
func (o *Options) Play(ctx context.Context, p ari.Player) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if o.result == nil {
		o.result = new(Result)
	}

	if o.uriList == nil {
		return errors.New("empty playback URI list")
	}

	// cancel if we go over the maximum time
	go o.watchMaxTime(ctx)

	// Listen for DTMF
	go o.listenDTMF(ctx, p)

	for i := 0; i < o.maxReplays+1; i++ {
		if ctx.Err() != nil {
			break
		}

		// Reset the digit cache
		o.result.mu.Lock()
		o.result.DTMF = ""
		o.result.mu.Unlock()

		// Play the sequence of audio URIs
		err := o.playSequence(ctx, p)
		if err != nil {
			return err
		}

		// Wait for digits in the silence after the playback sequence completes
		o.waitDigits(ctx)
	}

	return nil

}

// playSequence plays the complete audio sequence
func (o *Options) playSequence(ctx context.Context, p ari.Player) (err error) {
	s := newSequence(o)
	defer s.Stop()

	go s.Play(ctx, p)

	select {
	case <-ctx.Done():
	case <-o.digitChan:
	case err = <-s.Done():
	}

	return err
}

func (o *Options) waitDigits(ctx context.Context) {
	overallTimer := time.NewTimer(o.overallDigitTimeout)
	defer overallTimer.Stop()

	digitTimeout := o.firstDigitTimeout

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(digitTimeout):
			return
		case <-overallTimer.C:
			return
		case <-o.digitChan:
			digitTimeout = o.interDigitTimeout

			// Determine if a match was found
			if o.matchFunc != nil {
				o.result.mu.Lock()
				o.result.DTMF, o.result.MatchResult = o.matchFunc(o.result.DTMF)
				o.result.mu.Unlock()

				switch o.result.MatchResult {
				case Complete:
					// If we have a complete response, close the entire playback
					// and return
					o.Stop()
					return
				case Invalid:
					// If invalid, return without waiting
					// for any more digits
					return
				default:
					// Incomplete means we should wait for more
				}
			}
		}
	}

}

// Stop terminates the execution of a playback
func (o *Options) Stop() {
	if o.result == nil {
		o.result = new(Result)
	}

	if o.result.Status == InProgress {
		o.result.Status = Cancelled
	}

	if o.cancel != nil {
		o.cancel()
	}
}

func (o *Options) watchMaxTime(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(o.maxPlaybackTime):
		o.Stop()
	}
}
func (o *Options) listenDTMF(ctx context.Context, p ari.Player) {
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
			o.result.mu.Lock()
			o.result.DTMF += v.Digit
			o.result.mu.Unlock()

			// Signal receipt of digit, but never block in doing so
			select {
			case o.digitChan <- v.Digit:
			default:
			}

		}
	}
}
