package native

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/rid"
	"github.com/CyCoreSystems/ari/v6/stdbus"
	"github.com/rotisserie/eris"

	"github.com/inconshreveable/log15"
	"golang.org/x/net/websocket"
)

// Logger defaults to a discard handler (null output).
// If you wish to enable logging, you can set your own
// handler like so:
//
//	ari.Logger.SetHandler(log15.StderrHandler)
var Logger = log15.New()

func init() {
	// Null logger, by default
	Logger.SetHandler(log15.DiscardHandler())
}

// Options describes the options for connecting to
// a native Asterisk ARI server.
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

	// WebsocketOrigin is the origin to report for the websocket connection.
	// Defaults to http://localhost/
	WebsocketOrigin string

	// Username for ARI authentication
	Username string

	// Password for ARI authentication
	Password string

	// Allow subscribe to all events in Asterisk Server
	SubscribeAll bool
}

// ConnectWithContext creates and connects a new Client to Asterisk ARI.
// Providing a Context allows the caller to cancel the request, which is useful should the caller be in a go routine
// that may be cancelled by it's parent.
func ConnectWithContext(ctx context.Context, chanMsgConnected chan bool, opts *Options) (ari.Client, error) {
	c := New(opts, chanMsgConnected)

	err := c.ConnectWithContext(ctx)
	if err != nil {
		return c, err
	}

	info, err := c.Asterisk().Info(nil)
	if err != nil {
		return c, err
	}

	c.node = info.SystemInfo.EntityID

	return c, err
}

// Connect creates and connects a new Client to Asterisk ARI.
func Connect(opts *Options) (ari.Client, error) {
	c := New(opts, nil)

	err := c.Connect()
	if err != nil {
		return c, err
	}

	info, err := c.Asterisk().Info(nil)
	if err != nil {
		return c, err
	}

	c.node = info.SystemInfo.EntityID

	return c, err
}

// New creates a new ari.Client.  This function should not be used directly unless you need finer control.
// nolint: gocyclo
func New(opts *Options, chanMsgConnected chan bool) *Client {
	if opts == nil {
		opts = &Options{}
	}

	// Make sure we have an Application defined
	if opts.Application == "" {
		if os.Getenv("ARI_APPLICATION") != "" {
			opts.Application = os.Getenv("ARI_APPLICATION")
		} else {
			opts.Application = rid.New("")
		}
	}

	if opts.URL == "" {
		if os.Getenv("ARI_URL") != "" {
			opts.URL = os.Getenv("ARI_URL")
		} else {
			opts.URL = "http://localhost:8088/ari"
		}
	}

	if opts.WebsocketURL == "" {
		if os.Getenv("ARI_WSURL") != "" {
			opts.WebsocketURL = os.Getenv("ARI_WSURL")
		} else {
			opts.WebsocketURL = "ws://localhost:8088/ari/events"
		}
	}

	if opts.WebsocketOrigin == "" {
		if os.Getenv("ARI_WSORIGIN") != "" {
			opts.WebsocketOrigin = os.Getenv("ARI_WSORIGIN")
		} else {
			opts.WebsocketOrigin = "http://localhost/"
		}
	}

	if opts.Username == "" {
		opts.Username = os.Getenv("ARI_USERNAME")
	}

	if opts.Password == "" {
		opts.Password = os.Getenv("ARI_PASSWORD")
	}

	return &Client{
		appName:          opts.Application,
		Options:          opts,
		chanMsgConnected: chanMsgConnected,
	}
}

// Client describes a native ARI client, which connects directly to an Asterisk HTTP-based ARI service.
type Client struct {
	appName string

	node string

	// opts are the configuration options for the client
	Options *Options

	// WSConfig describes the configuration for the websocket connection to Asterisk, from which events will be received.
	WSConfig *websocket.Config

	// golang buffered channel to send Connection Up (true) and Down (false) messages to
	chanMsgConnected chan bool

	// connected is a flag indicating whether the Client is connected to Asterisk
	connected bool

	// Bus the event bus for the Client
	bus ari.Bus

	// httpClient is the reusable HTTP client on which commands to Asterisk are sent
	httpClient http.Client

	cancel context.CancelFunc
}

// ApplicationName returns the client's ARI Application name
func (c *Client) ApplicationName() string {
	return c.appName
}

// Connected indicates whether the websocket is connected
func (c *Client) Connected() bool {
	return c.connected
}

// Close shuts down the ARI client
func (c *Client) Close() {
	c.Bus().Close()

	if c.cancel != nil {
		c.cancel()
	}

	if c.connected {
		c.connected = false
		if c.chanMsgConnected != nil {
			c.chanMsgConnected <- false
		}
	}
}

// Application returns the ARI Application accessors for this client
func (c *Client) Application() ari.Application {
	return &Application{c}
}

// Asterisk returns the ARI Asterisk accessors for this client
func (c *Client) Asterisk() ari.Asterisk {
	return &Asterisk{c}
}

// Bridge returns the ARI Bridge accessors for this client
func (c *Client) Bridge() ari.Bridge {
	return &Bridge{c}
}

// Bus returns the Bus accessors for this client
func (c *Client) Bus() ari.Bus {
	return c.bus
}

// Channel returns the ARI Channel accessors for this client
func (c *Client) Channel() ari.Channel {
	return &Channel{c}
}

// DeviceState returns the ARI DeviceState accessors for this client
func (c *Client) DeviceState() ari.DeviceState {
	return &DeviceState{c}
}

// Endpoint returns the ARI Endpoint accessors for this client
func (c *Client) Endpoint() ari.Endpoint {
	return &Endpoint{c}
}

// LiveRecording returns the ARI LiveRecording accessors for this client
func (c *Client) LiveRecording() ari.LiveRecording {
	return &LiveRecording{c}
}

// Mailbox returns the ARI Mailbox accessors for this client
func (c *Client) Mailbox() ari.Mailbox {
	return &Mailbox{c}
}

// Playback returns the ARI Playback accessors for this client
func (c *Client) Playback() ari.Playback {
	return &Playback{c}
}

// Sound returns the ARI Sound accessors for this client
func (c *Client) Sound() ari.Sound {
	return &Sound{c}
}

// StoredRecording returns the ARI StoredRecording accessors for this client
func (c *Client) StoredRecording() ari.StoredRecording {
	return &StoredRecording{c}
}

// TextMessage returns the ARI TextMessage accessors for this client
func (c *Client) TextMessage() ari.TextMessage {
	return &TextMessage{c}
}

func (c *Client) createWSConfig() (err error) {
	// Construct the websocket connection url
	v := url.Values{}

	v.Set("app", c.Options.Application)

	if c.Options.SubscribeAll {
		v.Set("subscribeAll", "true")
	}

	wsurl := c.Options.WebsocketURL + "?" + v.Encode()

	// Construct a websocket config
	c.WSConfig, err = websocket.NewConfig(wsurl, c.Options.WebsocketOrigin)
	if err != nil {
		return eris.Wrap(err, "Failed to construct websocket config")
	}

	// Add the authorization header
	c.WSConfig.Header.Set("Authorization", "Basic "+basicAuth(c.Options.Username, c.Options.Password))

	return nil
}

// ConnectWithContext sets up and maintains and a websocket connection to Asterisk, passing any received events to the Bus
// Providing a Context allows the caller to cancel the request, which is useful should the caller be in a go routine
// that may be cancelled by it's parent.
func (c *Client) ConnectWithContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	if c.connected {
		cancel()
		return eris.New("already connected")
	}

	if c.Options.Username == "" {
		cancel()
		return eris.New("no username found")
	}

	if c.Options.Password == "" {
		cancel()
		return eris.New("no password found")
	}

	// Construct the websocket config, if we don't already have one
	if c.WSConfig == nil {
		if err := c.createWSConfig(); err != nil {
			cancel()
			return eris.Wrap(err, "failed to create websocket configuration")
		}
	}

	// Make sure the bus is set up
	c.bus = stdbus.New()

	// Setup and listen on the websocket
	wg := new(sync.WaitGroup)

	wg.Add(1)

	go c.listen(ctx, wg)

	wg.Wait()

	return nil
}

// Connect sets up and maintains and a websocket connection to Asterisk, passing any received events to the Bus
func (c *Client) Connect() error {
	return c.ConnectWithContext(context.Background())
}

func (c *Client) listen(ctx context.Context, wg *sync.WaitGroup) {
	var signalUp sync.Once

	for {
		// Exit if our context has been closed
		if ctx.Err() != nil {
			if wg != nil {
				signalUp.Do(wg.Done)
			}

			return
		}

		// Dial Asterisk
		ws, err := websocket.DialConfig(c.WSConfig)
		if err != nil {
			Logger.Error("failed to connect to Asterisk", "error", err)
			time.Sleep(time.Second)

			continue
		}

		info, err := c.Asterisk().Info(nil)
		if err != nil {
			Logger.Error("failed to get info from Asterisk", "error", err)

			time.Sleep(time.Second)

			continue
		}

		if c.node != "" && c.node != info.SystemInfo.EntityID {
			c.node = info.SystemInfo.EntityID
		}

		// We are connected
		c.connected = true
		if c.chanMsgConnected != nil {
			c.chanMsgConnected <- true
		}

		// Signal that we are connected (the first time only)
		if wg != nil {
			signalUp.Do(wg.Done)
		}

		// Wait for context closure or read error
		select {
		case <-ctx.Done():
			// Make sure our websocket connection is closed before looping
			c.connected = false
			if c.chanMsgConnected != nil {
				c.chanMsgConnected <- false
			}

		case err = <-c.wsRead(ws):
			Logger.Error("read failure on websocket", "error", err)

			c.connected = false
			if c.chanMsgConnected != nil {
				c.chanMsgConnected <- false
			}

			time.Sleep(10 * time.Millisecond)
		}

		err = ws.Close()
		if err != nil {
			Logger.Debug("failed to close websocket", "error", err)
		}

		if c.connected {
			c.connected = false
			if c.chanMsgConnected != nil {
				c.chanMsgConnected <- false
			}
		}
	}
}

// wsRead loops for the duration of a websocket connection,
// reading messages, decoding them to events, and passing
// them to the event bus.
func (c *Client) wsRead(ws *websocket.Conn) chan error {
	errChan := make(chan error, 1)

	go func() {
		for {
			var data []byte

			err := websocket.Message.Receive(ws, &data)
			if err != nil {
				errChan <- eris.Wrap(err, "failed to receive websocket message")
				return
			}

			e, err := ari.DecodeEvent(data)
			if err != nil {
				errChan <- eris.Wrap(err, "failed to devoce websocket message to event")
			}

			c.bus.Send(e)
		}
	}()

	return errChan
}

// stamp imprints the node metadata onto the given Key
func (c *Client) stamp(key *ari.Key) *ari.Key {
	if key == nil {
		key = &ari.Key{}
	}

	ret := *key
	ret.App = c.appName
	ret.Node = c.node

	return &ret
}

// basicAuth (stolen from net/http/client.go) creates a basic authentication header
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
