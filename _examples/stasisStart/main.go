package main

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
)

var log = log15.New()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect

	cl, err := native.Connect(&native.Options{
		Application:  "example",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://asterisk:8088/ari",
		WebsocketURL: "ws://asterisk:8088/ari/events",
	})
	if err != nil {
		log.Error("Failed to build nc ARI client", "error", err)
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
			w.Write([]byte("Failed to create call: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(h.ID()))
	}))

	log.Info("Listening for requests on port 9990")
	http.ListenAndServe(":9990", nil)

	return
}

func listenApp(ctx context.Context, cl ari.Client, handler func(cl ari.Client, h *ari.ChannelHandle)) {
	sub := cl.Bus().Subscribe(nil, "StasisStart")
	end := cl.Bus().Subscribe(nil, "StasisEnd")

	for {
		select {
		case e := <-sub.Events():
			log.Info("Got stasis start")
			v := e.(*ari.StasisStart)
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

func connect(ctx context.Context) (cl ari.Client, err error) {

	log.Info("Connecting")

	cl, err = native.Connect(&native.Options{})
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

		for {
			select {
			case <-stateChange.Events():
				log.Info("Got state change request")

				data, err = h.Data()
				if err != nil {
					log.Error("Error getting data", "error", err)
					continue
				}
				log.Info("New Channel State", "state", data.State)

				if data.State == "Up" {
					stateChange.Cancel() // stop subscription to state change events
					return
				}
			}
		}

	}()

	h.Answer()

	wg.Wait()

	h.Hangup()
}
