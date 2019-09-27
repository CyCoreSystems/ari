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

type playStagedTest struct {
	playbackStartedChan chan ari.Event
	playbackStarted     *arimocks.Subscription

	playbackEndChan chan ari.Event
	playbackEnd     *arimocks.Subscription

	handleExec func(_ *ari.PlaybackHandle) error

	playback *arimocks.Playback

	key *ari.Key

	handle *ari.PlaybackHandle
}

func (p *playStagedTest) Setup() {
	p.playbackStarted = &arimocks.Subscription{}
	p.playbackEnd = &arimocks.Subscription{}
	p.playback = &arimocks.Playback{}

	p.key = ari.NewKey(ari.PlaybackKey, "ph1")

	p.playbackStartedChan = make(chan ari.Event)
	p.playbackStarted.On("Events").Return((<-chan ari.Event)(p.playbackStartedChan))

	p.playbackStarted.On("Cancel").Times(1).Return(nil)
	p.playback.On("Subscribe", p.key, ari.Events.PlaybackStarted).Return(p.playbackStarted)
	p.playback.On("Stop", p.key).Times(1).Return(nil)

	p.playbackEndChan = make(chan ari.Event)
	p.playbackEnd.On("Events").Return((<-chan ari.Event)(p.playbackEndChan))
	p.playbackEnd.On("Cancel").Times(1).Return(nil)
	p.playback.On("Subscribe", p.key, ari.Events.PlaybackFinished).Return(p.playbackEnd)

	p.handle = ari.NewPlaybackHandle(p.key, p.playback, p.handleExec)
}

func TestPlayStaged(t *testing.T) {
	t.Run("noEventTimeout", testPlayStagedNoEventTimeout)
	t.Run("startFinishedEvent", testPlayStagedStartFinishedEvent)
	t.Run("finishedBeforeStart", testPlayStagedFinishedEvent)
	t.Run("failExec", testPlayStagedFailExec)
	t.Run("cancel", testPlayStagedCancel)
	t.Run("cancelAfterStart", testPlayStagedCancelAfterStart)
}

func testPlayStagedNoEventTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	st, err := playStaged(ctx, p.handle, 0)
	if err == nil || err.Error() != "timeout waiting for playback to start" {
		t.Errorf("Expected error '%v', got '%v'", "timeout waiting for playback to start", err)
	}
	if st != Timeout {
		t.Errorf("Expected status '%v', got '%v'", st, Timeout)
	}
}

func testPlayStagedStartFinishedEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		p.playbackStartedChan <- &ari.PlaybackStarted{}
		time.After(10 * time.Millisecond)
		p.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	st, err := playStaged(ctx, p.handle, 20*time.Millisecond)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Finished {
		t.Errorf("Expected status '%v', got '%v'", st, Finished)
	}
}

func testPlayStagedFinishedEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		p.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	st, err := playStaged(ctx, p.handle, 0)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Finished {
		t.Errorf("Expected status '%v', got '%v'", st, Finished)
	}
}

func testPlayStagedFailExec(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.handleExec = func(_ *ari.PlaybackHandle) error {
		return errors.New("err2")
	}
	p.Setup()

	st, err := playStaged(ctx, p.handle, 0)
	if err == nil || err.Error() != "failed to start playback: err2" {
		t.Errorf("Expected error '%v', got '%v'", "failed to start playback: err2", err)
	}
	if st != Failed {
		t.Errorf("Expected status '%v', got '%v'", st, Failed)
	}
}

// nolint
func XXXtestPlayStagedFinishBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		time.After(100 * time.Millisecond)
		p.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	st, err := playStaged(ctx, p.handle, 0)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Finished {
		t.Errorf("Expected status '%v', got '%v'", Finished, st)
	}
}

func testPlayStagedCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		<-time.After(10 * time.Millisecond)
		cancel()
	}()

	st, err := playStaged(ctx, p.handle, 20*time.Millisecond)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", Cancelled, st)
	}
}

func testPlayStagedCancelAfterStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		p.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(100 * time.Millisecond)
		cancel()
	}()

	st, err := playStaged(ctx, p.handle, 0)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", Cancelled, st)
	}
}

type playTest struct {
	ps playStagedTest

	dtmfChannel    chan ari.Event
	dtmfChannelSub *arimocks.Subscription
	player         *arimocks.Player
}

func (p *playTest) Setup() {
	p.ps.Setup()

	p.dtmfChannel = make(chan ari.Event)
	p.dtmfChannelSub = &arimocks.Subscription{}
	p.dtmfChannelSub.On("Events").Return((<-chan ari.Event)(p.dtmfChannel))
	p.dtmfChannelSub.On("Cancel").Return(nil)

	p.player = &arimocks.Player{}
	p.player.On("Subscribe", ari.Events.ChannelDtmfReceived).Return(p.dtmfChannelSub)
	p.player.On("StagePlay", mock.Anything, "sound:1").Return(p.ps.handle, nil)
}

func TestPlay(t *testing.T) {
	t.Run("testPlayNoURI", testPlayNoURI)
	t.Run("testPlay", testPlay)
	t.Run("testPlayDtmf", testPlayDtmf)
}

func testPlayNoURI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playTest
	p.Setup()

	res, err := Play(ctx, p.player).Result()
	if err == nil || err.Error() != "empty playback URI list" {
		t.Errorf("Expected error '%v', got '%v'", "empty playback URI list", err)
	}
	if res.DTMF != "" {
		t.Errorf("Unexpected DTMF: %s", res.DTMF)
	}
}

func testPlay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playTest
	p.Setup()

	go func() {
		p.ps.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(200 * time.Millisecond)
		p.ps.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	res, err := Play(ctx, p.player, URI("sound:1")).Result()
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if res.Status != Finished {
		t.Errorf("Expected status '%v', got '%v'", Finished, res.Status)
	}
	if res.DTMF != "" {
		t.Errorf("Unexpected DTMF: %s", res.DTMF)
	}
}

func testPlayDtmf(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playTest
	p.Setup()

	go func() {
		p.ps.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(200 * time.Millisecond)

		p.dtmfChannel <- &ari.ChannelDtmfReceived{
			Digit: "1",
		}
		<-time.After(200 * time.Millisecond)

		p.ps.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	res, err := Prompt(ctx, p.player, URI("sound:1")).Result()
	if err != nil {
		t.Error("Unexpected error", err)
	}

	if res.MatchResult != Complete {
		t.Errorf("Expected MatchResult '%v', got '%v'", Complete, res.MatchResult)
	}
	if res.DTMF != "1" {
		t.Errorf("Expected DTMF %s, got DTMF %s", "1", res.DTMF)
	}
}
