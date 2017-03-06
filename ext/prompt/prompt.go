package prompt

import (
	"time"

	"github.com/CyCoreSystems/ari/ext/audio"

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

	// LastTimeout is set to the timeout being used in the next stateFn to be called.
	// it is then being reference in Prompt in the st, err = st(.....) call
	LastTimeout = 0 * time.Second
)

type stateFn func(ctx context.Context) (stateFn, error)

type promptApp struct {
	dSub     ari.Subscription
	hSub     ari.Subscription
	oTimer   *time.Timer
	LTimeout time.Duration
	options  *Options
	player   audio.Player
	snds     []string
	retData  *Result
}

// Prompt plays the given sound and waits for user input.
func Prompt(ctx context.Context, p audio.Player, opts *Options, sounds ...string) (ret *Result, err error) {

	//snds = sounds

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

	overallTimer := time.NewTimer(opts.OverallTimeout)

	LastTimeout = 0 * time.Second
	/*if len(sounds) == 0 {
		LastTimeout = opts.FirstDigitTimeout
		st, err = waitDigit(ctx, opts.FirstDigitTimeout, opts, ret)
	}
	*/
	pa := &promptApp{
		dSub:     p.Subscribe(ari.Events.ChannelDtmfReceived),
		hSub:     p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed),
		oTimer:   overallTimer,
		LTimeout: LastTimeout,
		options:  opts,
		player:   p,
		snds:     sounds,
		retData:  &Result{},
	}
	st, err := pa.playPrompt(ctx)
	for st != nil && err == nil {
		st, err = st(ctx)
		select {
		case <-ctx.Done():
			Logger.Debug("Exiting on context closure")
			st = nil
		default:
		}
	}

	return pa.retData, err
}

func (pa *promptApp) playPrompt(ctx context.Context) (stateFn, error) {
	playCtx, playCancel := context.WithCancel(context.Background())
	opts := pa.options

	defer playCancel()

	doneCh := make(chan struct{})
	defer close(doneCh)

	lockCh := make(chan struct{})

	go func() {
		defer close(lockCh)
		defer playCancel()

		for {
			select {
			case _, ok := <-pa.hSub.Events():
				if ok {
					pa.retData.Status = Hangup
					return
				}
			case <-ctx.Done():
				pa.retData.Status = Canceled
				return
			case <-doneCh:
				return
			case e := <-pa.dSub.Events():
				pa.retData.Data += e.(*ari.ChannelDtmfReceived).Digit
				Logger.Debug("DTMF received", "digits", pa.retData.Data)
				match, res := opts.MatchFunc(pa.retData.Data)
				pa.retData.Data = match
				if res > 0 {
					pa.retData.Status = res
					playCancel() // cancel playback
					return
				}
			}
		}
	}()

	q := audio.NewQueue()
	q.Add(pa.snds...)
	//var st audio.Status
	st, err := q.Play(playCtx, pa.player, &audio.Options{ExitOnDTMF: ""})

	switch st {
	case audio.Canceled:
		//NOTE: since playCtx doesn't extend the parent context,
		// any audio cancel is considered a special case
		// and we don't overwrite the return status
		err = nil
	case audio.Hangup:
		pa.retData.Status = Hangup
	case audio.Failed:
		pa.retData.Status = Failed
	case audio.Timeout:
		pa.retData.Status = Failed
	}

	if pa.retData.Status > Incomplete {
		return nil, err
	}

	// helps us destroy the above goroutine and ensure its closure before reading ret.Data
	doneCh <- struct{}{}
	<-lockCh

	if pa.oTimer.Stop() {
		select {
		case <-pa.oTimer.C:
		default:
		}
	}
	pa.oTimer.Reset(opts.OverallTimeout)

	if pa.retData.Data == "" {
		pa.LTimeout = opts.FirstDigitTimeout
	} else {
		pa.LTimeout = opts.InterDigitTimeout
	}

	return pa.waitDigit, err
}

func (pa *promptApp) waitDigit(ctx context.Context) (stateFn, error) {
	opts := pa.options
	select {
	case _, ok := <-pa.hSub.Events():
		if ok {
			pa.retData.Status = Hangup
			return nil, nil
		}
	case <-ctx.Done():
		pa.retData.Status = Canceled
		return nil, nil
	case <-time.After(pa.LTimeout):
		pa.retData.Status = Timeout
		return nil, nil
	case <-pa.oTimer.C:
		pa.retData.Status = Timeout
		return nil, nil
	case e := <-pa.dSub.Events():
		pa.retData.Data += e.(*ari.ChannelDtmfReceived).Digit
		Logger.Debug("DTMF received", "digits", pa.retData.Data)
		match, res := opts.MatchFunc(pa.retData.Data)
		pa.retData.Data = match
		if res > 0 {
			pa.retData.Status = res
			return nil, nil
		}
	}

	LastTimeout = opts.InterDigitTimeout
	return pa.waitDigit, nil
}
