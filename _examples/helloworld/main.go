package main

import (
	"github.com/inconshreveable/log15"

	"github.com/CyCoreSystems/ari/client/native"
)

func main() {
	// OPTIONAL: setup logging
	log := log15.New()
	native.Logger = log

	log.Info("Connecting")

	cl, err := native.Connect(&native.Options{
		Application:  "example",
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
