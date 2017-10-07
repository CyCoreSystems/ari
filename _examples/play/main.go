package main

import (
	"golang.org/x/net/context"

	"github.com/inconshreveable/log15"

	"github.com/AVOXI/ari"
	"github.com/AVOXI/ari/client/native"
	"github.com/AVOXI/ari/ext/play"
)

var log = log15.New()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	native.Logger = log

	log.Info("Connecting to ARI")
	cl, err := native.Connect(&native.Options{
		Application:  "test",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Error("Failed to build ARI client", "error", err)
		return
	}

	// setup app

	log.Info("Listening for new calls")
	sub := cl.Bus().Subscribe(nil, "StasisStart")

	for {
		select {
		case e := <-sub.Events():
			v := e.(*ari.StasisStart)
			log.Info("Got stasis start", "channel", v.Channel.ID)
			go app(ctx, cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID)))
		case <-ctx.Done():
			return
		}
	}
}

func app(ctx context.Context, h *ari.ChannelHandle) {
	defer h.Hangup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Info("Running app", "channel", h.ID())

	end := h.Subscribe(ari.Events.StasisEnd)
	defer end.Cancel()

	// End the app when the channel goes away
	go func() {
		<-end.Events()
		cancel()
	}()

	if err := h.Answer(); err != nil {
		log.Error("failed to answer call", "error", err)
		return
	}

	if err := play.Play(ctx, h, play.URI("sound:tt-monkeys")).Err(); err != nil {
		log.Error("failed to play sound", "error", err)
		return
	}

	log.Info("completed playback")
	return

}
