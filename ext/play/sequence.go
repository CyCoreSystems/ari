package play

import (
	"context"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/rid"
	"github.com/pkg/errors"
)

// sequence represents an audio sequence playback session
type sequence struct {
	cancel context.CancelFunc
	s      *playSession

	done chan struct{}
}

func (s *sequence) Done() <-chan struct{} {
	return s.done
}

func (s *sequence) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func newSequence(s *playSession) *sequence {
	return &sequence{
		s:    s,
		done: make(chan struct{}),
	}
}

func (s *sequence) Play(ctx context.Context, p ari.Player) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	defer cancel()
	defer close(s.done)

	for u := s.s.o.uriList.First(); u != ""; u = s.s.o.uriList.Next() {
		pb, err := p.StagePlay(rid.New(rid.Playback), u)
		if err != nil {
			s.s.result.Status = Failed
			s.s.result.Error = errors.Wrap(err, "failed to stage playback")

			return
		}

		s.s.result.Status, err = playStaged(ctx, pb, s.s.o.playbackStartTimeout)
		if err != nil {
			s.s.result.Error = errors.Wrap(err, "failure in playback")

			return
		}
	}
}

// playStaged executes a staged playback, waiting for its completion
func playStaged(ctx context.Context, h *ari.PlaybackHandle, timeout time.Duration) (Status, error) {
	started := h.Subscribe(ari.Events.PlaybackStarted)
	defer started.Cancel()

	finished := h.Subscribe(ari.Events.PlaybackFinished)
	defer finished.Cancel()

	if timeout == 0 {
		timeout = DefaultPlaybackStartTimeout
	}

	if err := h.Exec(); err != nil {
		return Failed, errors.Wrap(err, "failed to start playback")
	}

	defer h.Stop() // nolint: errcheck

	select {
	case <-ctx.Done():
		return Cancelled, nil
	case <-started.Events():
	case <-finished.Events():
		return Finished, nil
	case <-time.After(timeout):
		return Timeout, errors.New("timeout waiting for playback to start")
	}

	// Wait for playback to complete
	select {
	case <-ctx.Done():
		return Cancelled, nil
	case <-finished.Events():
		return Finished, nil
	}
}
