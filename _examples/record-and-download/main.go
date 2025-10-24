package main

import (
	"context"
	"errors"
	"io"
	"os"

	"golang.org/x/exp/slog"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/client/native"
	"github.com/CyCoreSystems/ari/v6/ext/record"
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
		record.Format("wav"),
		record.Name("myrecording"),
	).Result()
	if err != nil {
		log.Error("failed to record", "error", err)
		return
	}

	log.Info("saving recording")

	rec, err := res.Download()
	if err != nil {
		log.Error("failed to download recording", "error", err)
		return
	}
	defer rec.Close() //nolint:errcheck

	// Create output file
	outFile, err := os.Create("audio.wav")
	if err != nil {
		log.Error("failed to create output file", "error", err)
		return
	}
	defer outFile.Close() //nolint:errcheck

	// Write the data to the output file
	buf := make([]byte, 1024)
	for {
		n, err := rec.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Error("failed to read from recording", "error", err)
			}
			break
		}
		if n > 0 {
			if _, err := outFile.Write(buf[:n]); err != nil {
				log.Error("failed to write to output file", "error", err)
				return
			}
		}
	}

	log.Info("completed recording")
}
