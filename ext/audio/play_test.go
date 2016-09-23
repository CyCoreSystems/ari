package audio

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/golang/mock/gomock"

	"golang.org/x/net/context"
)

func TestPlaySync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	go func() {
		playbackStartedChan <- &ari.EventData{}
		playbackFinishedChan <- &ari.EventData{}
	}()

	st, err := Play(ctx, playback, player, "sound:1")

	if st != Finished {
		t.Errorf("Expected playback status to be Finished, got '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected playback error '%v'", err)
	}

}

func TestPlay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	playbackStartedChan <- &ari.EventData{}

	<-pb.Started()

	if pb.Status() != InProgress {
		t.Errorf("Expected playback status to be InProgress, got '%v'", pb.Status())
	}

	playbackFinishedChan <- &ari.EventData{}

	<-pb.Stopped()

	if pb.Status() != Finished {
		t.Errorf("Expected playback status to be Finished, got '%v'", pb.Status())
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected playback error '%v'", pb.Err())
	}

}

func TestPlayHangup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(0).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(1).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	hangupChan <- &ari.EventData{}

	<-pb.Started()

	<-pb.Stopped()

	if pb.Status() != Hangup {
		t.Errorf("Expected playback status to be Hangup, got '%v'", pb.Status())
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected playback error '%v'", pb.Err())
	}

}

func TestPlayHangupAfterStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(1).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	playbackStartedChan <- &ari.EventData{}

	<-pb.Started()

	hangupChan <- &ari.EventData{}

	<-pb.Stopped()

	if pb.Status() != Hangup {
		t.Errorf("Expected playback status to be Hangup, got '%v'", pb.Status())
	}

	if pb.Err() != nil {
		t.Errorf("Unexpected playback error '%v'", pb.Err())
	}

}

func TestPlayContextCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(0).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(1).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	cancel()

	<-pb.Started()
	<-pb.Stopped()

	if pb.Status() != Canceled {
		t.Errorf("Expected playback status to be Canceled, got '%v'", pb.Status())
	}

	if pb.Err() == nil && pb.Err().Error() != "context canceled" {
		t.Errorf("Expected playback error 'context canceled', got '%v'", pb.Err())
	}

}

func TestPlayContextCancelAfterStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(1).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	playbackStartedChan <- &ari.EventData{}

	<-pb.Started()

	cancel()

	<-pb.Stopped()

	if pb.Status() != Canceled {
		t.Errorf("Expected playback status to be Canceled, got '%v'", pb.Status())
	}

	if pb.Err() == nil && pb.Err().Error() != "context canceled" {
		t.Errorf("Expected playback error 'context canceled', got '%v'", pb.Err())
	}

}

func TestPlayCancelCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(1).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	playbackStartedChan <- &ari.EventData{}

	<-pb.Started()

	pb.Cancel()

	<-pb.Stopped()

	if pb.Status() != Canceled {
		t.Errorf("Expected playback status to be Canceled, got '%v'", pb.Status())
	}

	if pb.Err() == nil && pb.Err().Error() != "context canceled" {
		t.Errorf("Expected playback error 'context canceled', got '%v'", pb.Err())
	}

}

func TestPlayTimeoutStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(0).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(1).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	<-pb.Started()
	<-pb.Stopped()

	if pb.Status() != Timeout {
		t.Errorf("Expected playback status to be Timeout, got '%v'", pb.Status())
	}

	cancel()
	<-time.After(1 * time.Millisecond)
}

func TestPlayTimeoutFinished(t *testing.T) {
	MaxPlaybackTime = 300 * time.Millisecond

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(1).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(2).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	pb := PlayAsync(ctx, playback, player, "sound:1")

	playbackStartedChan <- &ari.EventData{}
	<-pb.Started()

	<-pb.Stopped()

	if pb.Status() != Timeout {
		t.Errorf("Expected playback status to be Timeout, got '%v'", pb.Status())
	}

}

func TestPlayFailToPlay(t *testing.T) {
	MaxPlaybackTime = 300 * time.Millisecond

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playback.EXPECT().Get(gomock.Any()).Return(handle)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Times(0).Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	playback.EXPECT().Subscribe("pb1", ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Times(0).Return(playbackFinishedChan)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().Times(0).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(nil, errors.New("Failed to play"))

	pb := PlayAsync(ctx, playback, player, "sound:1")

	<-pb.Started()
	<-pb.Stopped()

	if pb.Status() != Failed {
		t.Errorf("Expected playback status to be Failed, got '%v'", pb.Status())
	}

	if pb.Err() == nil || pb.Err().Error() != "Failed to play" {
		t.Errorf("Expected error 'Failed to play', got '%v'", pb.Err())
	}

}
