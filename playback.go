package ari

// Playback represents a communication path for interacting
// with an Asterisk server for playback resources
type Playback interface {

	// Get gets the handle to the given playbacl ID
	Get(id string) PlaybackHandle

	// Data gets the playback data
	Data(id string) (*PlaybackData, error)

	// Control performs the given operation on the current playback
	Control(id string, op string) error

	// Stop stops the playback
	Stop(id string) error

	// Subscribe subscribes on the playback events
	Subscribe(id string, n ...string) Subscription
}

// Playbacker contains a playback transport
type Playbacker interface {
	Playback() Playback
}

// PlaybackData represents the state of a playback
type PlaybackData struct {
	ID        string `json:"id"` // Unique ID for this playback session
	Language  string `json:"language,omitempty"`
	MediaURI  string `json:"media_uri"`  // URI for the media which is to be played
	State     string `json:"state"`      // State of the playback operation
	TargetURI string `json:"target_uri"` // URI of the channel or bridge on which the media should be played (follows format of 'type':'name')
}

// PlaybackHandle is the handle for performing playback operations
type PlaybackHandle interface {

	// ID returns the identifier for the playback
	ID() string

	// Data gets the playback data
	Data() (pd *PlaybackData, err error)

	// Control performs the given operation
	Control(op string) (err error)

	// Stop stops the playback
	Stop() (err error)

	// Match returns true if the event matches the playback
	Match(e Event) bool

	// Subscribe subscribes the list of channel events
	Subscribe(n ...string) Subscription
}
