package nc

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// Options is the list options
type Options struct {

	// URL is the nats URL
	URL string

	// ReadOperationRetryCount is the amount of times to retry a read operation
	ReadOperationRetryCount int

	// RequestTimeout is the timeout duration of a request
	RequestTimeout time.Duration
}

// New creates a new ari.Client connected to a gateway ARI server via NATS
func New(opts Options) (cl *ari.Client, err error) {

	var nc *nats.Conn
	nc, err = nats.Connect(opts.URL)
	if err != nil {
		return
	}

	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = DefaultRequestTimeout
	}

	conn := &Conn{
		opts: opts,
		conn: nc,
	}

	playback := natsPlayback{conn}
	bus := &natsBus{conn}
	liveRecording := &natsLiveRecording{conn}
	storedRecording := &natsStoredRecording{conn}
	logging := &natsLogging{conn}
	modules := &natsModules{conn}
	config := &natsConfig{conn}

	cl = &ari.Client{
		Cleanup:     func() error { nc.Close(); return nil },
		Asterisk:    &natsAsterisk{conn, logging, modules, config},
		Application: &natsApplication{conn},
		Bridge:      &natsBridge{conn, &playback, liveRecording},
		Channel:     &natsChannel{conn, &playback, liveRecording},
		DeviceState: &natsDeviceState{conn},
		Mailbox:     &natsMailbox{conn},
		Sound:       &natsSound{conn},
		Playback:    &playback,
		Recording: &ari.Recording{
			Live:   liveRecording,
			Stored: storedRecording,
		},
		Bus: bus,
	}

	return
}
