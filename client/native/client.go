package native

import (
	"os"

	"github.com/CyCoreSystems/ari"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

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

	// Username for ARI authentication
	Username string

	// Password for ARI authentication
	Password string

	// Optional context to act as parent
	Context context.Context
}

// New creates a new ari.Client connected to a native ARI server
func New(opts Options) (ari.Client, error) {

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

	if opts.Username == "" {
		opts.Username = os.Getenv("ARI_USERNAME")
	}
	if opts.Password == "" {
		opts.Password = os.Getenv("ARI_PASSWORD")
	}

	conn := newConn(opts)

	// Connect to Asterisk (websocket)
	if err := conn.Listen(); err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		appName: opts.Application,
		conn:    conn,
	}, nil
}

// Client describes a native ARI client, which connects directly to an Asterisk HTTP-based ARI service.
type Client struct {
	appName string

	conn *Conn
}

// ApplicationName returns the client's ARI Application name
func (c *Client) ApplicationName() string {
	return c.appName
}

// Close shuts down the ARI client
func (c *Client) Close() {
	c.Bus().Close()
	c.conn.Close()
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
	return c.conn.Bus
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
