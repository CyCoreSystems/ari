package prompt

import (
	"time"

	"github.com/CyCoreSystems/ari/ext/audio"

	"github.com/CyCoreSystems/ari"
	//"github.com/CyCoreSystems/ari/client/native"

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

type stateFn func(ctx context.Context) (stateFn, error)

type stateObject struct {
	dSub    ari.Subscription
	hSub    ari.Subscription
	oTimer  *time.Timer
	options *Options
	player  audio.Player
	snds    []string
	retData *Result

	// digitReceived is a channel on which an event is sent every time a digit is received by the monitor.
	// It is used by the state functions to step to the next state, where DTMF receipt does so.
	digitReceived chan struct{}
}

// Prompt plays the given sound and waits for user input.
func Prompt(ctx context.Context, p audio.Player, opts *Options, sounds ...string) (ret *Result, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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

	s := &stateObject{
		dSub:          p.Subscribe(ari.Events.ChannelDtmfReceived),
		hSub:          p.Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed),
		oTimer:        time.NewTimer(opts.OverallTimeout),
		options:       opts,
		player:        p,
		snds:          sounds,
		retData:       &Result{},
		digitReceived: make(chan struct{}),
	}

	// Start with a 'Canceled' status, in order to catch early context cancellations
	s.retData.Status = Canceled

	// Monitor the DTMF input
	go s.monitor(ctx, cancel)

	st := s.playPrompt

	for st != nil {
		if err = ctx.Err(); err != nil {
			break
		}

		st, err = st(ctx)
		if err != nil {
			Logger.Error("failure in prompt state machine", "error", err)
			break
		}
	}

	return s.retData, err
}

func (s *stateObject) playPrompt(ctx context.Context) (stateFn, error) {
	playCtx, playCancel := context.WithCancel(context.Background())
	opts := s.options

	defer playCancel()

	q := audio.NewQueue()
	q.Add(s.snds...)
	//var st audio.Status
	st, err := q.Play(playCtx, s.player, &audio.Options{ExitOnDTMF: ""})

	switch st {
	case audio.Canceled:
		//NOTE: since playCtx doesn't extend the parent context,
		// any audio cancel is considered a special case
		// and we don't overwrite the return status
		err = nil
	case audio.Hangup:
		s.retData.Status = Hangup
		return nil, err
	case audio.Failed, audio.Timeout:
		s.retData.Status = Failed
		return nil, err
	}

	// Stop and reset the overall timer for Prompt
	if s.oTimer.Stop() {
		select {
		case <-s.oTimer.C:
		default:
		}
	}
	s.oTimer.Reset(opts.OverallTimeout)

	return s.waitDigit, err
}

func (s *stateObject) monitor(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		select {
		case _, ok := <-s.hSub.Events():
			if ok {
				s.retData.Status = Hangup
				return
			}
		case <-ctx.Done():
			return
		case e := <-s.dSub.Events():
			v, ok := e.(*ari.ChannelDtmfReceived)
			if !ok {
				continue
			}
			s.retData.Data += v.Digit
			Logger.Debug("DTMF received", "digit", v.Digit, "digits", s.retData.Data)

			var status Status
			s.retData.Data, status = s.options.MatchFunc(s.retData.Data)
			s.digitReceived <- struct{}{}
			if status > 0 {
				s.retData.Status = status
				return
			}
		}
	}
}

func (s *stateObject) waitDigit(ctx context.Context) (stateFn, error) {
	lTimeout := s.options.InterDigitTimeout
	if s.retData.Data == "" {
		lTimeout = s.options.FirstDigitTimeout
	}

	select {
	case <-ctx.Done():
		return nil, nil
	case <-time.After(lTimeout):
		s.retData.Status = Timeout
		return nil, nil
	case <-s.oTimer.C:
		s.retData.Status = Timeout
		return nil, nil
	case <-s.digitReceived:
	}
	return s.waitDigit, nil
}
