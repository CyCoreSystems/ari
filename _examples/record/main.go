package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"

	"github.com/PolyAI-LDN/ari/v6"
	"github.com/PolyAI-LDN/ari/v6/client/native"
	"github.com/PolyAI-LDN/ari/v6/ext/record"
)

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Connecting to ARI")

	cl, err := native.Connect(&native.Options{
		Application:  "test",
		Logger:       log,
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
	defer h.Hangup() //nolint:errcheck

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
		record.WithLogger(log.With("app", "recorder")),
	).Result()
	if err != nil {
		log.Error("failed to record", "error", err)
		return
	}

	if err = res.Save("test-recording"); err != nil {
		log.Error("failed to save recording", "error", err)
	}

	log.Info("completed recording")

	h.Hangup() //nolint:errcheck
}
