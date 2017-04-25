package play

import (
	"context"
	"errors"
	"testing"

	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/arimocks"
)

type playStagedTest struct {
	playbackStartedChan chan ari.Event
	playbackStarted     *arimocks.Subscription

	playbackEndChan chan ari.Event
	playbackEnd     *arimocks.Subscription

	handleExeced bool
	handleExec   func(_ *ari.PlaybackHandle) error

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

type timeoutTest struct {
	playStagedTest
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

	st, err := playStaged(ctx, p.handle, nil)
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
		time.After(200 * time.Millisecond)
		p.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	st, err := playStaged(ctx, p.handle, nil)
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

	st, err := playStaged(ctx, p.handle, nil)
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

	st, err := playStaged(ctx, p.handle, nil)
	if err == nil || err.Error() != "failed to start playback: err2" {
		t.Errorf("Expected error '%v', got '%v'", "failed to start playback: err2", err)
	}
	if st != Failed {
		t.Errorf("Expected status '%v', got '%v'", st, Failed)
	}
}

func testPlayStagedFinishBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		time.After(100 * time.Millisecond)
		p.playbackEndChan <- &ari.PlaybackFinished{}
	}()

	st, err := playStaged(ctx, p.handle, nil)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Finished {
		t.Errorf("Expected status '%v', got '%v'", st, Finished)
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

	st, err := playStaged(ctx, p.handle, nil)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", st, Cancelled)
	}
}

func testPlayStagedCancelAfterStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var p playStagedTest
	p.Setup()

	go func() {
		p.playbackStartedChan <- &ari.PlaybackStarted{}
		<-time.After(200 * time.Millisecond)
		cancel()
	}()

	st, err := playStaged(ctx, p.handle, nil)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}
	if st != Cancelled {
		t.Errorf("Expected status '%v', got '%v'", st, Cancelled)
	}
}

/*
func init() {
	PlaybackStartTimeout = 5 * time.Millisecond
}

func TestWaitStartCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var c Control

	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	select {
	case <-time.After(2 * time.Millisecond):
		t.Error("waitStart failed to detect context closure")
	case <-doneCh:
	}

	if c.status != Canceled {
		t.Error("waitStart returned the wrong state")
	}
}

func TestWaitStartTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c Control

	// Prepare mock subscriptions
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	select {
	case <-time.After(PlaybackStartTimeout * 2):
		t.Error("waitStart failed to detect timeout")
	case <-doneCh:
	}

	if c.status != Timeout {
		t.Error("waitStart returned the wrong state")
	}

}

func TestWaitStartHangup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c Control

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = sub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect hangup")
	case <-doneCh:
	}

	if c.status != Hangup {
		t.Error("waitStart returned the wrong state")
	}
}

func TestWaitStartFinished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := Control{
		stopCh: make(chan struct{}),
	}

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel().AnyTimes()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = sub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect playback started")
	case <-doneCh:
	}

	if c.status != Finished {
		t.Error("waitStart returned the wrong state", c.status)
	}
}
func TestWaitStartStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := Control{
		startCh: make(chan struct{}),
	}

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel().AnyTimes()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = sub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f == nil {
			t.Error("waitStart returned a nil state instead of c.waitStop")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect playback started")
	case <-doneCh:
	}

	if c.status != InProgress {
		t.Error("waitStart returned the wrong state", c.status)
	}
}

*/
