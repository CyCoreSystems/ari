package prompt

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext"
	"github.com/CyCoreSystems/ari/ext/audio"

	"golang.org/x/net/context"
)

// Options describes options for the prompt
type Options struct {
	// FirstDigitTimeout is the maximum length of time to wait
	// after the prompt sequence ends for the user to enter
	// a response.
	// If not specified, the default is DefaultFirstDigitTimeout.
	FirstDigitTimeout time.Duration

	// InterDigitTimeout is the maximum length of time to wait
	// for an additional digit after a digit is received.
	// If not specified, the default is DefaultInterDigitTimeout.
	InterDigitTimeout time.Duration

	// OverallTimeout is the maximum length of time to wait
	// for a response regardless of digits received.
	// If not specified, the default is DefaultOverallTimeout.
	OverallTimeout time.Duration

	// EchoData is the flag for saying each digit as it is input
	EchoData bool

	// MatchFunc is an optional function which, if supplied, returns
	// a string and an int.
	//
	// The string is allows the MatchFunc to return a different number
	// to be used as `result.Data`.  This is commonly used for prompts
	// which look for a terminator.  In such a practice, the terminator
	// would be stripped from the match and this argument would be populated
	// with the result.  Otherwise, the original string should be returned.
	// NOTE: Whatever is returned here will become `result.Data`.
	//
	// The int parameter indicates the result of the match, and it should
	// be one of:
	//  Incomplete (0) : insufficient digits to determine match.
	//  Complete (1) : A match was found.
	//  Invalid (2) : A match could not be found, given the digits received.
	// If this function returns a non-zero int, then the prompt will be stopped.
	// If not specified MatchAny will be used.
	MatchFunc func(string) (string, ext.Status)

	// Which type of word to use when playing '#'
	SoundHash string // pound or hash
}

var (
	// DefaultFirstDigitTimeout is the maximum time to wait for the
	// first digit after a prompt, if not otherwise set.
	DefaultFirstDigitTimeout = 4 * time.Second

	// DefaultInterDigitTimeout is the maximum time to wait for additional
	// digits after the first is received.
	DefaultInterDigitTimeout = 3 * time.Second

	// DefaultOverallTimeout is the maximum time to wait for a response
	// regardless of the number of received digits or pattern matching.
	DefaultOverallTimeout = 3 * time.Minute
)

type stateFn func() stateFn

// Prompt plays the given sound and waits for user input.
// nolint:gocyclo
func Prompt(ctx context.Context, p ari.Player, opts *Options, sounds ...string) (ret *ext.Result, err error) {
	ret = &ext.Result{}

	// Handle default options
	if opts == nil {
		opts = &Options{}
	}
	if opts.FirstDigitTimeout == 0 {
		opts.FirstDigitTimeout = DefaultFirstDigitTimeout
	}
	if opts.InterDigitTimeout == 0 {
		opts.InterDigitTimeout = DefaultInterDigitTimeout
	}
	if opts.OverallTimeout == 0 {
		opts.OverallTimeout = DefaultOverallTimeout
	}
	if opts.MatchFunc == nil {
		opts.MatchFunc = MatchAny
	}
	if opts.SoundHash == "" {
		opts.SoundHash = "hash"
	}

	dtmfSub := p.Subscribe(ari.Events.ChannelDtmfReceived)
	defer dtmfSub.Cancel()

	hangupSub := p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)
	defer hangupSub.Cancel()

	var playPrompt func() stateFn
	var waitDigit func() stateFn
	var waitFirstDigit func() stateFn
	overallTimer := time.NewTimer(opts.OverallTimeout)

	playPrompt = func() stateFn {
		playCtx, playCancel := context.WithCancel(context.Background())
		defer playCancel()

		doneCh := make(chan struct{})
		defer close(doneCh)

		lockCh := make(chan struct{})

		go func() {
			defer close(lockCh)
			defer playCancel()

			for {
				select {
				case _, ok := <-hangupSub.Events():
					if ok {
						ret.Status = ext.Hangup
						return
					}
				case <-ctx.Done():
					ret.Status = ext.Canceled
					return
				case <-doneCh:
					return
				case e := <-dtmfSub.Events():
					ret.Data += e.(*ari.ChannelDtmfReceived).Digit
					Logger.Debug("DTMF received", "digits", ret.Data)
					match, res := opts.MatchFunc(ret.Data)
					ret.Data = match
					if res > 0 {
						ret.Status = res
						playCancel() // cancel playback
						return
					}
				}
			}
		}()

		var st ext.Status
		st, err = audio.Queue(playCtx, p, sounds...)
		if st == ext.Canceled {
			err = nil
		} else if st != ext.Complete {
			ret.Status = st
			return nil
		}

		// helps us destroy the above goroutine and ensure its closure before reading ret.Data
		doneCh <- struct{}{}
		<-lockCh

		if !overallTimer.Stop() {
			select {
			case <-overallTimer.C:
			default:
			}
		}
		overallTimer.Reset(opts.OverallTimeout)

		if ret.Data == "" {
			return waitFirstDigit
		}

		return waitDigit
	}

	waitFirstDigit = func() stateFn {
		select {
		case _, ok := <-hangupSub.Events():
			if ok {
				ret.Status = ext.Hangup
				return nil
			}
		case <-ctx.Done():
			ret.Status = ext.Canceled
			return nil
		case <-time.After(opts.FirstDigitTimeout):
			ret.Status = ext.Timeout
			return nil
		case <-overallTimer.C:
			ret.Status = ext.Timeout
			return nil
		case e := <-dtmfSub.Events():
			ret.Data += e.(*ari.ChannelDtmfReceived).Digit
			Logger.Debug("DTMF received", "digits", ret.Data)
			match, res := opts.MatchFunc(ret.Data)
			ret.Data = match
			if res > 0 {
				ret.Status = res
				return nil
			}
		}

		return waitDigit
	}

	waitDigit = func() stateFn {
		select {
		case _, ok := <-hangupSub.Events():
			if ok {
				ret.Status = ext.Hangup
				return nil
			}
		case <-ctx.Done():
			ret.Status = ext.Canceled
			return nil
		case <-time.After(opts.InterDigitTimeout):
			ret.Status = ext.Timeout
			return nil
		case <-overallTimer.C:
			ret.Status = ext.Timeout
			return nil
		case e := <-dtmfSub.Events():
			ret.Data += e.(*ari.ChannelDtmfReceived).Digit
			Logger.Debug("DTMF received", "digits", ret.Data)
			match, res := opts.MatchFunc(ret.Data)
			ret.Data = match
			if res > 0 {
				ret.Status = res
				return nil
			}
		}

		return waitDigit
	}

	var st = playPrompt
	if len(sounds) == 0 {
		st = waitFirstDigit
	}

	for st != nil {
		st = st()
	}

	return
}
