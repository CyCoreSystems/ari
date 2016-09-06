package native

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/generic"
	"golang.org/x/net/context"
)

// Options describes the options for connecting to
// a generic. Asterisk ARI server.
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

// New creates a new ari.Client connected to a generic. ARI server
func New(opts *Options) (*ari.Client, error) {

	conn := newConn(opts)

	if err := conn.Listen(context.Background()); err != nil {
		conn.Close()
		return nil, err
	}

	playback := &generic.Playback{conn}

	return &ari.Client{
		Playback:    playback,
		Channel:     &generic.Channel{conn, playback},
		Bridge:      &generic.Bridge{conn, playback},
		Asterisk:    &generic.Asterisk{conn},
		Application: &generic.Application{conn},
		Mailbox:     &generic.Mailbox{conn},
		Endpoint:    &generic.Endpoint{conn},
		DeviceState: &generic.DeviceState{conn},
		TextMessage: &generic.TextMessage{conn},
		Sound:       &generic.Sound{conn},
		Bus:         conn.Bus,
		Recording: &ari.Recording{
			Live:   &generic.LiveRecording{conn},
			Stored: &generic.StoredRecording{conn},
		},
	}, nil
}
