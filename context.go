package ari

import "context"

// ChannelContextOptions describes the set of options to be used when creating a channel-bound context.
type ChannelContextOptions struct {
	ctx    context.Context
	cancel context.CancelFunc

	hangupOnEnd bool

	sub Subscription
}

// ChannelContextOptionFunc describes a function which modifies channel context options.
type ChannelContextOptionFunc func(o *ChannelContextOptions)

// ChannelContext returns a context which is closed when the provided channel leaves the ARI application or the parent context is closed.  The parent context is optional, and if it is `nil`, a new background context will be created.
func ChannelContext(h *ChannelHandle, opts ...ChannelContextOptionFunc) (context.Context, context.CancelFunc) {

	o := new(ChannelContextOptions)
	for _, opt := range opts {
		opt(o)
	}

	if o.ctx == nil {
		o.ctx, o.cancel = context.WithCancel(context.Background())
	}

	if o.sub == nil {
		o.sub = h.Subscribe(Events.StasisEnd)
	}

	go func() {
		defer o.cancel()

		select {
		case <-o.ctx.Done():
		case <-o.sub.Events():
		}

		if o.hangupOnEnd {
			h.Hangup()
		}
	}()

	return o.ctx, o.cancel
}

// WithParentContext requests that the generated channel context be created from the given parent context.
func WithParentContext(parent context.Context) ChannelContextOptionFunc {
	return func(o *ChannelContextOptions) {
		if parent != nil {
			o.ctx, o.cancel = context.WithCancel(parent)
		}
	}
}

// HangupOnEnd indicates that the channel should be terminated when the channel context is terminated.  Note that this also provides an easy way to create a time scope on a channel by putting a deadline on the parent context.
func HangupOnEnd() ChannelContextOptionFunc {
	return func(o *ChannelContextOptions) {
		o.hangupOnEnd = true
	}
}
