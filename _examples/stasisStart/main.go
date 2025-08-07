package main

import (
	"context"
	"net/http"
	"os"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/client/native"
)

var log = slog.New(slog.NewTextHandler(os.Stderr, nil))

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Connecting to ARI")

	cl, err := native.Connect(&native.Options{
		Application:  "test",
		Logger:       log.With("app", "test"),
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

	log.Info("Starting listener app")

	go listenApp(ctx, cl, channelHandler)

	// start call start listener

	log.Info("Starting HTTP Handler")

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make call
		log.Info("Make sample call")

		h, err := createCall(cl)
		if err != nil {
			log.Error("Failed to create call", "error", err)

			w.WriteHeader(http.StatusBadGateway)

			w.Write([]byte("Failed to create call: " + err.Error())) //nolint:errcheck

			return
		}

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(h.ID())) //nolint:errcheck
	}))

	log.Info("Listening for requests on port 9990")

	http.ListenAndServe(":9990", nil) //nolint:errcheck
}

func listenApp(ctx context.Context, cl ari.Client, handler func(cl ari.Client, h *ari.ChannelHandle)) {
	sub := cl.Bus().Subscribe(nil, "StasisStart")
	end := cl.Bus().Subscribe(nil, "StasisEnd")

	for {
		select {
		case e := <-sub.Events():
			v := e.(*ari.StasisStart)

			log.Info("Got stasis start", "channel", v.Channel.ID)

			go handler(cl, cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID)))
		case <-end.Events():
			log.Info("Got stasis end")
		case <-ctx.Done():
			return
		}
	}
}

func createCall(cl ari.Client) (h *ari.ChannelHandle, err error) {
	h, err = cl.Channel().Create(nil, ari.ChannelCreateRequest{
		Endpoint: "Local/1000",
		App:      "example",
	})

	return
}

func channelHandler(cl ari.Client, h *ari.ChannelHandle) {
	log.Info("Running channel handler")

	stateChange := h.Subscribe(ari.Events.ChannelStateChange)
	defer stateChange.Cancel()

	data, err := h.Data()
	if err != nil {
		log.Error("Error getting data", "error", err)
		return
	}

	log.Info("Channel State", "state", data.State)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		log.Info("Waiting for channel events")

		defer wg.Done()

		defer stateChange.Cancel()

		for ev := range stateChange.Events() {
			if ev == nil {
				return
			}

			log.Info("Got state change request")

			data, err = h.Data()
			if err != nil {
				log.Error("Error getting data", "error", err)
				continue
			}

			log.Info("New Channel State", "state", data.State)

			if data.State == "Up" {
				return
			}
		}
	}()

	h.Answer() //nolint:errcheck

	wg.Wait()

	h.Hangup() //nolint:errcheck
}
