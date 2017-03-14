package native

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync"
	"time"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/stdbus"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

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
}

// Connect creates and connects a new Client to Asterisk ARI.
func Connect(ctx context.Context, opts Options) (ari.Client, error) {
	c := New(opts)

	err := c.Connect(ctx)

	return c, err
}

// New creates a new ari.Client.  This function should not be used directly unless you need finer control.
func New(opts Options) *Client {

	// Make sure we have an Application defined
	if opts.Application == "" {
		if os.Getenv("ARI_APPLICATION") != "" {
			opts.Application = os.Getenv("ARI_APPLICATION")
		} else {
			opts.Application = uuid.NewV1().String()
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
		opts.WebsocketOrigin = "http://localhost/"
	}

	if opts.Username == "" {
		opts.Username = os.Getenv("ARI_USERNAME")
	}
	if opts.Password == "" {
		opts.Password = os.Getenv("ARI_PASSWORD")
	}

	return &Client{
		appName: opts.Application,
		Options: &opts,
	}
}

// Client describes a native ARI client, which connects directly to an Asterisk HTTP-based ARI service.
type Client struct {
	appName string

	// opts are the configuration options for the client
	Options *Options

	// WSConfig describes the configuration for the websocket connection to Asterisk, from which events will be received.
	WSConfig *websocket.Config

	// Connected is a flag indicating whether the Client is connected to Asterisk
	Connected bool

	// Bus the event bus for the Client
	bus ari.Bus

	// httpClient is the reusable HTTP client on which commands to Asterisk are sent
	httpClient http.Client

	// wsConn is the current websocket connection
	wsConn *websocket.Conn

	cancel context.CancelFunc
	mu     sync.Mutex
}

// ApplicationName returns the client's ARI Application name
func (c *Client) ApplicationName() string {
	return c.appName
}

// Close shuts down the ARI client
func (c *Client) Close() {
	c.Bus().Close()

	if c.cancel != nil {
		c.cancel()
	}

	c.Connected = false
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
	wsurl := c.Options.WebsocketURL + "?" + v.Encode()

	// Construct a websocket config
	c.WSConfig, err = websocket.NewConfig(wsurl, c.Options.WebsocketOrigin)
	if err != nil {
		return errors.Wrap(err, "Failed to construct websocket config")
	}

	// Add the authorization header
	c.WSConfig.Header.Set("Authorization", "Basic "+basicAuth(c.Options.Username, c.Options.Password))
	return nil
}

// Connect sets up and maintains and a websocket connection to Asterisk, passing any received events to the Bus
func (c *Client) Connect(ctx context.Context) error {
	if c.Connected {
		return errors.New("already connected")
	}

	if c.Options.Username == "" {
		return errors.New("no username found")
	}
	if c.Options.Password == "" {
		return errors.New("no password found")
	}

	// Construct the websocket config, if we don't already have one
	if c.WSConfig == nil {
		if err := c.createWSConfig(); err != nil {
			return errors.Wrap(err, "failed to create websocket configuration")
		}
	}

	// Make sure the bus is set up
	c.bus = stdbus.New()

	// Setup and listen on the websocket
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go c.listen(ctx, wg)
	wg.Wait()
	c.Connected = true

	return nil
}

func (c *Client) listen(ctx context.Context, wg *sync.WaitGroup) {
	var signalUp sync.Once

	for {
		// Exit if our context has been closed
		if ctx.Err() != nil {
			return
		}

		// Dial Asterisk
		ws, err := websocket.DialConfig(c.WSConfig)
		if err != nil {
			Logger.Error("failed to connect to Asterisk", "error", err)
			time.Sleep(time.Second)
			continue
		}

		// Signal that we are connected (the first time only)
		if wg != nil {
			signalUp.Do(wg.Done)
		}

		// Wait for context closure or read error
		select {
		case <-ctx.Done():
		case err := <-c.wsRead(ws):
			Logger.Error("read failure on websocket", "error", err)
			time.Sleep(10 * time.Millisecond)
		}

		// Make sure our websocket connection is closed before looping
		ws.Close()
	}
}

// wsRead loops for the duration of a websocket connection,
// reading messages, decoding them to events, and passing
// them to the event bus.
func (c *Client) wsRead(ws *websocket.Conn) chan error {
	errChan := make(chan error, 1)

	go func() {
		for {
			var msg ari.Message
			err := AsteriskCodec.Receive(ws, &msg)
			if err != nil {
				errChan <- errors.Wrap(err, "failed to decode websocket message")
				return
			}

			c.bus.Send(&msg)
		}
	}()

	return errChan
}

// basicAuth (stolen from net/http/client.go) creates a basic authentication header
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Marshal is a no-op to implement websocket.Codec.  Asterisk
// websocket connections should never have the client send any data
func marshal(v interface{}) (data []byte, payloadType byte, err error) {
	return
}

// Unmarshal implements websocket.Codec
func unmarshal(data []byte, payloadType byte, v interface{}) error {
	data = append(data, '\n')

	e, ok := v.(*ari.Message)
	if !ok {
		return fmt.Errorf("Cannot cast receiver to a Message when it is of type %v", reflect.TypeOf(v))
	}

	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	// Store the raw data
	e.SetRaw(data)

	return nil
}

// AsteriskCodec is a websocket Codec for Asterisk messages
var AsteriskCodec = websocket.Codec{
	Marshal:   marshal,
	Unmarshal: unmarshal,
}
