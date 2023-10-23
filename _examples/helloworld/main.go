package main

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/PolyAI-LDN/ari/v6/client/native"
)

func main() {
	// OPTIONAL: setup logging
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	log.Info("Connecting")

	cl, err := native.Connect(&native.Options{
		Application:  "example",
		Logger:       log.With("app", "example"),
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Error("Failed to build native ARI client", "error", err)
		return
	}

	defer cl.Close()

	log.Info("Connected")

	info, err := cl.Asterisk().Info(nil)
	if err != nil {
		log.Error("Failed to get Asterisk Info", "error", err)
		return
	}

	log.Info("Asterisk Info", "info", info)
}
