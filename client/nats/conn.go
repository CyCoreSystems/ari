package nats

import "github.com/nats-io/nats"

// Conn is a connection to a native ARI server
type Conn struct {
	c *nats.EncodedConn
}

func newConn(c *nats.EncodedConn) *Conn {
	return &Conn{
		c: c,
	}
}

// Close closes the ARI client
func (c *Conn) Close() {
	c.c.Close()
}
