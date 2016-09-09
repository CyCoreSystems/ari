package natsgw

import (
	"golang.org/x/net/context"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Options are the group of options for the natsgw server
type Options struct {
	//TODO: include nats options

	URL string

	Logger log15.Logger
	Parent context.Context
}
