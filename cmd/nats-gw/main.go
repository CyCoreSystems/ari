package main

import (
	"strings"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/nats-io/nats"
)

func main() {

	opts := native.Options{
		Application: "my-app",
		URL:         "...",
	}

	cl, err := native.New(&opts)
	if err != nil {
		panic(err)
	}

	//TODO: forward ARI events to NATS

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
}

func server(cl *ari.Client, e *nats.EncodedConn) {

	e.Subscribe("applications.get", func(_ string, reply string, _ string) {
		cl, err := cl.Application.List()
		for _, c := range cl {

		}
	})

	e.Subscribe("applications.get.>", func(subj string, reply string, _ string) {
		id := strings.Join(strings.Split(subj, ".")[3:], ".")
		e.Publish(reply, cl.Application.Data(id))
	})

	e.Subscribe("get./applications/.>", func(subj string, reply string, _ string) {
		id := strings.Join(strings.Split(subj, ".")[3:], ".")
		e.Publish(reply, cl.Application.Data(id))
	})

}
