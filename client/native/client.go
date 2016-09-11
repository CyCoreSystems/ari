package native

import "github.com/CyCoreSystems/ari"

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

// New creates a new ari.Client connected to a native ARI server
func New(opts *Options) (*ari.Client, error) {

	conn := newConn(opts)

	if err := conn.Listen(nil); err != nil {
		conn.Close()
		return nil, err
	}

	playback := &nativePlayback{conn}

	return &ari.Client{
		Playback:    playback,
		Channel:     &nativeChannel{conn, playback},
		Bridge:      &nativeBridge{conn, playback},
		Asterisk:    &nativeAsterisk{conn},
		Application: &nativeApplication{conn},
		Mailbox:     &nativeMailbox{conn},
		Endpoint:    &nativeEndpoint{conn},
		DeviceState: &nativeDeviceState{conn},
		TextMessage: &nativeTextMessage{conn},
		Sound:       &nativeSound{conn},
		Event:       &nativeEvent{conn},
		Bus:         conn.Bus,
		Recording: &ari.Recording{
			Live:   &nativeLiveRecording{conn},
			Stored: &nativeStoredRecording{conn},
		},
	}, nil
}
