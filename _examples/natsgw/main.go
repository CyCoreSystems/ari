package main

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/CyCoreSystems/ari/server/natsgw"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

func main() {

	<-time.After(5 * time.Second)

	native.Logger = log15.New()

	cl, err := connect()
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	opts := natsgw.Options{
		URL: "nats://nats:4222",
	}

	srv, err := natsgw.NewServer(cl, &opts)
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	srv.Listen()
}

func connect() (cl *ari.Client, err error) {

	opts := native.Options{
		Application:  "example",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://asterisk:8088/ari",
		WebsocketURL: "ws://asterisk:8088/ari/events",
	}

	cl, err = native.New(opts)
	return
}
