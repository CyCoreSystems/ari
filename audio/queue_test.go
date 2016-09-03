package audio

import (
	"strings"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"

	"golang.org/x/net/context"
)

type multiDummyPlayer struct {
	players []dummyPlayer
	C       chan struct{}
}

func (mdp *multiDummyPlayer) Play(mediaURI string) (h *ari.PlaybackHandle, err error) {
	h = mdp.players[0].H
	err = mdp.players[0].Err
	mdp.players = mdp.players[1:]
	mdp.C <- struct{}{}
	return
}

type multiDummySubscriber struct {
	subs map[string]*v2.Subscription
}

func (mdm *multiDummySubscriber) Subscribe(n ...string) *v2.Subscription {
	a := v2.NewSubscription("")
	mdm.subs[strings.Join(n, ";")] = a
	return a
}

func TestQueueTimeout(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	err := q.Play(ctx, player, nil)

	if te, ok := err.(timeoutErrI); !ok || !te.IsTimeout() {
		t.Errorf("Expected timeout error, got: '%v'", err)
	} else if err.Error() != "Timeout waiting for start of playback" {
		t.Errorf("Expected timeout waiting for start of playback error, got: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueTimeoutSecond(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.C // wait for play request
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
	}()

	err := q.Play(ctx, player, nil)

	if te, ok := err.(timeoutErrI); !ok || !te.IsTimeout() {
		t.Errorf("Expected timeout error, got: '%v'", err)
	} else if err.Error() != "Timeout waiting for stop of playback" {
		t.Errorf("Expected timeout waiting for stop of playback error, got: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}
}

func TestQueueTimeoutThird(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.C // wait for play request
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for play request
	}()

	err := q.Play(ctx, player, nil)

	if te, ok := err.(timeoutErrI); !ok || !te.IsTimeout() {
		t.Errorf("Expected timeout error, got: '%v'", err)
	} else if err.Error() != "Timeout waiting for start of playback" {
		t.Errorf("Expected timeout waiting for start of playback error, got: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueTimeoutFourth(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")

	}()

	err := q.Play(ctx, player, nil)

	if te, ok := err.(timeoutErrI); !ok || !te.IsTimeout() {
		t.Errorf("Expected timeout error, got: '%v'", err)
	} else if err.Error() != "Timeout waiting for stop of playback" {
		t.Errorf("Expected timeout waiting for stop of playback error, got: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueSuccess(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb2")

	}()

	err := q.Play(ctx, player, nil)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueSuccessWithEmpty(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("") // empty should just be skipped
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb2")

	}()

	err := q.Play(ctx, player, nil)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueExitOnDTMF(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	opts := &Options{
		ExitOnDTMF: "3",
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("2")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("3")
		//bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb2")
	}()

	err := q.Play(ctx, player, opts)

	if err == nil || err.Error() != "context canceled" {
		t.Errorf("Expected error 'context canceled', got '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "23" {
		t.Errorf("Expected DTMF '23' during playback, got '%s'", dtmf)
	}

}

func TestQueueExitOnDTMF100(t *testing.T) {
	for i := 0; i != 100; i++ {
		TestQueueExitOnDTMF(t)
	}
}

func TestQueueDoneTrigger(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	done := make(chan struct{})

	opts := &Options{
		Done: done,
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("2")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("3")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb2")
	}()

	var err error
	go func() {
		err = q.Play(ctx, player, opts)
	}()

	select {
	case <-done:
	case <-time.After(MaxPlaybackTime * 2): // 2 because two audio clips
		t.Errorf("options.Done never got triggered")
	}

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "23" {
		t.Errorf("Expected DTMF '23' during playback, got '%s'", dtmf)
	}
}

func TestQueueDTMF(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	dtmfChan := make(chan *v2.ChannelDtmfReceived, 2)

	opts := &Options{
		DTMF: dtmfChan,
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("2")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		<-player.C // wait for second play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb2")
		bus.subs["ChannelDtmfReceived"].C <- channelDtmf("3")
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb2")
	}()

	err := q.Play(ctx, player, opts)
	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "23" {
		t.Errorf("Expected DTMF '23' during playback, got '%s'", dtmf)
	}

	select {
	case e := <-dtmfChan:
		if e.Digit != "2" {
			t.Errorf("Expected first DTMF digit to be '2', was '%s'", e.Digit)
		}
	default:
		t.Errorf("Unexpected fallthrough checking opts.DTMF output")
	}

	select {
	case e := <-dtmfChan:
		if e.Digit != "3" {
			t.Errorf("Expected first DTMF digit to be '3', was '%s'", e.Digit)
		}
	default:
		t.Errorf("Unexpected fallthrough checking opts.DTMF output")
	}

	select {
	case e := <-dtmfChan:
		t.Errorf("Unexpected third item in dtmfChan: '%v'", e)
	default:
	}

}

func TestQueueFlush(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.C // wait for first play request

		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		q.Flush()
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackFinishedGood("pb1")

		select {
		case <-player.C: // wait for second play request
			t.Errorf("Unexpected second play after flush")
		case <-time.After(1 * time.Second):
		}

	}()

	err := q.Play(ctx, player, nil)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}

func TestQueueCancel(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &multiDummySubscriber{
		subs: make(map[string]*v2.Subscription),
	}

	player := &multiDummyPlayer{
		C: make(chan struct{}, 10),
		players: []dummyPlayer{
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}),
				Err: nil,
			},
			dummyPlayer{
				H:   ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}),
				Err: nil,
			},
		},
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.C // wait for first play request
		bus.subs["PlaybackStarted;PlaybackFinished"].C <- playbackStartedGood("pb1")
		cancel()
	}()

	err := q.Play(ctx, player, nil)

	if err == nil || err.Error() != "context canceled" {
		t.Errorf("Expected error 'context canceled', got '%v'", err)
	}

	dtmf := q.ReceivedDTMF()
	if dtmf != "" {
		t.Errorf("Unexpected DTMF during playback: '%s'", dtmf)
	}

}
