package nats

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 20 * time.Millisecond

// New creates a new ari.Client connected to a NATS gateway ARI server
func New(url string) (*ari.Client, error) {

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return &ari.Client{
		Application: &natsApplication{c},
	}, nil
}
