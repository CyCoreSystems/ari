package native

import (
	"net/http"
	"sync"

	"github.com/CyCoreSystems/ari"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
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
}

// Conn is a connection to a native ARI server
type Conn struct {
	Options *Options // client options

	WSConfig *websocket.Config // websocket connection configuration

	ReadyChan chan struct{}

	Bus    *ari.Bus    // event bus
	events chan *Event // chan on which events are sent

	httpClient *http.Client

	cancel context.CancelFunc
	mu     sync.Mutex
}

// New creates a new ari.Client connected to a native ARI server
func New(_ *Options) (*ari.Client, error) {

	var conn Conn //TODO: create connection from opts

	//TODO: populate client
	return &ari.Client{
		Channel:     &nativeChannel{&conn},
		Asterisk:    &nativeAsterisk{&conn},
		Application: &nativeApplication{&conn},
	}, nil
}
