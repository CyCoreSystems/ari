package prompt

import (
	"errors"
	"testing"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	v2 "github.com/CyCoreSystems/ari/v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext/audio"
	"github.com/CyCoreSystems/ari/internal/testutils"
	"golang.org/x/net/context"
)

func TestPromptPlayError(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1", failData: true}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))
	}()

	res, err := Prompt(ctx, bus, player, nil, "sound:1", "sound:2")

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
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
	}()

	cancel()

	res, err := Prompt(ctx, bus, player, nil, "sound:1", "sound:2")

	if err == nil || err.Error() != "context canceled" {
		t.Errorf("Expected error 'context cancelled', got '%v'", err)
	}

	if res.Status != Canceled {
		t.Errorf("Expected Failed result, got '%v'", res)
	}

	if res.Data != "" {
		t.Errorf("Expected Data to be empty, got, got '%v'", res.Data)
	}
}

func TestPromptNoInput(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))
	}()

	res, err := Prompt(ctx, bus, player, nil, "sound:1", "sound:2")

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

func TestPromptHangup(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(&v2.ChannelHangupRequest{Event: v2.Event{Message: v2.Message{Type: "ChannelHangupRequest"}}})
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

func TestPromptMatchHashEchoData(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

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
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		<-time.After(50 * time.Millisecond)

		bus.Send(playbackStartedGood("d1"))
		bus.Send(playbackFinishedGood("d1"))

		bus.Send(channelDtmf("3"))
		<-time.After(50 * time.Millisecond)

		bus.Send(playbackStartedGood("d2"))
		bus.Send(playbackFinishedGood("d2"))

		bus.Send(channelDtmf("1"))
		<-time.After(50 * time.Millisecond)

		bus.Send(playbackStartedGood("d3"))
		bus.Send(playbackFinishedGood("d3"))

		bus.Send(channelDtmf("4"))

		bus.Send(playbackStartedGood("d4"))
		bus.Send(playbackFinishedGood("d4"))

		<-time.After(50 * time.Millisecond)
		bus.Send(channelDtmf("#"))

		bus.Send(playbackStartedGood("d5"))
		bus.Send(playbackFinishedGood("d5"))

	}()

	var opts Options
	opts.MatchFunc = MatchHash
	opts.EchoData = true

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

func TestPromptMatchHash(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("#"))
	}()

	var opts Options
	opts.MatchFunc = MatchHash

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

func TestPromptMatchAny(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("#"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("#"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("9"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("9"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("9"))
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

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))

		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("#"))

		bus.Send(playbackFinishedGood("pb1"))

		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

	}()

	var opts Options
	opts.MatchFunc = MatchHash

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

func TestPromptPostPromptHangup(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))
		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

		<-time.After(1 * time.Second)

		bus.Send(&v2.ChannelHangupRequest{Event: v2.Event{Message: v2.Message{Type: "ChannelHangupRequest"}}})
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

func TestPromptNoSound100(t *testing.T) {
	Logger.SetHandler(log15.DiscardHandler())
	for i := 0; i != 100; i++ {
		TestPromptNoSound(t)
	}
}

func TestPromptNoSound(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()

	go func() {
		bus.Send(channelDtmf("2"))
		bus.Send(channelDtmf("3"))
		bus.Send(channelDtmf("1"))
		bus.Send(channelDtmf("4"))
		bus.Send(channelDtmf("#"))
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
	audio.Logger = log15.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {

		// complete prompt

		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))

		// send initial DTMF
		bus.Send(channelDtmf("2"))

		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

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

func TestPromptInterDigitTimeout2(t *testing.T) {
	audio.MaxPlaybackTime = 3 * time.Second
	audio.Logger = log15.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {

		// send initial DTMF
		bus.Send(channelDtmf("2"))

		// complete prompt

		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))

		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

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
	audio.Logger = log15.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		// complete prompt

		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))

		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

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
	audio.Logger = log15.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := testutils.NewBus()

	player := testutils.NewPlayer()
	player.Append(ari.NewPlaybackHandle("pb1", &testPlayback{id: "pb1"}), nil)
	player.Append(ari.NewPlaybackHandle("pb2", &testPlayback{id: "pb2"}), nil)

	go func() {
		// complete prompt

		<-player.Next // play request
		bus.Send(playbackStartedGood("pb1"))
		bus.Send(playbackFinishedGood("pb1"))

		<-player.Next
		bus.Send(playbackStartedGood("pb2"))
		bus.Send(playbackFinishedGood("pb2"))

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

var channelDtmf = func(dtmf string) v2.Eventer {
	return &v2.ChannelDtmfReceived{
		Event: v2.Event{
			Message: v2.Message{
				Type: "ChannelDtmfReceived",
			},
		},
		Digit: dtmf,
	}
}

var playbackStartedGood = func(id string) v2.Eventer {
	return &v2.PlaybackStarted{
		Event: v2.Event{
			Message: v2.Message{
				Type: "PlaybackStarted",
			},
		},
		Playback: v2.Playback{
			ID: id,
		},
	}
}

var playbackFinishedGood = func(id string) v2.Eventer {
	return &v2.PlaybackFinished{
		Event: v2.Event{
			Message: v2.Message{
				Type: "PlaybackFinished",
			},
		},
		Playback: v2.Playback{
			ID: id,
		},
	}
}

var playbackStartedBadMessageType = &v2.PlaybackStarted{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackStarted2",
		},
	},
	Playback: v2.Playback{
		ID: "pb1",
	},
}

var playbackFinishedBadMessageType = &v2.PlaybackFinished{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackFinished2",
		},
	},
	Playback: v2.Playback{
		ID: "pb1",
	},
}

var playbackStartedDifferentPlaybackID = &v2.PlaybackStarted{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackStarted",
		},
	},
	Playback: v2.Playback{
		ID: "pb2",
	},
}

var playbackFinishedDifferentPlaybackID = &v2.PlaybackFinished{
	Event: v2.Event{
		Message: v2.Message{
			Type: "PlaybackFinished",
		},
	},
	Playback: v2.Playback{
		ID: "pb2",
	},
}
