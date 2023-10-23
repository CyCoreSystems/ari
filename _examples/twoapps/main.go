package main

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/PolyAI-LDN/ari/v6/client/native"
)

func main() {
	// OPTIONAL: setup logging
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	log.Info("connecting")

	appA, err := native.Connect(&native.Options{
		Application:  "example",
		Logger:       log.With("app", "exampleA"),
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Error("failed to build native ARI client", "error", err)
		return
	}
	defer appA.Close()

	log.Info("Connected A")

	appB, err := native.Connect(&native.Options{
		Application:  "exampleB",
		Logger:       log.With("app", "exampleB"),
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Error("failed to build native ARI client", "error", err)
		return
	}
	defer appB.Close()

	log.Info("Connected B")

	infoA, err := appA.Asterisk().Info(nil)
	if err != nil {
		log.Error("Failed to get Asterisk Info", "error", err)
		return
	}

	log.Info("AppA AsteriskInfo", "info", infoA)

	infoB, err := appB.Asterisk().Info(nil)
	if err != nil {
		log.Error("Failed to get Asterisk Info", "error", err)
		return
	}

	log.Info("AppB AsteriskInfo", "info", infoB)
}
