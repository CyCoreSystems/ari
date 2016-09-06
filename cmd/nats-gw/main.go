package main

import (
	"strings"

	"github.com/CyCoreSystems/ari/client/native"
	"github.com/nats-io/nats"
	"golang.org/x/net/context"
)

func main() {

	ctx := context.Background()

	cl, _ := native.New(nil)
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer c.Close()

	c.Subscribe("ari.applications", func(_ string, reply string, _ string) {
		apps, err := cl.Application.List()
		if err != nil {
			//TODO: what do?
			return
		}

		var items []string
		for _, app := range apps {
			items = append(items, app.ID())
		}

		c.Publish(reply, items)
	})

	c.Subscribe("ari.applications.data.>", func(subj string, reply string, _ string) {

		appName := strings.Join(strings.Split(subj, ".")[3:], ".")

		d, err := cl.Application.Data(appName)
		if err != nil {
			//TODO: what do?
			return
		}

		c.Publish(reply, d)
	})

	c.Subscribe("ari.applications.subscribe.>", func(subj string, reply string, source string) {

		appName := strings.Join(strings.Split(subj, ".")[3:], ".")

		err := cl.Application.Subscribe(appName, source)
		if err != nil {
			c.Publish(reply, err.Error())
			return
		}

		c.Publish(reply, "OK")
	})

	c.Subscribe("ari.applications.unsubscribe.>", func(subj string, reply string, source string) {

		appName := strings.Join(strings.Split(subj, ".")[3:], ".")

		err := cl.Application.Unsubscribe(appName, source)
		if err != nil {
			c.Publish(reply, err.Error())
			return
		}

		c.Publish(reply, "OK")
	})

	<-ctx.Done()

}
