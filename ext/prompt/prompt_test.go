package prompt

import (
	"context"
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/ext/audio"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

/*
// TestWaitDigit tests the WaidDigit func
func TestWaitDigit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := Options{
		FirstDigitTimeout: 4 * time.Second,
		InterDigitTimeout: 3 * time.Second,
		OverallTimeout:    0,
		EchoData:          true,
		MatchFunc:         nil,
		SoundHash:         "",
	}
	ret := &Result{}

	ctlr := gomock.NewController(t)
	defer ctlr.Finish()

	//p := mock.NewMockPlayer(ctlr)

	//dtmfSub := p.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived)

	//hangupSub := p.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed)

	overallTimer = time.NewTimer(opts.OverallTimeout)

	//bus := mock.NewMockBus(ctlr)

	sub2 := mock.NewMockSubscription(ctlr)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	go func() {
		ctx2, cancel2 := context.WithTimeout(ctx, time.Millisecond)
		defer cancel2()
		fn, err := waitDigit(ctx2, 3*time.Second, &opts, ret)
		if err != nil ||
			fn == nil {
			t.Errorf("Unexpected error. '%v'", err)
		}
		if ret.Status != Timeout {
			t.Errorf("Expected Timeout but got '%v'", ret.Status)
		}
	}()

	ch2 <- channelDtmf("1")
	<-ctx.Done()
}


func TestPromptPlayError(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	dtmfA := mock.NewMockSubscription(ctrl)
	dtmfB := mock.NewMockSubscription(ctrl)
	hangup := mock.NewMockSubscription(ctrl)

	sub.EXPECT().Cancel().Times(1)
	dtmfA.EXPECT().Cancel().Times(1)
	dtmfB.EXPECT().Cancel().Times(1)
	hangup.EXPECT().Cancel().Times(1)

	gomock.InOrder(
		bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(hangup),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfA),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfB),
		bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Times(1).Return(sub),
	)

	dtmfA.EXPECT().Events().Times(0) // called in prompt
	dtmfB.EXPECT().Events().Times(1) // called in play
	hangup.EXPECT().Events().Times(1)
	sub.EXPECT().Events().Times(0)

	player := mock.NewMockPlayer(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)
	//player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1", failData: true}), nil)
	//player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	res, err := Prompt(ctx, player, nil, "sound:1", "sound:2")

	<-time.After(30 * time.Millisecond)

	if err == nil || err.Error() != "Dummy error getting playback data" {
		t.Errorf("Expected dummy error getting playback error, got: '%v'", err)
	}

	if res.Status != Failed {
		t.Errorf("Expected Failed result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}
}


func TestPromptCancelBeforePromptComplete(t *testing.T) {
	//audio.MaxPlaybackTime = 3 * time.Second

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	dtmfA := mock.NewMockSubscription(ctrl)
	dtmfB := mock.NewMockSubscription(ctrl)
	sub3 := mock.NewMockSubscription(ctrl)

	sub.EXPECT().Cancel().Times(1)
	dtmfA.EXPECT().Cancel().Times(1)
	dtmfB.EXPECT().Cancel().Times(1)
	sub3.EXPECT().Cancel().Times(1)

	gomock.InOrder(
		bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfA),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfB),
		bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Times(1).Return(sub),
	)

	ch := make(chan ari.Event)

	dtmfA.EXPECT().Events().MinTimes(0) // called in prompt
	dtmfB.EXPECT().Events().MinTimes(1) // called in play
	sub3.EXPECT().Events().MinTimes(1)
	sub.EXPECT().Events().MinTimes(1).Return(ch)

	player := mock.NewMockPlayer(ctrl)
	//player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	//player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		//<-player. // play request
		ch <- playbackStartedGood("pb1")
		cancel()
	}()

	res, err := Prompt(ctx, player, nil, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err == nil || err.Error() != "context canceled" {
		t.Errorf("Expected error 'context cancelled', got '%v'", err)
	}

	if res.Status != Canceled {
		t.Errorf("Expected Canceled result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}
}

func TestPromptNoInput(t *testing.T) {
	//audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().Times(2)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := mock.NewMockPlayer(ctrl)
	//player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	//player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next // play request
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")
	}()

	res, err := Prompt(ctx, player, nil, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Timeout {
		t.Errorf("Expected Timeout result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}

}
*/

func TestPromptHangup(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//bus := mock.NewMockBus(ctrl)

	//sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	//sub.EXPECT().Events().MinTimes(1).Return(ch)
	//sub.EXPECT().Cancel().Times(1)

	//bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	/*sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)
	*/
	//bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	//sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	//sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	//sub3.EXPECT().Cancel().Times(1)

	//bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	//hangupChan := make(chan ari.Event)
	//hangupSub := mock.NewMockSubscription(ctrl)
	//dtmfSub := mock.NewMockSubscription(ctrl)
	player := mock.NewMockPlayer(ctrl)
	//player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived)
	//player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	_ = player.Subscribe(ari.Events.ChannelDtmfReceived)

	//defer dtmfSub.Cancel()
	//hangupSub.EXPECT().Cancel()
	//hangupSub.EXPECT().Events().MinTimes(1).Return(hangupChan)

	//player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	//player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)
	//playback := mock.NewMockPlayback(ctrl)
	//handle := ari.NewPlaybackHandle("pb1", playback)
	//handle2 := ari.NewPlaybackHandle("pb2", playback)

	go func() {
		//<-player.Next // play request
		//player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)
		ch <- playbackStartedGood("pb1")
		ch3 <- &ari.ChannelHangupRequest{EventData: ari.EventData{Message: ari.Message{Type: "ChannelHangupRequest"}}}
	}()

	res, err := Prompt(ctx, player, nil, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Hangup {
		t.Errorf("Expected Hangup result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}
	dtmfSub.Cancel()
}

/*
func TestPromptMatchHashEchoData(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	player.Append(ari.NewPlaybackHandle("d1", &testPlayback{id: "d1"}), nil)
	player.Append(ari.NewPlaybackHandle("d2", &testPlayback{id: "d2"}), nil)
	player.Append(ari.NewPlaybackHandle("d3", &testPlayback{id: "d3"}), nil)
	player.Append(ari.NewPlaybackHandle("d4", &testPlayback{id: "d4"}), nil)
	player.Append(ari.NewPlaybackHandle("d5", &testPlayback{id: "d5"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		<-time.After(50 * time.Millisecond)

		ch <- playbackStartedGood("d1")
		ch <- playbackFinishedGood("d1")

		ch2 <- channelDtmf("3")
		<-time.After(50 * time.Millisecond)

		ch <- playbackStartedGood("d2")
		ch <- playbackFinishedGood("d2")

		ch2 <- channelDtmf("1")
		<-time.After(50 * time.Millisecond)

		ch <- playbackStartedGood("d3")
		ch <- playbackFinishedGood("d3")

		ch2 <- channelDtmf("4")

		ch <- playbackStartedGood("d4")
		ch <- playbackFinishedGood("d4")

		<-time.After(50 * time.Millisecond)
		ch2 <- channelDtmf("#")

		ch <- playbackStartedGood("d5")
		ch <- playbackFinishedGood("d5")

	}()

	var opts Options
	opts.MatchFunc = MatchHash
	opts.EchoData = true

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptMatchHash(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("#")
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptMatchAny(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("#")
	}()

	var opts Options
	opts.MatchFunc = MatchAny

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2" {
		t.Errorf("Expected Data to be '2', got, got '%v'", res.Data)
	}
}

func TestPromptMatchLenFunc(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("#")
	}()

	var opts Options
	opts.MatchFunc = MatchLenFunc(3)

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "231" {
		t.Errorf("Expected Data to be '231', got, got '%v'", res.Data)
	}
}

func TestPromptMatchLenOrTerminatorFuncTerm(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("9")
	}()

	var opts Options
	opts.MatchFunc = MatchLenOrTerminatorFunc(8, "9")

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptMatchLenOrTerminatorFunc(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("9")
	}()

	var opts Options
	opts.MatchFunc = MatchLenOrTerminatorFunc(3, "9")

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "231" {
		t.Errorf("Expected Data to be '231', got, got '%v'", res.Data)
	}
}

func TestPromptMatchTerminatorFunc(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("9")
	}()

	var opts Options
	opts.MatchFunc = MatchTerminatorFunc("9")

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptMatchHashPrePromptComplete(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	dtmfA := mock.NewMockSubscription(ctrl)
	dtmfB := mock.NewMockSubscription(ctrl)
	hangup := mock.NewMockSubscription(ctrl)

	sub.EXPECT().Cancel().Times(1)
	dtmfA.EXPECT().Cancel().Times(1)
	dtmfB.EXPECT().Cancel().Times(1)
	hangup.EXPECT().Cancel().Times(1)

	gomock.InOrder(
		bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(hangup),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfA),
		bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Times(1).Return(dtmfB),
		bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).Times(1).Return(sub),
	)

	dtmfAChan := make(chan ari.Event)
	dtmfBChan := make(chan ari.Event)
	ch := make(chan ari.Event)

	dtmfA.EXPECT().Events().MinTimes(1).Return(dtmfAChan) // called in prompt
	dtmfB.EXPECT().Events().MinTimes(1).Return(dtmfBChan) // called in play
	hangup.EXPECT().Events().MinTimes(1)
	sub.EXPECT().Events().MinTimes(1).Return(ch)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")

		dtmfBChan <- channelDtmf("2")
		dtmfAChan <- channelDtmf("2")

		<-time.After(1 * time.Millisecond)

		dtmfAChan <- channelDtmf("3")
		<-time.After(1 * time.Millisecond)

		dtmfAChan <- channelDtmf("1")
		<-time.After(1 * time.Millisecond)

		dtmfAChan <- channelDtmf("4")
		<-time.After(1 * time.Millisecond)

		dtmfAChan <- channelDtmf("#")
		<-time.After(1 * time.Millisecond)

		ch <- playbackFinishedGood("pb1")

		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	<-time.After(1 * time.Millisecond)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptPostPromptHangup(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")
		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		ch3 <- &ari.ChannelHangupRequest{EventData: ari.EventData{Message: ari.Message{Type: "ChannelHangupRequest"}}}
	}()

	res, err := Prompt(ctx, bus, player, nil, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Hangup {
		t.Errorf("Expected Hangup result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}
}

func TestPromptNoSound(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()

	go func() {
		ch2 <- channelDtmf("2")
		ch2 <- channelDtmf("3")
		ch2 <- channelDtmf("1")
		ch2 <- channelDtmf("4")
		ch2 <- channelDtmf("#")
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Complete {
		t.Errorf("Expected Complete result, got '%v'", res)
	}

	if res.Data != "2314" {
		t.Errorf("Expected Data to be '2314', got, got '%v'", res.Data)
	}
}

func TestPromptInterDigitTimeout(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {

		// complete prompt

		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")

		// send initial DTMF
		ch2 <- channelDtmf("2")

		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		// inter-digit timeout should trigger here
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Timeout {
		t.Errorf("Expected Failed result, got '%v'", res.Status)
	}

	if res.Data != "2" {
		t.Errorf("Expected Data to be '2', got, got '%v'", res.Data)
	}
}

func TestPromptInterDigitTimeout210(t *testing.T) {
	for i := 0; i != 10; i++ {
		TestPromptInterDigitTimeout2(t)
	}
}

func TestPromptInterDigitTimeout2(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {

		<-time.After(40 * time.Millisecond)

		// send initial DTMF
		ch2 <- channelDtmf("2")

		// complete prompt

		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")

		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		// inter-digit timeout should trigger here
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Timeout {
		t.Errorf("Expected Timeout result, got '%v'", res.Status)
	}

	if res.Data != "2" {
		t.Errorf("Expected Data to be '2', got, got '%v'", res.Data)
	}
}

func TestPromptOverrallTimeout(t *testing.T) {
	DefaultOverallTimeout = 3 * time.Second
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		// complete prompt

		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")

		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		// overall timeout should trigger
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if res.Status != Timeout {
		t.Errorf("Expected Timeout result, got '%v'", res.Status)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be '', got, got '%v'", res.Data)
	}
}

func TestPromptCancelAfterPlaybackFinished(t *testing.T) {
	DefaultOverallTimeout = 3 * time.Second
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	sub := mock.NewMockSubscription(ctrl)
	ch := make(chan ari.Event)
	sub.EXPECT().Events().MinTimes(1).Return(ch)
	sub.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.PlaybackStarted, ari.Events.PlaybackFinished).MinTimes(1).Return(sub)

	sub2 := mock.NewMockSubscription(ctrl)
	ch2 := make(chan ari.Event)
	sub2.EXPECT().Events().MinTimes(1).Return(ch2)
	sub2.EXPECT().Cancel().MinTimes(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).MinTimes(1).Return(sub2)

	sub3 := mock.NewMockSubscription(ctrl)
	ch3 := make(chan ari.Event)
	sub3.EXPECT().Events().MinTimes(1).Return(ch3)
	sub3.EXPECT().Cancel().Times(1)

	bus.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(1).Return(sub3)

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		// complete prompt

		<-player.Next // play request
		ch <- playbackStartedGood("pb1")
		ch <- playbackFinishedGood("pb1")

		<-player.Next
		ch <- playbackStartedGood("pb2")
		ch <- playbackFinishedGood("pb2")

		<-time.After(1 * time.Second)

		cancel()
	}()

	var opts Options
	opts.MatchFunc = MatchHash

	res, err := Prompt(ctx, bus, player, &opts, "sound:1", "sound:2")

	if err == nil && err.Error() != "context canceled" {
		t.Errorf("Expected error context canceled, got '%v'", err)
	}

	if res.Status != Canceled {
		t.Errorf("Expected Canceled result, got '%v'", res.Status)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be '', got, got '%v'", res.Data)
	}
}
*/
type testPlayback struct {
	id       string
	failData bool
}

func (p *testPlayback) Get(id string) *ari.PlaybackHandle {
	panic("not implemented")
}

func (p *testPlayback) Data(id string) (pd ari.PlaybackData, err error) {
	if p.failData {
		//err = errors.New("Dummy error getting playback data")
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
