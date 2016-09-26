package cmd

import (
	"os"
	"os/signal"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/CyCoreSystems/ari/client/nc"
	"github.com/CyCoreSystems/ari/server/natsgw"
	"github.com/spf13/viper"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

func runServer(log log15.Logger) int {

	// inject logger

	native.Logger = log
	nc.Logger = log

	log.Info("Starting ari-natsgw server")

	log.Debug("Connecting to ARI")

	cl, err := connect()
	if err != nil {
		log.Error("Failed to connect to ARI", "error", err)
		return -1
	}
	defer cl.Close()

	opts := natsgw.Options{
		URL:    viper.GetString("nats.url"),
		Logger: log,
	}

	log.Debug("Connecting to NATS")

	srv, err := natsgw.NewServer(cl, &opts)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		return -1
	}
	defer srv.Close()

	srv.Start()

	log.Info("Started ari-natsgw server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c

	log.Debug("Got signal", "signal", s)
	log.Info("Shutting down ari-natsgw server")

	srv.Close()

	return 0
}

func connect() (cl *ari.Client, err error) {

	opts := native.Options{
		Application:  viper.GetString("ari.application"),
		Username:     viper.GetString("ari.username"),
		Password:     viper.GetString("ari.password"),
		URL:          viper.GetString("ari.http_url"),
		WebsocketURL: viper.GetString("ari.websocket_url"),
	}

	cl, err = native.New(opts)
	return
}
