package nc

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 20 * time.Millisecond

// New creates a new ari.Client connected to a gateway ARI server via NATS
func New(url string) (cl *ari.Client, err error) {
	var nc *nats.Conn
	nc, err = nats.Connect(url)
	if err != nil {
		return
	}

	var c *nats.EncodedConn
	c, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return
	}
	c.Close() // close for now, since we aren't using it anywhere

	cl = &ari.Client{}

	return
}
