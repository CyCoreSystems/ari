package main

import (
	"github.com/inconshreveable/log15"

	"github.com/CyCoreSystems/ari/v5/client/native"
)

func main() {
	// OPTIONAL: setup logging
	log := log15.New()
	native.Logger = log

	log.Info("connecting")

	appA, err := native.Connect(&native.Options{
		Application:  "example",
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
