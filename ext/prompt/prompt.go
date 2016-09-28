package prompt

import (
	"time"

	"github.com/CyCoreSystems/ari/ext/audio"
	"github.com/CyCoreSystems/ari/ext/audio/audiouri"

	"github.com/CyCoreSystems/ari"

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
	MatchFunc func(string) (string, Status)

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

// Prompt plays the given sound and waits for user input.
func Prompt(ctx context.Context, p audio.Player, opts *Options, sounds ...string) (ret *Result, err error) {
	ret = &Result{}

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

	hangupSub := p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)
	defer hangupSub.Cancel()

	dtmfSub := p.Subscribe(ari.Events.ChannelDtmfReceived)
	defer dtmfSub.Cancel()

	playCtx, playCancel := context.WithCancel(context.Background())
	defer playCancel()

	go func() {

		select {
		case <-hangupSub.Events():
			ret.Status = Hangup
			playCancel()
		case <-ctx.Done():
			ret.Status = Canceled
			playCancel()
		}
	}()

	// Play the prompt, if we have one
	if sounds != nil && len(sounds) != 0 {
		q := audio.NewQueue()
		q.Add(sounds...)
		var st audio.Status
		st, err = q.Play(playCtx, p, &audio.Options{
			ExitOnDTMF: audio.AllDTMF,
		})

		switch st {
		case audio.Canceled:
			//NOTE: since playCtx doesn't extend the parent context,
			// any audio cancel is considered a special case
			// and we don't overwrite the return status
			err = nil
			return
		case audio.Hangup:
			ret.Status = Hangup
			return
		case audio.Failed:
			ret.Status = Failed
			return
		case audio.Timeout:
			ret.Status = Failed
			return
		}
	}

	// Construct timeouts
	timeoutChan := make(chan time.Time)
	var interDigitTimeout = func(digits string) *time.Timer {
		idTimer := time.AfterFunc(opts.InterDigitTimeout, func() {
			if len(digits) == len(ret.Data) {
				Logger.Debug("Inter-digit timeout elapsed")
				timeoutChan <- time.Now()
			}
		})
		return idTimer
	}
	otimer := time.AfterFunc(opts.OverallTimeout, func() {
		// Overall timeout
		Logger.Debug("Overall timeout elapsed")
		timeoutChan <- time.Now()
	})
	defer otimer.Stop()

	Logger.Debug("current data", "data", ret.Data)

	var timer *time.Timer

	// Apply first-digit timeout or inter-digit timeout, as appropriate
	if len(ret.Data) == 0 {
		timer = time.AfterFunc(opts.FirstDigitTimeout, func() {
			// First digit timeout
			if ret.Data == "" {
				Logger.Debug("First-digit timeout elapsed")
				timeoutChan <- time.Now()
			}
		})
		defer timer.Stop()
	} else {
		// We already have digits, so the inter-digit timeout applies
		timer = interDigitTimeout(ret.Data)
		defer timer.Stop()
	}

	// Wait for response
	for {
		select {
		case <-ctx.Done():
			Logger.Debug("Prompt wait canceled")
			ret.Status = Canceled
			err = ctx.Err()
			return
		case <-hangupSub.Events():
			Logger.Debug("Hangup after prompt playback")
			ret.Status = Hangup
			return
		case e := <-dtmfSub.Events():
			ret.Data += e.(*ari.ChannelDtmfReceived).Digit
			Logger.Debug("DTMF received", "digits", ret.Data)
			match, res := opts.MatchFunc(ret.Data)
			ret.Data = match
			if res > 0 {
				ret.Status = res
				return
			}

			// If we are set to echo the digits back, do so
			if opts.EchoData {
				go func(currentData string) {
					Logger.Debug("Echoing digits", "data", currentData)

					digitsQueue := audio.NewQueue()
					digitsQueue.Add(audiouri.DigitsURI(ret.Data, opts.SoundHash)...)

					if _, err := digitsQueue.Play(ctx, p, nil); err != nil {
						Logger.Error("Error saying digits", "error", err)
					}

					// Set inter-digit timeout after completion of playback
					interDigitTimeout(currentData)
				}(ret.Data)
			} else {
				// Set inter-digit timeout after completion of playback
				interDigitTimeout(ret.Data)
			}

			// Set inter-digit timeout
			interDigitTimeout(ret.Data)

		case <-timeoutChan:
			Logger.Debug("Timeout waiting for prompt response")
			ret.Status = Timeout
			return
		}
	}
}
