package audio

import (
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/internal/testutils"
	v2 "github.com/CyCoreSystems/ari/v2"

	"golang.org/x/net/context"
)

func TestQueueTimeout(t *testing.T) {
	MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	err := q.Play(ctx, player, nil)

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for start of playback" {
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for play request
		bus.Send(playbackStartedGood("pb1"))
	}()

	err := q.Play(ctx, player, nil)

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for stop of playback" {
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for play request
	}()

	err := q.Play(ctx, player, nil)

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for start of playback" {
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
	}()

	err := q.Play(ctx, player, nil)

	if !isTimeout(err) {
		t.Errorf("Expected timeout error, got: '%v'", err)
	}

	if err != nil && err.Error() != "Timeout waiting for stop of playback" {
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("") // empty should just be skipped
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	opts := &Options{
		ExitOnDTMF: "3",
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(channelDtmf("2"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(channelDtmf("3"))
	}()

	err := q.Play(ctx, player, opts)

	if err != nil {
		t.Errorf("Unexpected error '%v'", err)
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	done := make(chan struct{})

	opts := &Options{
		Done: done,
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(channelDtmf("2"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(channelDtmf("3"))
		bus.Send(playbackFinishedGood("pb2"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	dtmfChan := make(chan *v2.ChannelDtmfReceived, 2)

	opts := &Options{
		DTMF: dtmfChan,
	}

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(channelDtmf("2"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // wait for second play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(channelDtmf("3"))
		bus.Send(playbackFinishedGood("pb2"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {

		<-player.Next // wait for first play request

		bus.Send(playbackStartedGood("pb1"))
		q.Flush()
		bus.Send(playbackFinishedGood("pb1"))
	}()

	err := q.Play(ctx, player, nil)

	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	select {
	case <-player.Next: // wait for second play request
		t.Errorf("Unexpected second play after flush")
	default:
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	q := NewQueue(bus)
	q.Add("sound:1")
	q.Add("sound:2")

	go func() {
		<-player.Next // wait for first play request
		bus.Send(playbackStartedGood("pb1"))
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
