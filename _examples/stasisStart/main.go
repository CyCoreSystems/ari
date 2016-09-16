package main

import (
	"os"
	"sync"

	"golang.org/x/net/context"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"

	"github.com/CyCoreSystems/ari/client/native"
)

var log = log15.New()
var wg sync.WaitGroup

func main() {
	if i := run(); i != 0 {
		os.Exit(i)
	}
}

func channelHandler(cl *ari.Client, h *ari.ChannelHandle) {
	log.Info("Running channel handler")
	defer wg.Done()

	h.Answer()

	data, err := h.Data()
	if err != nil {
		log.Error("Error getting data", "error", err)
		return
	}

	log.Info("Channel Data", "data", data)

	h.Hangup()
}

func run() int {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logging

	native.Logger = log15.New()

	// connect

	cl, err := connect(ctx)
	if err != nil {
		log.Error("Failed to build native ARI client", "error", err)
		return -1
	}

	// setup app

	log.Info("Starting listener app")

	go listenApp(ctx, cl, channelHandler)

	// make sample call

	wg.Add(1)
	log.Info("Make sample call")

	_, err = createCall(cl)
	if err != nil {
		log.Error("Failed to create call", "error", err)
	}

	wg.Wait()

	return 0
}

func listenApp(ctx context.Context, cl *ari.Client, handler func(cl *ari.Client, h *ari.ChannelHandle)) {
	sub := cl.Bus.Subscribe("StasisStart")

	select {
	case e := <-sub.C:
		log.Info("Got stasis start")
		stasisStartEvent := e.(*v2.StasisStart)
		go handler(cl, cl.Channel.Get(stasisStartEvent.Channel.Id))
	case <-ctx.Done():
		return
	}

}

func createCall(cl *ari.Client) (h *ari.ChannelHandle, err error) {
	h, err = cl.Channel.Create(ari.OriginateRequest{
		Endpoint: "Local/1000",
		App:      "example",
	})

	return
}

func connect(ctx context.Context) (cl *ari.Client, err error) {

	opts := native.Options{
		Application:  "stasis-start",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
		Context:      ctx,
	}

	log.Info("Connecting")

	cl, err = native.New(opts)
	return
}
