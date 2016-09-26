package native

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// Conn is a connection to a native ARI server
type Conn struct {
	Options Options // client options

	WSConfig *websocket.Config // websocket connection configuration

	ReadyChan chan struct{}

	Bus    ari.Bus        // event bus
	events chan ari.Event // chan on which events are sent

	httpClient http.Client

	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func newConn(opts Options) (c *Conn) {

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	c = &Conn{}
	c.Options = opts
	c.ctx, c.cancel = context.WithCancel(opts.Context)

	return
}

// Close closes the ARI client
func (c *Conn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

// Listen maintains and listens to a websocket connection until told to stop.
func (c *Conn) Listen() (err error) {

	// Construct the websocket config, if we don't already have one
	if c.WSConfig == nil {
		// Construct the websocket connection url
		v := url.Values{}
		v.Set("app", c.Options.Application)
		wsurl := c.Options.WebsocketURL + "?" + v.Encode()

		// Construct a websocket.Config
		c.WSConfig, err = websocket.NewConfig(wsurl, "http://localhost/")
		if err != nil {
			Logger.Error("Failed to construct a valid websocket config:", err.Error())
			return fmt.Errorf("Failed to construct websocket config: %s", err.Error())
		}

		// Add the authorization header
		if c.Options.Username != "" && c.Options.Password != "" {
			c.WSConfig.Header.Set("Authorization", "Basic "+basicAuth(c.Options.Username, c.Options.Password))
		} else if os.Getenv("ARI_USERNAME") != "" {
			c.WSConfig.Header.Set("Authorization", "Basic "+basicAuth(os.Getenv("ARI_USERNAME"), os.Getenv("ARI_PASSWORD")))
		} else {
			Logger.Warn("No credentials found; expect failure")
		}
	}

	// Make sure the bus is set up
	if c.Bus == nil {
		c.Bus = startbus(c.ctx)
	}

	// Make sure we have a readychan to signal the websocket is up
	if c.ReadyChan == nil {
		c.ReadyChan = make(chan struct{})
	}

	// If we are in test mode, do not connect the websocket, but
	// return and close the ready channel.
	//if c.TestMode {
	//	close(c.ReadyChan)
	//	return nil
	//}

	// Setup and listen on the websocket
	go c.listen(c.ctx)

	// Wait for the websocket connection to connect or for the context to be cancelled
	select {
	case <-c.ReadyChan:
	case <-c.ctx.Done():
		return c.ctx.Err()
	}

	return nil
}

func (c *Conn) listen(ctx context.Context) {
	var ws *websocket.Conn
	var err error

	// Close the websocket if our context is closed
	go func() {
		<-ctx.Done()
		Logger.Info("closing websocket on request")
		if ws != nil {
			ws.Close()
		}
	}()

	for {

		select {
		case <-ctx.Done():
			return
		default:
		}

		Logger.Debug("Connecting to websocket")

		ws, err = websocket.DialConfig(c.WSConfig)
		if err != nil {
			Logger.Error("Failed to create websocket connection to Asterisk", "error", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		close(c.ReadyChan)

		err = c.wsRead(ws)
		if err != nil {
			Logger.Error("Failure reading from websocket", "error", err.Error())
		}

		// Clean up
		ws.Close()
		ws = nil

		c.ReadyChan = make(chan struct{})

		// Don't restart too quickly
		Logger.Info("Waiting 10ms to restart websocket")
		time.Sleep(10 * time.Millisecond)
	}

}

// wsRead loops for the duration of a websocket connection,
// reading messages, decoding them to events, and passing
// them to the event bus.
func (c *Conn) wsRead(ws *websocket.Conn) (err error) {
	for {
		var msg ari.Message
		err = AsteriskCodec.Receive(ws, &msg)
		if err != nil {
			Logger.Error("Error decoding websocket message", "error", err)
			return
		}

		c.Bus.Send(&msg)
	}
}

// basicAuth (stolen from net/http/client.go) creates a basic authentication header
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Logger defaults to a discard handler (null output).
// If you wish to enable logging, you can set your own
// handler like so:
// 		ari.Logger.SetHandler(log15.StderrHandler)
//
var Logger = log15.New()

func init() {
	// Null logger, by default
	Logger.SetHandler(log15.DiscardHandler())
}
