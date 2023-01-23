package play

import (
	"context"

	"github.com/CyCoreSystems/ari/v6"
)

// AllDTMF is a string which contains all possible
// DTMF digits.
const AllDTMF = "0123456789ABCD*#"

// NewPlay creates a new audio Options suitable for general audio playback
func NewPlay(ctx context.Context, p ari.Player, opts ...OptionFunc) (*Options, error) {
	o := NewDefaultOptions()
	err := o.ApplyOptions(opts...)

	return o, err
}

// NewPrompt creates a new audio Options suitable for prompt-style playback-and-get-response situations
func NewPrompt(ctx context.Context, p ari.Player, opts ...OptionFunc) (*Options, error) {
	o := NewPromptOptions()
	err := o.ApplyOptions(opts...)

	return o, err
}

// Play plays a set of media URIs.  Pass these URIs in with the `URI` OptionFunc.
func Play(ctx context.Context, p ari.Player, opts ...OptionFunc) Session {
	o, err := NewPlay(ctx, p, opts...)
	if err != nil {
		return errorSession(err)
	}

	return o.Play(ctx, p)
}

// Prompt plays the given media URIs and waits for a DTMF response.  The
// received DTMF is available as `DTMF` in the Result object.  Further
// customize the behaviour with Match type OptionFuncs.
func Prompt(ctx context.Context, p ari.Player, opts ...OptionFunc) Session {
	o, err := NewPrompt(ctx, p, opts...)
	if err != nil {
		return errorSession(err)
	}

	return o.Play(ctx, p)
}

// Play starts a new Play Session from the existing Options
func (o *Options) Play(ctx context.Context, p ari.Player) Session {
	s := newPlaySession(o)

	go s.play(ctx, p)

	return s
}
