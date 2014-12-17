package ari

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"code.google.com/p/go.net/websocket"
)

// Client describes an ARI connection to an Asterisk server
// Create one client for each ARI application
type Client struct {
	Application string // ARI Application name
	Url         string // root URL of ARI server (asterisk box), e.g. http://localhost:8088/ari
	WsUrl       string // URL of ARI Websocket events, e.g. ws://localhost:8088/ari/events
	username    string // username for ARI authentication
	password    string // password for ARI authentication

	WSConfig  *websocket.Config // websocket connection configuration
	WebSocket *websocket.Conn   // websocket connection used to receive events

	Events chan *Event // chan on which events are sent

	StopRequested bool // Whether client has been requested to stop
}

// NewClient creates a new Asterisk client
// This function does not attempt to connect to Asterisk itself.
func NewClient(appName, aurl, wsurl, username, password string) (Client, error) {
	c := Client{
		Application: appName,
		Url:         aurl,
		WsUrl:       wsurl,
		username:    username,
		password:    password,
	}

	// Construct the websocket connection url
	v := url.Values{}
	v.Set("app", c.Application)
	wsurl = c.WsUrl + "?" + v.Encode()

	// Construct a websocket.Config
	wsConfig, err := websocket.NewConfig(wsurl, "http://localhost/")
	if err != nil {
		return c, fmt.Errorf("Failed to construct websocket config:", err.Error())
	}

	// Add the authorization header
	wsConfig.Header.Set("Authorization", "Basic "+basicAuth(c.username, c.password))

	// Store the websocket configuration to our struct
	c.WSConfig = wsConfig

	// Create the (buffered) Events chan
	c.Events = make(chan *Event, 100)

	return c, nil
}

// Go maintains a websocket connection until told to stop
func (c *Client) Go() {
	// Dereference so we don't block
	go c._goloop()

	return
}

// GoWait maintains a websocket connection until told to stop
// but does not return until a successful connection is established
func (c *Client) GoWait() error {
	go c._goloop()

	// Wait until we see the WebSocket come up
	for c.WebSocket == nil {
		time.Sleep(100 * time.Millisecond)
		// Exit if a stop has been requested
		if c.StopRequested {
			return fmt.Errorf("Timed out waiting for websocket")
		}
	}

	return nil
}

// _goloop maintains the websocket connection
func (c *Client) _goloop() {
	for c.StopRequested == false {
		glog.V(9).Infoln("Connecting to websocket")
		err := c.Connect()

		// Exit if we were _requested_ to stop
		if c.StopRequested == true {
			glog.V(3).Infoln("Websocket connection closed by request; exiting")
			return
		}

		// An error indicates we failed to connect, not that we lost
		// an active connection.  Hence, insert a delay before the
		// next attempt
		if err != nil {
			glog.Errorln("Failed to open websocket connection: ", err)
			time.Sleep(1 * time.Second)
			continue
		}

		glog.Infoln("Websocket connection died (reconnecting)")
	}
}

// Attempt to connect to the websocket connection
func (c *Client) Connect() error {
	// Connect to the websocket
	ws, err := websocket.DialConfig(c.WSConfig)
	if err != nil {
		return fmt.Errorf("Failed to create websocket connection to Asterisk:", err.Error())
	}
	c.WebSocket = ws

	// Start listening on the websocket
	return c.Listen()
}

// Close an active client connection
func (c *Client) Close() {
	// Flag that a stop has been requested
	c.StopRequested = true

	// Call the real close routine
	c._close()
}

// Internal close function; does not set StopRequested flag
func (c *Client) _close() {
	// Close the websocket connection
	if c.WebSocket != nil {
		err := c.WebSocket.Close()
		if err != nil {
			glog.Warningln("Failed to close websocket")
		}
		c.WebSocket = nil
	}
}

// Listen waits for events on the Websocket connection
func (c *Client) Listen() error {
	for c.WebSocket != nil {
		var data []byte

		// Wait for a message
		err := websocket.Message.Receive(c.WebSocket, &data)

		// If we got an error, signal failure and exit
		if err != nil {
			glog.Warningln("Failure in websocket (or connection lost):", err)
			c._close()
			return err
		}

		// If we got a non-empty message, parse it into an event
		if len(data) > 0 {
			go c.ParseMessage(data)
		}
	}
	glog.V(3).Infoln("Websocket is gone")
	return nil
}

// Parse a websocket message and send it to the Events chan
func (c *Client) ParseMessage(data []byte) {
	// Attempt to construct a message out of the data
	m, err := NewMessage(data)
	if err != nil {
		glog.Errorln("Failed to read message from websocket:", err.Error())
		return
	}

	// Attempt to decode the message as an event
	var e Event
	err = m.DecodeAs(&e)
	if err != nil {
		glog.Errorln("Failed to decode message as an event.  Unhandled message type:", m.Type)
		return
	}

	// Send the event to the Events chan
	c.Events <- &e

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
