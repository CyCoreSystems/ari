package natsgw

import (
	"errors"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
	"golang.org/x/net/context"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Server is the nats gateway server
type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	upstream *ari.Client
	conn     *nats.Conn
	log      log15.Logger
}

// NewServer creates a new nats gw server
func NewServer(client *ari.Client, opts *Options) (srv *Server, err error) {

	if client == nil {
		err = errors.New("No client provided")
		return
	}

	if opts == nil {
		opts = &Options{}
	}

	if opts.Logger == nil {
		opts.Logger = log15.New()
	}

	if opts.Parent == nil {
		opts.Parent = context.Background()
	}

	if opts.URL == "" {
		opts.URL = nats.DefaultURL
	}

	srv = &Server{}
	defer func() {
		if err != nil {
			srv = nil // don't return and garbage collect srv on error
		}
	}()

	srv.conn, err = nats.Connect(opts.URL)
	if err != nil {
		return
	}

	srv.ctx, srv.cancel = context.WithCancel(opts.Parent)
	srv.log = opts.Logger
	srv.upstream = client

	return
}

// Listen listens for nats requests and delegates them to the upstream ARI client
func (srv *Server) Listen() error {
	defer srv.conn.Close()

	go srv.application()
	go srv.asterisk()

	<-srv.ctx.Done()

	return nil
}

// Close closes the gateway server
func (srv *Server) Close() {
	if srv == nil {
		return
	}
	srv.cancel()
}
