package audio

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/internal/testutils"
	"github.com/golang/mock/gomock"

	"golang.org/x/net/context"
)

func TestPlayAsync(t *testing.T) {

	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	pb, err := PlayAsync(ctx, bus, player, "audio:hello-world")
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if pb == nil {
		t.Errorf("Expected playback object to be non-nil")
		return
	}

	if pb.Handle() == nil {
		t.Errorf("Expected playback.Handle to be non-nil")
	}

	select {
	case <-pb.StartCh():
		t.Errorf("Unexpected trigger of Start channel")
	case <-pb.StopCh():
		t.Errorf("Unexpected trigger of Stop channel")
	case <-time.After(1 * time.Second):
	}

	// wait for timeout
	<-time.After(MaxPlaybackTime)

	select {
	case <-pb.StartCh():
	default:
		t.Errorf("Expected trigger of start channel after MaxPlaybackTime")
	}

	select {
	case <-pb.StopCh():
	default:
		t.Errorf("Expected trigger of stop channel after MaxPlaybackTime")
	}

	if err := pb.Err(); err == nil {
		t.Errorf("Expected non-nil error")
	}
}

func TestPlayTimeoutStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	err := Play(ctx, bus, player, "audio:hello-world")

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for start of playback" {
		t.Errorf("Expected timeout waiting for start of playback error, got: '%v'", err)
	}
}

func TestPlayTimeoutStop(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		ch <- playbackStartedGood("pb1")
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for stop of playback" {
		t.Errorf("Expected timeout waiting for stop of playback error, got: '%v'", err)
	}
}

func TestPlaySuccess(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayNilEvents(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		ch <- nil
		ch <- playbackStartedGood("pb1")
		ch <- nil
		ch <- nil
		ch <- playbackFinishedGood("pb1")
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayUnrelatedEvents(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		ch <- playbackFinishedDifferentPlaybackID
		ch <- playbackStartedDifferentPlaybackID
		ch <- playbackStartedGood("pb1")

		<-time.After(1 * time.Second)

		ch <- playbackFinishedDifferentPlaybackID
		ch <- playbackFinishedGood("pb1")
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayStopBeforeStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	go func() {
		ch <- playbackFinishedGood("pb1")
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}
}

func TestPlayContextCancellation(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	cancel()

	err := Play(ctx, bus, player, "audio:hello-world")

	<-time.After(1 * time.Millisecond) // causes the other goroutines to 'wake up' and see the cancellation, which causes proper cleanup

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "context canceled" { //TODO: should be an interface to cast to here instead of string comparison
		t.Errorf("Expected context cancellation error, got '%v'", err)
	}
}

func TestPlayContextCancellation100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestPlayContextCancellation(t)
	}
}

func TestPlayContextCancellationAfterStart(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().AnyTimes().Return(ch)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)

	go func() {
		ch <- playbackStartedGood("pb1")
		cancel()
	}()

	err := Play(ctx, bus, player, "audio:hello-world")

	<-time.After(1 * time.Millisecond) // causes the other goroutines to 'wake up' and see the cancellation, which causes proper cleanup

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "context canceled" { //TODO: should be an interface to cast to here instead of string comparison
		t.Errorf("Expected context cancellation error, got '%v'", err)
	}
}

func TestPlayContextCancellationAfterStart100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestPlayContextCancellationAfterStart(t)
	}
}

func TestErrorInPlayer(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(nil, errors.New("Dummy error playing to dummy player"))

	err := Play(ctx, bus, player, "audio:hello-world")

	<-time.After(1 * time.Millisecond) // causes the other goroutines to 'wake up' and see the cancellation, which causes proper cleanup

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "Dummy error playing to dummy player" {
		t.Errorf("Expected dummy error, got '%v'", err)
	}
}

func TestErrorInDataFetch(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Return(sub)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1", failData: true}), nil)

	err := Play(ctx, bus, player, "audio:hello-world")

	<-time.After(1 * time.Millisecond) // causes the other goroutines to 'wake up' and see the cancellation, which causes proper cleanup

	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "Dummy error getting playback data" {
		t.Errorf("Expected dummy error, got '%v'", err)
	}
}

// messages

var channelDtmf = func(dtmf string) ari.Event {
	return &ari.ChannelDtmfReceived{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "ChannelDtmfReceived",
			},
		},
		Digit: dtmf,
	}
}

var playbackStartedGood = func(id string) ari.Event {
	return &ari.PlaybackStarted{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "PlaybackStarted",
			},
		},
		Playback: ari.PlaybackData{
			ID: id,
		},
	}
}

var playbackFinishedGood = func(id string) ari.Event {
	return &ari.PlaybackFinished{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "PlaybackFinished",
			},
		},
		Playback: ari.PlaybackData{
			ID: id,
		},
	}
}

var playbackStartedBadMessageType = &ari.PlaybackStarted{
	EventData: ari.EventData{
		Message: ari.Message{
			Type: "PlaybackStarted2",
		},
	},
	Playback: ari.PlaybackData{
		ID: "pb1",
	},
}

var playbackFinishedBadMessageType = &ari.PlaybackFinished{
	EventData: ari.EventData{
		Message: ari.Message{
			Type: "PlaybackFinished2",
		},
	},
	Playback: ari.PlaybackData{
		ID: "pb1",
	},
}

var playbackStartedDifferentPlaybackID = &ari.PlaybackStarted{
	EventData: ari.EventData{
		Message: ari.Message{
			Type: "PlaybackStarted",
		},
	},
	Playback: ari.PlaybackData{
		ID: "pb2",
	},
}

var playbackFinishedDifferentPlaybackID = &ari.PlaybackFinished{
	EventData: ari.EventData{
		Message: ari.Message{
			Type: "PlaybackFinished",
		},
	},
	Playback: ari.PlaybackData{
		ID: "pb2",
	},
}

// test playback ari transport

type testPlayback struct {
	id       string
	failData bool
}

func (p *testPlayback) Get(id string) *ari.PlaybackHandle {
	panic("not implemented")
}

func (p *testPlayback) Data(id string) (pd ari.PlaybackData, err error) {
	if p.failData {
		err = errors.New("Dummy error getting playback data")
	}
	pd.ID = p.id
	return
}

func (p *testPlayback) Control(id string, op string) error {
	panic("not implemented")
}

func (p *testPlayback) Stop(id string) error {
	panic("not implemented")
}

func isTimeout(err error) bool {

	type timeout interface {
		Timeout() bool
	}

	te, ok := err.(timeout)
	return ok && te.Timeout()
}
