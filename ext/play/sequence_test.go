package play

import (
	"context"
	"errors"
	"testing"

	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/arimocks"
	"github.com/stretchr/testify/mock"
)

type sequenceTest struct {
	playbackStartedChan chan ari.Event
	playbackStarted     *arimocks.Subscription

	playbackEndChan chan ari.Event
	playbackEnd     *arimocks.Subscription

	handleExeced bool
	handleExec   func(_ *ari.PlaybackHandle) error

	playback *arimocks.Playback

	key *ari.Key
}

func (p *sequenceTest) Setup(handle string) {

	p.playbackStarted = &arimocks.Subscription{}
	p.playbackEnd = &arimocks.Subscription{}
	p.playback = &arimocks.Playback{}

	p.key = ari.NewKey(ari.PlaybackKey, handle)

	p.playbackStartedChan = make(chan ari.Event)
	p.playbackStarted.On("Events").Return((<-chan ari.Event)(p.playbackStartedChan))

	p.playbackStarted.On("Cancel").Times(1).Return(nil)
	p.playback.On("Subscribe", p.key, ari.Events.PlaybackStarted).Return(p.playbackStarted)
	p.playback.On("Stop", p.key).Times(1).Return(nil)

	p.playbackEndChan = make(chan ari.Event)
	p.playbackEnd.On("Events").Return((<-chan ari.Event)(p.playbackEndChan))
	p.playbackEnd.On("Cancel").Times(1).Return(nil)
	p.playback.On("Subscribe", p.key, ari.Events.PlaybackFinished).Return(p.playbackEnd)
}

func TestSequence(t *testing.T) {
	t.Run("noItems", testSequenceNoItems)
	t.Run("someItemsTimeoutStart", testSequenceSomeItemsTimeoutStart)
	t.Run("someItems", testSequenceSomeItems)
	t.Run("someItemsStopEarly", testSequenceSomeItemsStopEarly)
	t.Run("someItemsCancelEarly", testSequenceSomeItemsCancelEarly)
	t.Run("someItemsStagePlayFailure", testSequenceSomeItemsStagePlayFailure)
}

func testSequenceNoItems(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seq := newSequence(NewDefaultOptions())

	seq.Play(ctx, player)

	player.AssertNotCalled(t, "StagePlay")
}

func testSequenceSomeItemsTimeoutStart(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var s sequenceTest
	s.Setup("ph1")

	var s2 sequenceTest
	s2.Setup("ph2")

	ph1 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph1"), s.playback, nil)
	ph2 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph2"), s2.playback, nil)

	player.On("StagePlay", mock.Anything, "sound:1").Return(ph1, nil)
	player.On("StagePlay", mock.Anything, "sound:2").Return(ph2, nil)

	opts := NewDefaultOptions()
	opts.uriList.Add("sound:1")
	opts.uriList.Add("sound:2")
	seq := newSequence(opts)

	seq.Play(ctx, player)

	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:1")
	player.AssertNotCalled(t, "StagePlay", mock.Anything, "sound:2")

	if err := seq.opts.result.Err(); err == nil || err.Error() != "failure in playback: timeout waiting for playback to start" {
		t.Errorf("Expected error: %s, got %v", "failure in playback: timeout waiting for playback to start", err)
	}
	if seq.opts.result.Status != Timeout {
		t.Errorf("Expected status '%v', got '%v'", Timeout, seq.opts.result.Status)
	}
}

func testSequenceSomeItems(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var s sequenceTest
	s.Setup("ph1")

	var s2 sequenceTest
	s2.Setup("ph2")

	ph1 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph1"), s.playback, nil)
	ph2 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph2"), s2.playback, nil)

	player.On("StagePlay", mock.Anything, "sound:1").Return(ph1, nil)
	player.On("StagePlay", mock.Anything, "sound:2").Return(ph2, nil)

	opts := NewDefaultOptions()
	opts.uriList.Add("sound:1")
	opts.uriList.Add("sound:2")
	seq := newSequence(opts)

	go func() {
		s.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)
		s.playbackEndChan <- &ari.PlaybackFinished{}

		<-time.After(20 * time.Millisecond)

		s2.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)
		s2.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	seq.Play(ctx, player)

	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:1")
	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:2")

	if seq.opts.result.Err() != nil {
		t.Errorf("Unexpected error: %v", seq.opts.result.Err())
	}
	if seq.opts.result.Status != Finished {
		t.Errorf("Expected status '%v', got '%v'", Finished, seq.opts.result.Status)
	}
}

func testSequenceSomeItemsCancelEarly(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var s sequenceTest
	s.Setup("ph1")

	var s2 sequenceTest
	s2.Setup("ph2")

	ph1 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph1"), s.playback, nil)
	ph2 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph2"), s2.playback, nil)

	player.On("StagePlay", mock.Anything, "sound:1").Return(ph1, nil)
	player.On("StagePlay", mock.Anything, "sound:2").Return(ph2, nil)

	opts := NewDefaultOptions()
	opts.uriList.Add("sound:1")
	opts.uriList.Add("sound:2")
	seq := newSequence(opts)

	go func() {
		s.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)
		s.playbackEndChan <- &ari.PlaybackFinished{}

		<-time.After(20 * time.Millisecond)

		s2.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)

		cancel()

	}()

	seq.Play(ctx, player)

	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:1")
	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:2")

	if seq.opts.result.Err() != nil {
		t.Errorf("Unexpected error: %v", seq.opts.result.Err())
	}
	if seq.opts.result.Status != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", Cancelled, seq.opts.result.Status)
	}
}

func testSequenceSomeItemsStopEarly(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var s sequenceTest
	s.Setup("ph1")

	var s2 sequenceTest
	s2.Setup("ph2")

	ph1 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph1"), s.playback, nil)
	ph2 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph2"), s2.playback, nil)

	player.On("StagePlay", mock.Anything, "sound:1").Return(ph1, nil)
	player.On("StagePlay", mock.Anything, "sound:2").Return(ph2, nil)

	opts := NewDefaultOptions()
	opts.uriList.Add("sound:1")
	opts.uriList.Add("sound:2")
	seq := newSequence(opts)

	go func() {
		s.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)
		s.playbackEndChan <- &ari.PlaybackFinished{}

		<-time.After(20 * time.Millisecond)

		s2.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)

		seq.Stop()

	}()

	seq.Play(ctx, player)

	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:1")
	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:2")

	if seq.opts.result.Err() != nil {
		t.Errorf("Unexpected error: %v", seq.opts.result.Err())
	}
	if seq.opts.result.Status != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", Cancelled, seq.opts.result.Status)
	}
}

func testSequenceSomeItemsStagePlayFailure(t *testing.T) {
	player := &arimocks.Player{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var s sequenceTest
	s.Setup("ph1")

	var s2 sequenceTest
	s2.Setup("ph2")

	ph1 := ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "ph1"), s.playback, nil)

	player.On("StagePlay", mock.Anything, "sound:1").Return(ph1, nil)
	player.On("StagePlay", mock.Anything, "sound:2").Return(nil, errors.New("unknown error"))

	opts := NewDefaultOptions()
	opts.uriList.Add("sound:1")
	opts.uriList.Add("sound:2")
	seq := newSequence(opts)

	go func() {
		s.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(20 * time.Millisecond)
		s.playbackEndChan <- &ari.PlaybackFinished{}

		<-time.After(20 * time.Millisecond)
	}()

	seq.Play(ctx, player)

	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:1")
	player.AssertCalled(t, "StagePlay", mock.Anything, "sound:2")

	if err := seq.opts.result.Err(); err == nil || err.Error() != "failed to stage playback: unknown error" {
		t.Errorf("Expected error: %v, got %v", "failed to stage playback: unknown error", err)
	}
	if seq.opts.result.Status != Failed {
		t.Errorf("Expected status '%v', got '%v'", Failed, seq.opts.result.Status)
	}
}
