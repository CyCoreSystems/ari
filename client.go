package ari

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

// Client connection options
type Options struct {
	Application string // ARI Application name
	Url         string // root URL of ARI server (asterisk box), e.g. http://localhost:8088/ari
	WsUrl       string // URL of ARI Websocket events, e.g. ws://localhost:8088/ari/events
	Username    string // username for ARI authentication
	Password    string // password for ARI authentication
}

// Client describes an ARI connection to an Asterisk server
// Create one client for each ARI application
type Client struct {
	Options *Options // client options

	WSConfig *websocket.Config // websocket connection configuration

	Bus    *Bus        // event bus
	events chan *Event // chan on which events are sent

	httpClient *gorequest.SuperAgent // reusable HTTP client

	mu sync.Mutex
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
	if opts.Url == "" {
		if ariUrl := os.Getenv("ARI_URL"); ariUrl != "" {
			opts.Url = ariUrl
		} else {
			opts.Url = "http://localhost:8088/ari"
		}
	}

	// Websocket URL should default to be derived from Url
	if opts.WsUrl == "" {
		if ariWsUrl := os.Getenv("ARI_WSURL"); ariWsUrl != "" {
			opts.WsUrl = ariWsUrl
		} else {
			opts.WsUrl = "ws" + strings.TrimPrefix(opts.Url, "http") + "/events"
		}
	}

	return &Client{Options: opts}
}

// Listen maintains and listens to a websocket connection until told to stop.
func (c *Client) Listen(ctx context.Context) (err error) {
	// Construct the websocket config, if we don't already have one
	if c.WSConfig == nil {
		// Construct the websocket connection url
		v := url.Values{}
		v.Set("app", c.Options.Application)
		wsurl := c.Options.WsUrl + "?" + v.Encode()

		// Construct a websocket.Config
		c.WSConfig, err = websocket.NewConfig(wsurl, "http://localhost/")
		if err != nil {
			Logger.Error("Failed to construct a calid websocket config:", err.Error())
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

	go c.listen(ctx)
	return nil
}

func (c *Client) listen(ctx context.Context) {
	var err error
	var ws *websocket.Conn
	var stop bool

	go func() {
		for !stop {
			Logger.Debug("Connecting to websocket")
			ws, err = websocket.DialConfig(c.WSConfig)
			if err != nil {
				Logger.Error("Failed to create websocket connection to Asterisk:", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

		ReadLoop:
			for !stop {
				var msg Message
				err := AsteriskCodec.Receive(ws, &msg)
				if err != nil {
					Logger.Error("Failure in websocket connection:", err.Error())
					break ReadLoop
				}
				c.Bus.send(&msg)
			}

			// Clean up
			ws.Close()
			ws = nil

			// Don't restart too quickly
			Logger.Info("Waiting 10ms to restart websocket")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Wait for stop
	<-ctx.Done()
	stop = true
	ws.Close()
	return
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
