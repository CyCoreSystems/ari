package ari

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// Options describes client connection options
type Options struct {
	// Application is the the name of this ARI application
	Application string

	// URL is the root URL of the ARI server (asterisk box).
	// Default to http://localhost:8088/ari
	URL string

	// WebsocketURL is the URL for ARI Websocket events.
	// Defaults to the events directory of URL, with a protocol of ws.
	// Usually ws://localhost:8088/ari/events.
	WebsocketURL string

	// Username for ARI authentication
	Username string

	// Password for ARI authentication
	Password string
}

// Client describes an ARI connection to an Asterisk server
// Create one client for each ARI application
type Client struct {
	Options *Options // client options

	WSConfig *websocket.Config // websocket connection configuration

	ReadyChan chan struct{}

	Bus    *Bus        // event bus
	events chan *Event // chan on which events are sent

	httpClient *http.Client

	cancel context.CancelFunc
	mu     sync.Mutex

	TestMode bool // Client is in test mode: no actual HTTP requests will be made
}

// NewClient creates a new Asterisk client
// This function does not attempt to connect to Asterisk itself.
// The ARI URL and websocket URL may also be defined by environment
// variables ARI_URL and ARI_WSURL, respectively; explicitly-supplied
// values for these in the supplied `Options` struct will override
// any environment variables.  Defaults for each are to connect to
// `localhost` at the normal locations for each.
//
// Additionally, username and password for the ARI connection may also
// be supplied by environment variables ARI_USERNAME and ARI_PASSWORD,
// respectively.  There are no defaults for these values.
func NewClient(opts *Options) *Client {
	if opts == nil {
		opts = &Options{}
	}

	// Make sure we have an application name
	if opts.Application == "" {
		opts.Application = uuid.NewV1().String()
	}

	// URL should default to localhost
	if opts.URL == "" {
		if ariURL := os.Getenv("ARI_URL"); ariURL != "" {
			opts.URL = ariURL
		} else {
			opts.URL = "http://localhost:8088/ari"
		}
	}

	// Websocket URL should default to be derived from Url
	if opts.WebsocketURL == "" {
		if ariWsURL := os.Getenv("ARI_WSURL"); ariWsURL != "" {
			opts.WebsocketURL = ariWsURL
		} else {
			opts.WebsocketURL = "ws" + strings.TrimPrefix(opts.URL, "http") + "/events"
		}
	}

	// Add the authorization settings
	if opts.Username == "" && opts.Password == "" {
		if os.Getenv("ARI_USERNAME") != "" {
			opts.Username = os.Getenv("ARI_USERNAME")
			opts.Password = os.Getenv("ARI_PASSWORD")
		} else {
			Logger.Warn("No credentials found; expect failure")
		}
	}

	return &Client{Options: opts, ReadyChan: make(chan struct{})}
}

// Close closes the ARI client
func (c *Client) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Listen maintains and listens to a websocket connection until told to stop.
func (c *Client) Listen(ctx context.Context) (err error) {
	if c.cancel == nil {
		ctx, c.cancel = context.WithCancel(ctx)
	}

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
		c.Bus = StartBus(ctx)
	}

	// Make sure we have a readychan to signal the websocket is up
	if c.ReadyChan == nil {
		c.ReadyChan = make(chan struct{})
	}

	// If we are in test mode, do not connect the websocket, but
	// return and close the ready channel.
	if c.TestMode {
		close(c.ReadyChan)
		return nil
	}

	// Setup and listen on the websocket
	go c.listen(ctx)

	// Wait for the websocket connection to connect or for the context to be cancelled
	select {
	case <-c.ReadyChan:
	case <-ctx.Done():
	}

	return nil
}

func (c *Client) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			Logger.Debug("Exiting websocket on request")
			return
		default:
		}

		Logger.Debug("Connecting to websocket")
		ws, err := websocket.DialConfig(c.WSConfig)
		if err != nil {
			Logger.Error("Failed to create websocket connection to Asterisk:", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		close(c.ReadyChan)

		select {
		case <-ctx.Done():
			Logger.Debug("Closing websocket on request")
			ws.Close()
			ws = nil
			return
		case err := <-c.wsRead(ws):
			Logger.Error("Failure reading from websocket:", "error", err.Error())
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
func (c *Client) wsRead(ws *websocket.Conn) chan error {
	var err error
	ch := make(chan error, 1)

	for {
		var msg Message
		err = AsteriskCodec.Receive(ws, &msg)
		if err != nil {
			ch <- err
			close(ch)
		}
		c.Bus.send(&msg)
	}
}

// basicAuth (stolen from net/http/client.go) creates a basic authentication header
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//
//  Context-related items
//

// clientKey is the key type for contexts
type clientKey string

// NewClientContext returns a context with the client attached
func NewClientContext(ctx context.Context, c *Client) context.Context {
	return NewClientContextWithKey(ctx, c, "_default")
}

// NewClientContextWithKey returns a context with the client attached
// as the given key
func NewClientContextWithKey(ctx context.Context, c *Client, name string) context.Context {
	return context.WithValue(ctx, clientKey(name), c)
}

// ClientFromContext returns the Client stored in the context
// with the default key
func ClientFromContext(ctx context.Context) (*Client, bool) {
	return ClientFromContextWithKey(ctx, "_default")
}

// ClientFromContextWithKey returns the Client stored in the context
// with the given keyname
func ClientFromContextWithKey(ctx context.Context, name string) (*Client, bool) {
	c, ok := ctx.Value(clientKey(name)).(*Client)
	return c, ok
}
