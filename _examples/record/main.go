package main

import (
	"context"

	"github.com/inconshreveable/log15"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/CyCoreSystems/ari/ext/record"
)

var log = log15.New()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	native.Logger = log
	record.Logger = log

	// connect
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

	res, err := record.Record(ctx, h,
		record.TerminateOn("any"),
		record.IfExists("overwrite"),
	).Result()
	if err != nil {
		log.Error("failed to record", "error", err)
		return
	}

	if err = res.Save("test-recording"); err != nil {
		log.Error("failed to save recording", "error", err)
	}

	log.Info("completed recording")

	h.Hangup()
	return
}
