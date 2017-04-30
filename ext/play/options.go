package play

import (
	"container/list"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	// DefaultPlaybackStartTimeout is the default amount of time to wait for a playback to start before declaring that the playback has failed.
	DefaultPlaybackStartTimeout = 2 * time.Second

	// DefaultMaxPlaybackTime is the default maximum amount of time any playback is allowed to run.  If this time is exeeded, the playback will be cancelled.
	DefaultMaxPlaybackTime = 10 * time.Minute

	// DefaultFirstDigitTimeout is the default amount of time to wait, after the playback for all audio completes, for the first digit to be received.
	DefaultFirstDigitTimeout = 4 * time.Second

	// DefaultInterDigitTimeout is the maximum time to wait for additional
	// digits after the first is received.
	DefaultInterDigitTimeout = 3 * time.Second

	// DefaultOverallDigitTimeout is the default maximum time to wait for a
	// response, after the playback for all audio is complete, regardless of the
	// number of received digits or pattern matching.
	DefaultOverallDigitTimeout = 3 * time.Minute

	// DigitBufferSize is the number of digits stored in the received-digit
	// event buffer before further digit events are ignored.  NOTE that digits
	// overflowing this buffer are still stored in the digits received buffer.
	// This only affects the digit _signaling_ buffer.
	DigitBufferSize = 20
)

// Result describes the result of a playback operation
type Result struct {
	mu sync.Mutex

	// Duration indicates how long the playback execution took, from start to finish
	Duration time.Duration

	// DTMF records any DTMF which was received by the playback, as modified by any match functions
	DTMF string

	// Error indicates any error encountered which caused the termination of the playback
	Error error

	// MatchResult indicates the final result of any applied match function for DTMF digits which were received
	MatchResult MatchResult

	// Status indicates the resulting status of the playback, why it was stopped
	Status Status
}

// Status indicates the final status of a playback, be it individual of an entire sequence.  This Status indicates the reason the playback stopped.
type Status int

const (
	// InProgress indicates that the audio is currently playing or is staged to play
	InProgress Status = iota

	// Cancelled indicates that the audio was cancelled.  This cancellation could be due
	// to anything from the control context being closed or a DTMF Match being found
	Cancelled

	// Failed indicates that the audio playback failed.  This indicates that one
	// or more of the audio playbacks failed to be played.  This could be due to
	// a system, network, or Asterisk error, but it could also be due to an
	// invalid audio URI.  Check the returned error for more details.
	Failed

	// Finished indicates that the playback completed playing all bound audio
	// URIs in full.  Note that for a prompt-style execution, this also means
	// that no DTMF was matched to the match function.
	Finished

	// Hangup indicates that the audio playback was interrupted due to a hangup.
	Hangup

	// Timeout indicates that audio playback timed out.  It is not known whether this was due to a failure in the playback, a network loss, or some other problem.
	Timeout
)

// MatchResult indicates the status of a match for the received DTMF of a playback
type MatchResult int

const (
	// Incomplete indicates that there are not enough digits to determine a match
	Incomplete MatchResult = iota

	// Complete indicates that a match was found and the current DTMF pattern is complete
	Complete

	// Invalid indicates that a match cannot befound from the current DTMF received set
	Invalid
)

type uriList struct {
	list    *list.List
	current *list.Element
	mu      sync.Mutex
}

func (u *uriList) Empty() bool {
	if u == nil || u.list == nil || u.list.Len() == 0 {
		return true
	}
	return false
}

func (u *uriList) Add(uri string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.list == nil {
		u.list = list.New()
	}

	u.list.PushBack(uri)

	if u.current == nil {
		u.current = u.list.Front()
	}
}

func (u *uriList) First() string {
	if u.list == nil {
		return ""
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.current = u.list.Front()
	return u.val()
}

func (u *uriList) Next() string {
	if u.list == nil {
		return ""
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.current == nil {
		return ""
	}

	u.current = u.current.Next()
	return u.val()
}

func (u *uriList) val() string {
	if u.current == nil {
		return ""
	}

	ret, ok := u.current.Value.(string)
	if !ok {
		return ""
	}
	return ret
}

// Options represent the various playback options which can modify the operation of a Playback.
type Options struct {
	// uriList is the list of audio URIs to play
	uriList *uriList

	// playbackStartTimeout defines the amount of time to wait for a playback to
	// start before declaring it failed.
	//
	// This value is important because ARI does NOT report playback failures in
	// any usable way.
	//
	// If not specified, the default is DefaultPlaybackStartTimeout
	playbackStartTimeout time.Duration

	// maxPlaybackTime is the maximum amount of time to wait for a playback
	// session to complete, everything included.  The playback will be
	// terminated if this time is exceeded.
	maxPlaybackTime time.Duration

	// firstDigitTimeout is the maximum length of time to wait
	// after the prompt sequence ends for the user to enter
	// a response.
	//
	// If not specified, the default is DefaultFirstDigitTimeout.
	firstDigitTimeout time.Duration

	// interDigitTimeout is the maximum length of time to wait
	// for an additional digit after a digit is received.
	//
	// If not specified, the default is DefaultInterDigitTimeout.
	interDigitTimeout time.Duration

	// overallDigitTimeout is the maximum length of time to wait
	// for a response regardless of digits received after the completion
	// of all audio playbacks.
	// If not specified, the default is DefaultOverallTimeout.
	overallDigitTimeout time.Duration

	// matchFunc is an optional function which, if supplied, returns
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
	matchFunc func(string) (string, MatchResult)

	// maxReplays is the maximum number of times the audio sequence will be
	// replayed if there is no response.  By default, the audio sequence is
	// played only once.
	maxReplays int
}

// NewDefaultOptions returns a set of options which represent reasonable defaults for most simple playbacks.
func NewDefaultOptions() *Options {
	opts := &Options{
		playbackStartTimeout: DefaultPlaybackStartTimeout,
		maxPlaybackTime:      DefaultMaxPlaybackTime,
		uriList:              new(uriList),
	}

	MatchAny()(opts) // nolint  No error is possible with MatchAny

	return opts
}

// ApplyOptions applies a set of OptionFuncs to the Playback
func (o *Options) ApplyOptions(opts ...OptionFunc) (err error) {
	for _, f := range opts {
		err = f(o)
		if err != nil {
			return errors.Wrap(err, "failed to apply option")
		}
	}
	return nil
}

// NewPromptOptions returns a set of options which represent reasonable defaults for most prompt playbacks.  It will terminate when any single DTMF digit is received.
func NewPromptOptions() *Options {
	opts := NewDefaultOptions()

	opts.firstDigitTimeout = DefaultFirstDigitTimeout
	opts.interDigitTimeout = DefaultInterDigitTimeout
	opts.overallDigitTimeout = DefaultOverallDigitTimeout

	return opts
}

// OptionFunc defines an interface for functions which can modify a play session's Options
type OptionFunc func(*Options) error

// NoExitOnDTMF disables exiting the playback when DTMF is received.  Note that
// this is just a wrapper for MatchFunc(nil), so it is mutually exclusive with
// MatchFunc; whichever comes later will win.
func NoExitOnDTMF() OptionFunc {
	return func(o *Options) error {
		o.matchFunc = nil
		return nil
	}
}

// URI adds a set of audio URIs to a playback
func URI(uri ...string) OptionFunc {
	return func(o *Options) error {
		if o.uriList == nil {
			o.uriList = new(uriList)
		}

		for _, u := range uri {
			if u != "" {
				o.uriList.Add(u)
			}
		}
		return nil
	}
}

// PlaybackStartTimeout overrides the default playback start timeout
func PlaybackStartTimeout(timeout time.Duration) OptionFunc {
	return func(o *Options) error {
		o.playbackStartTimeout = timeout
		return nil
	}
}

// DigitTimeouts sets the digit timeouts.  Passing a negative value to any of these indicates that the default value (shown in parentheses below) should be used.
//
//  - First digit timeout (4 sec):  The time (after the stop of the audio) to wait for the first digit to be received
//
//  - Inter digit timeout (3 sec):  The time (after receiving a digit) to wait for the _next_ digit to be received
//
//  - Overall digit timeout (3 min):  The maximum amount of time to wait (after the stop of the audio) for digits to be received, regardless of the digit frequency
//
func DigitTimeouts(first, inter, overall time.Duration) OptionFunc {
	return func(o *Options) error {
		if first >= 0 {
			o.firstDigitTimeout = first
		}
		if inter >= 0 {
			o.interDigitTimeout = inter
		}
		if overall >= 0 {
			o.overallDigitTimeout = overall
		}
		return nil
	}
}

// Replays sets the number of replays of the audio sequence before exiting
func Replays(count int) OptionFunc {
	return func(o *Options) error {
		o.maxReplays = count
		return nil
	}
}

// MatchAny indicates that the playback should be considered Matched and terminated if
// any DTMF digit is received during the playback or post-playback time.
func MatchAny() OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			if len(pat) > 0 {
				return pat, Complete
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchDiscrete indicates that the playback should be considered Matched and terminated if
// the received DTMF digits match any of the discrete list of strings.
func MatchDiscrete(list []string) OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			var maxLen int
			for _, t := range list {
				if t == pat {
					return pat, Complete
				}
				if len(t) > maxLen {
					maxLen = len(t)
				}
			}
			if len(pat) > maxLen {
				return pat, Invalid
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchHash indicates that the playback should be considered Matched and terminated if it contains a hash (#).  The hash (and any subsequent digits) is removed from the final result.
func MatchHash() OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			if strings.Contains(pat, "#") {
				return strings.Split(pat, "#")[0], Complete
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchTerminator indicates that the playback shoiuld be considered Matched and terminated if it contains the provided Terminator string.  The terminator (and any subsequent digits) is removed from the final result.
func MatchTerminator(terminator string) OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			if strings.Contains(pat, terminator) {
				return strings.Split(pat, terminator)[0], Complete
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchLen indicates that the playback should be considered Matched and terminated if the given number of DTMF digits are receieved.
func MatchLen(length int) OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			if len(pat) >= length {
				return pat, Complete
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchLenOrTerminator indicates that the playback should be considered Matched and terminated if the given number of DTMF digits are receieved or if the given terminator is received.  If the terminator is present, it and any subsequent digits will be removed from the final result.
func MatchLenOrTerminator(length int, terminator string) OptionFunc {
	return func(o *Options) error {
		o.matchFunc = func(pat string) (string, MatchResult) {
			if len(pat) >= length {
				return pat, Complete
			}
			if strings.Contains(pat, terminator) {
				return strings.Split(pat, terminator)[0], Complete
			}
			return pat, Incomplete
		}
		return nil
	}
}

// MatchFunc uses the provided match function to determine when the playback should be terminated based on DTMF input.
func MatchFunc(f func(string) (string, MatchResult)) OptionFunc {
	return func(o *Options) error {
		o.matchFunc = f
		return nil
	}
}
