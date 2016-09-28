package audio

import (
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/golang/mock/gomock"
	"golang.org/x/net/context"
)

func TestQueueSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().MinTimes(1)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	var st Status
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1")
		st, err = q.Play(ctx, player, &Options{
			Done: doneCh,
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	playbackFinishedChan <- &ari.PlaybackFinished{Playback: ari.PlaybackData{ID: "pb1"}}

	<-doneCh

	if st != Finished {
		t.Errorf("Expected status to be Finished, was '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	cancel()
	<-time.After(1 * time.Millisecond)
}

func TestQueueMultiples(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)
	handle2 := ari.NewPlaybackHandle("pb2", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel().Times(2)
	playbackStartedSub.EXPECT().Events().Times(2).Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel().Times(2)
	playbackFinishedSub.EXPECT().Events().Times(2).Return(playbackFinishedChan)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Times(2).Return(hangupSub)
	hangupSub.EXPECT().Cancel().Times(2)
	hangupSub.EXPECT().Events().Times(4)

	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)
	player.EXPECT().Play(gomock.Any(), "sound:2").Return(handle2, nil)

	var st Status = 100
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1", "", "sound:2")
		st, err = q.Play(ctx, player, &Options{
			Done: doneCh,
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	playbackFinishedChan <- &ari.PlaybackFinished{Playback: ari.PlaybackData{ID: "pb1"}}

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb2"}}

	playbackFinishedChan <- &ari.PlaybackFinished{Playback: ari.PlaybackData{ID: "pb2"}}

	<-doneCh

	if st != Finished {
		t.Errorf("Expected status to be Finished, was '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	cancel()

}

func TestQueueCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)
	//handle2 := ari.NewPlaybackHandle("pb2", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	//player.EXPECT().Subscribe("pb2", ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel().Times(1)
	playbackStartedSub.EXPECT().Events().Times(1).Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	//player.EXPECT().Subscribe("pb2", ari.Events.PlaybackFinished).Return(playbackFinishedSub)

	playbackFinishedSub.EXPECT().Cancel().Times(1)
	playbackFinishedSub.EXPECT().Events().Times(1).Return(playbackFinishedChan)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).MinTimes(1).Return(hangupSub)
	hangupSub.EXPECT().Cancel().MinTimes(1)
	hangupSub.EXPECT().Events().MinTimes(1)

	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)
	//	player.EXPECT().Play(gomock.Any(), "sound:2").Return(handle2, nil)

	var st Status = 1000
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1", "", "sound:2")
		st, err = q.Play(ctx, player, &Options{
			Done: doneCh,
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	cancel()

	<-doneCh

	if st != Canceled {
		t.Errorf("Expected status to be Finished, was '%v'", st)
	}

	cancel()

}

func TestQueueSimpleDTMF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	dtmfChan := make(chan ari.Event)
	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1).Return(dtmfChan)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().MinTimes(1)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	var st Status
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1")
		st, err = q.Play(ctx, player, &Options{
			Done: doneCh,
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "1",
	}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "3",
	}

	playbackFinishedChan <- &ari.PlaybackFinished{Playback: ari.PlaybackData{ID: "pb1"}}

	<-doneCh

	if st != Finished {
		t.Errorf("Expected status to be Finished, was '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	if q.ReceivedDTMF() != "13" {
		t.Errorf("Expected DTMF to be '13', was '%v'", q.ReceivedDTMF())
	}

	cancel()
	<-time.After(1 * time.Millisecond)
}

func TestQueueDTMFExit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	dtmfChan := make(chan ari.Event)
	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1).Return(dtmfChan)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().MinTimes(1)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	var st Status
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1")
		st, err = q.Play(ctx, player, &Options{
			Done:       doneCh,
			ExitOnDTMF: "4",
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "1",
	}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "3",
	}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "2",
	}

	dtmfChan <- &ari.ChannelDtmfReceived{
		Digit: "4",
	}

	<-doneCh

	if st != Finished {
		t.Errorf("Expected status to be Finished, was '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	if q.ReceivedDTMF() != "1324" {
		t.Errorf("Expected DTMF to be '1324', was '%v'", q.ReceivedDTMF())
	}

	cancel()
	<-time.After(1 * time.Millisecond)
}

func TestQueueHangup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playback := mock.NewMockPlayback(ctrl)
	player := mock.NewMockPlayer(ctrl)

	handle := ari.NewPlaybackHandle("pb1", playback)

	playbackStartedChan := make(chan ari.Event)
	playbackStartedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackStarted).Return(playbackStartedSub)
	playbackStartedSub.EXPECT().Cancel()
	playbackStartedSub.EXPECT().Events().Return(playbackStartedChan)

	playbackFinishedChan := make(chan ari.Event)
	playbackFinishedSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.PlaybackFinished).Return(playbackFinishedSub)
	playbackFinishedSub.EXPECT().Cancel()
	playbackFinishedSub.EXPECT().Events().Return(playbackFinishedChan)

	dtmfChan := make(chan ari.Event)
	dtmfSub := mock.NewMockSubscription(ctrl)
	dtmfSub.EXPECT().Events().MinTimes(1).Return(dtmfChan)
	dtmfSub.EXPECT().Cancel()
	player.EXPECT().Subscribe(ari.Events.ChannelDtmfReceived).Return(dtmfSub)

	hangupChan := make(chan ari.Event)
	hangupSub := mock.NewMockSubscription(ctrl)
	player.EXPECT().Subscribe(ari.Events.ChannelHangupRequest, ari.Events.ChannelDestroyed).Return(hangupSub)
	hangupSub.EXPECT().Cancel()
	hangupSub.EXPECT().Events().MinTimes(1).Return(hangupChan)

	player.EXPECT().Play(gomock.Any(), "sound:1").Return(handle, nil)

	var st Status
	var err error

	var q *Queue
	q = NewQueue()

	doneCh := make(chan struct{})

	go func() {
		q.Add("sound:1")
		st, err = q.Play(ctx, player, &Options{
			Done:       doneCh,
			ExitOnDTMF: "4",
		})
	}()

	playbackStartedChan <- &ari.PlaybackStarted{Playback: ari.PlaybackData{ID: "pb1"}}

	hangupChan <- &ari.EventData{}

	<-doneCh

	if st != Hangup {
		t.Errorf("Expected status to be Hangup, was '%v'", st)
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	cancel()
	<-time.After(1 * time.Millisecond)
}
