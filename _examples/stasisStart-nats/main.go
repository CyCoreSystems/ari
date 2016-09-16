package main

import (
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/nc"
	v2 "github.com/CyCoreSystems/ari/v2"
)

var log = log15.New()
var wg sync.WaitGroup

func main() {

	<-time.After(20 * time.Second)

	if i := run(); i != 0 {
		os.Exit(i)
	}
}

func channelHandler(cl *ari.Client, h *ari.ChannelHandle) {
	log.Info("Running channel handler")
	defer wg.Done()

	// TODO: this subscription /should/ be channel specific. Not all Channel events are setup this way yet
	sub := h.Subscribe("ChannelHangupRequest")
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Info("Waiting for channel hangup request")
		<-sub.Events()
		log.Info("Got Channel hangup request")
	}()

	data, err := h.Data()
	if err != nil {
		log.Error("Error getting data", "error", err)
		return
	}
	log.Info("Channel Data", "data", data)

	h.Answer()

	h.Hangup()
}

func run() int {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logging

	nc.Logger = log15.New()

	// connect

	cl, err := connect(ctx)
	if err != nil {
		log.Error("Failed to build nc ARI client", "error", err)
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
	case e := <-sub.Events():
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

	opts := nc.Options{
		URL: "nats://nats:4222",
	}

	log.Info("Connecting")

	cl, err = nc.New(opts)
	return
}
