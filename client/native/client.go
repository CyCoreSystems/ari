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
func New(opts Options) (*ari.Client, error) {

	// Make sure we have an Application defined
	if opts.Application == "" {
		if os.Getenv("ARI_APPLICATION") != "" {
			opts.Application = os.Getenv("ARI_APPLICATION")
		} else {
			opts.Application = uuid.NewV1().String()
		}
	}

	conn := newConn(opts)

	// Connect to Asterisk (websocket)
	if err := conn.Listen(); err != nil {
		conn.Close()
		return nil, err
	}

	playback := &nativePlayback{conn, conn.Bus}
	liveRecording := &nativeLiveRecording{conn}
	logging := &nativeLogging{conn}
	modules := &nativeModules{conn}
	config := &nativeConfig{conn}

	return &ari.Client{
		Cleanup:     conn.Close,
		Playback:    playback,
		Channel:     &nativeChannel{conn, conn.Bus, playback, liveRecording},
		Bridge:      &nativeBridge{conn, conn.Bus, playback, liveRecording},
		Asterisk:    &nativeAsterisk{conn, logging, modules, config},
		Application: &nativeApplication{conn},
		Mailbox:     &nativeMailbox{conn},
		Endpoint:    &nativeEndpoint{conn},
		DeviceState: &nativeDeviceState{conn},
		TextMessage: &nativeTextMessage{conn},
		Sound:       &nativeSound{conn},
		Bus:         conn.Bus,
		Recording: &ari.Recording{
			Live:   liveRecording,
			Stored: &nativeStoredRecording{conn},
		},
	}, nil
}
