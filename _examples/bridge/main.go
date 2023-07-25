package main

import (
	"context"
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/rotisserie/eris"

	"github.com/Amtelco-Software/ari/v6"
	"github.com/Amtelco-Software/ari/v6/client/native"
	"github.com/Amtelco-Software/ari/v6/ext/play"
	"github.com/Amtelco-Software/ari/v6/rid"
)

var ariApp = "test"

var log = log15.New()

var bridge *ari.BridgeHandle

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	log.Info("Connecting to ARI")
	cl, err := native.Connect(&native.Options{
		Application:  "test",
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

	log.Info("Listening for new calls")
	sub := cl.Bus().Subscribe(nil, "StasisStart")

	for {
		select {
		case e := <-sub.Events():
			v := e.(*ari.StasisStart)
			log.Info("Got stasis start", "channel", v.Channel.ID)
			go app(ctx, cl, cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID)))
		case <-ctx.Done():
			return
		}
	}
}

func app(ctx context.Context, cl ari.Client, h *ari.ChannelHandle) {
	log.Info("running app", "channel", h.Key().ID)

	if err := h.Answer(); err != nil {
		log.Error("failed to answer call", "error", err)
		// return
	}

	if err := ensureBridge(ctx, cl, h.Key()); err != nil {
		log.Error("failed to manage bridge", "error", err)
		return
	}

	if err := bridge.AddChannel(h.Key().ID); err != nil {
		log.Error("failed to add channel to bridge", "error", err)
		return
	}

	log.Info("channel added to bridge")
	return
}

type bridgeManager struct {
	h *ari.BridgeHandle
}

func ensureBridge(ctx context.Context, cl ari.Client, src *ari.Key) (err error) {
	if bridge != nil {
		log.Debug("Bridge already exists")
		return nil
	}

	key := src.New(ari.BridgeKey, rid.New(rid.Bridge))
	bridge, err = cl.Bridge().Create(key, "mixing", key.ID)
	if err != nil {
		bridge = nil
		return eris.Wrap(err, "failed to create bridge")
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go manageBridge(ctx, bridge, wg)
	wg.Wait()

	return nil
}

func manageBridge(ctx context.Context, h *ari.BridgeHandle, wg *sync.WaitGroup) {
	// Delete the bridge when we exit
	defer h.Delete()

	destroySub := h.Subscribe(ari.Events.BridgeDestroyed)
	defer destroySub.Cancel()

	enterSub := h.Subscribe(ari.Events.ChannelEnteredBridge)
	defer enterSub.Cancel()

	leaveSub := h.Subscribe(ari.Events.ChannelLeftBridge)
	defer leaveSub.Cancel()

	wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-destroySub.Events():
			log.Debug("bridge destroyed")
			return
		case e, ok := <-enterSub.Events():
			if !ok {
				log.Error("channel entered subscription closed")
				return
			}
			v := e.(*ari.ChannelEnteredBridge)
			log.Debug("channel entered bridge", "channel", v.Channel.Name)
			go func() {
				if err := play.Play(ctx, h, play.URI("sound:confbridge-join")).Err(); err != nil {
					log.Error("failed to play join sound", "error", err)
				}
			}()
		case e, ok := <-leaveSub.Events():
			if !ok {
				log.Error("channel left subscription closed")
				return
			}
			v := e.(*ari.ChannelLeftBridge)
			log.Debug("channel left bridge", "channel", v.Channel.Name)
			go func() {
				if err := play.Play(ctx, h, play.URI("sound:confbridge-leave")).Err(); err != nil {
					log.Error("failed to play leave sound", "error", err)
				}
			}()
		}
	}
}
