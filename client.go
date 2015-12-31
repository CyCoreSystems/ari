package ari

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
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

	StopChan <-chan struct{} // Stop signal channel

	httpClient *gorequest.SuperAgent // reusable HTTP client
}

// NewClient creates a new Asterisk client
// This function does not attempt to connect to Asterisk itself.
func NewClient(appName, aurl, wsurl, username, password string) (Client, error) {
	return NewClientWithStop(appName, aurl, wsurl, username, password, nil)
}

// NewClientWithStop creates a new Asterisk client with a stop channel
// This function does not attempt to connect to Asterisk itself.
func NewClientWithStop(appName, aurl, wsurl, username, password string, stopChan <-chan struct{}) (Client, error) {
	c := Client{
		Application: appName,
		Url:         aurl,
		WsUrl:       wsurl,
		username:    username,
		password:    password,
		StopChan:    stopChan,
	}

	// Construct the websocket connection url
	v := url.Values{}
	v.Set("app", c.Application)
	wsurl = c.WsUrl + "?" + v.Encode()

	// Construct a websocket.Config
	wsConfig, err := websocket.NewConfig(wsurl, "http://localhost/")
	if err != nil {
		return c, fmt.Errorf("Failed to construct websocket config: %s", err.Error())
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

// _goloop maintains the websocket connection
func (c *Client) _goloop() {
	for {
		Logger.Debug("Connecting to websocket")
		select {
		case <-c.listen():
			Logger.Warn("Websocket connection lost")
		case <-c.StopChan:
			Logger.Info("Stop requested; exiting")
			return
		}
		Logger.Info("Waiting 100ms to restart websocket")
		time.Sleep(100 * time.Millisecond)
	}
}

// listen opens a websocket, processes messages,
// and returns a chan bool to indicate the socket
// is closed
func (c *Client) listen() chan bool {
	closedChan := make(chan bool)

	// Connect to the websocket
	ws, err := websocket.DialConfig(c.WSConfig)
	if err != nil {
		Logger.Error("Failed to create websocket connection to Asterisk:", err.Error())
		close(closedChan)
		return closedChan
	}
	c.WebSocket = ws

	// Loop, receiving messages
	go c.Listen(closedChan)

	// TODO: Signal that we are ready

	return closedChan
}

// Close an active client connection
func (c *Client) Close() {
	c.closeWebsocket()
}

// closeWebsocket closes the websocket connection, if it exists
func (c *Client) closeWebsocket() {
	// Close the websocket connection
	if c.WebSocket != nil {
		err := c.WebSocket.Close()
		if err != nil {
			Logger.Error("Failed to close websocket")
		}
		c.WebSocket = nil
	}
}

// Listen waits for events on the Websocket connection
func (c *Client) Listen(closedChan chan bool) {
	defer func() {
		c.closeWebsocket()
		close(closedChan)
	}()

	for c.WebSocket != nil {
		var data []byte

		// Wait for a message
		err := websocket.Message.Receive(c.WebSocket, &data)

		// If we got an error, signal failure and exit
		if err != nil {
			Logger.Error("Failure in websocket (or connection lost):", err)
			return
		}

		// If we got a non-empty message, parse it into an event
		if len(data) > 0 {
			go c.ParseMessage(data)
		}
	}
	Logger.Info("Websocket is gone")
	return
}

// Parse a websocket message and send it to the Events chan
func (c *Client) ParseMessage(data []byte) {
	// Attempt to construct a message out of the data
	m, err := NewMessage(data)
	if err != nil {
		Logger.Error("Failed to read message from websocket:", err.Error())
		return
	}

	// Attempt to decode the message as an event
	var e Event
	err = m.DecodeAs(&e)
	if err != nil {
		Logger.Error("Failed to decode message as an event.  Unhandled message type:", m.Type)
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
