package audio

import (
	"context"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// sequence represents an audio sequence playback session
type sequence struct {
	cancel context.CancelFunc
	opts   *Options

	done chan error
}

func (s *sequence) Done() <-chan error {
	return s.done
}

func (s *sequence) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func newSequence(o *Options) *sequence {
	return &sequence{
		opts: o,
		done: make(chan error),
	}
}

func (s *sequence) Play(ctx context.Context, p ari.Player) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	defer cancel()

	for u := s.opts.uriList.First(); u != ""; u = s.opts.uriList.Next() {
		pb, err := p.StagePlay(uuid.NewV1().String(), u)
		if err != nil {
			s.opts.result.Status = Failed
			s.opts.result.Error = errors.Wrap(err, "failed to stage playback")
			return
		}

		s.opts.result.Status, err = playStaged(ctx, pb, s.opts)
		if err != nil {
			s.opts.result.Error = errors.Wrap(err, "failure in playback")
			return
		}
	}
}
