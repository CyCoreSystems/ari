package main

import (
	"os"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
)

func main() {
	if i := run(); i != 0 {
		os.Exit(i)
	}
}

func run() int {

	// setup logging
	native.Logger = log15.New()
	log := log15.New()

	opts := native.Options{
		Application:  "example",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	}

	log.Info("Connecting")

	cl := native.New(&opts)

	log.Info("Connected")

	// info, err := cl.Asterisk().Info("")
	info, err := cl.Asterisk().Info(&ari.Key{Kind: "build"})
	if err != nil {
		log.Error("Failed to get Asterisk Info", "error", err)
		return -1
	}

	log.Info("Asterisk Info", "info", info)

	return 0
}
