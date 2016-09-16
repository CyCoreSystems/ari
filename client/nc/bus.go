package nc

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"
	"github.com/nats-io/nats"
)

type natsBus struct {
	conn *Conn
}

func (b *natsBus) Send(msg *v2.Message) {
	panic("Send unsupported")
}

func (b *natsBus) Subscribe(nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan v2.Eventer, 10)
	ns.closeChan = make(chan struct{})

	go func() {
		for _, n := range nx {
			subj := fmt.Sprintf("ari.events.%s", n)
			sub, err := b.conn.conn.Subscribe(subj, func(msg *nats.Msg) {
				eventType := msg.Subject[len("ari.events."):]

				var ariMessage v2.Message
				ariMessage.SetRaw(&msg.Data)
				ariMessage.Type = eventType

				evt := v2.Parse(&ariMessage)
				ns.events <- evt
			})
			if err != nil {
				//TODO: handle error
				panic(err)
			}
			defer sub.Unsubscribe()
		}

		<-ns.closeChan
		ns.closeChan = nil
	}()

	return &ns
}

type natsSubscription struct {
	closeChan chan struct{}
	events    chan v2.Eventer
}

func (ns *natsSubscription) Events() chan v2.Eventer {
	return ns.events
}

func (ns *natsSubscription) Cancel() {
	if ns.closeChan != nil {
		close(ns.closeChan)
	}
}
