package nats

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/generic"
	"github.com/nats-io/nats"
)

// New creates a new ari.Client connected to a nats based ARI server
func New(url string) (*ari.Client, error) {

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	conn := newConn(c)

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
		Bus:         nil,
		Recording: &ari.Recording{
			Live:   &generic.LiveRecording{conn},
			Stored: &generic.StoredRecording{conn},
		},
	}, nil
}
