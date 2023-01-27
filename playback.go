package ari

// Playback represents a communication path for interacting
// with an Asterisk server for playback resources
type Playback interface {

	// Get gets the handle to the given playback ID
	Get(key *Key) *PlaybackHandle

	// Data gets the playback data
	Data(key *Key) (*PlaybackData, error)

	// Control performs the given operation on the current playback.  Available operations are:
	//   - restart
	//   - pause
	//   - unpause
	//   - reverse
	//   - forward
	Control(key *Key, op string) error

	// Stop stops the playback
	Stop(key *Key) error

	// Subscribe subscribes on the playback events
	Subscribe(key *Key, n ...string) Subscription
}

// A Player is an entity which can play an audio URI
type Player interface {
	// Play plays the audio using the given playback ID and media URI
	Play(string, interface{}) (*PlaybackHandle, error)

	// StagePlay stages a `Play` operation
	StagePlay(string, interface{}) (*PlaybackHandle, error)

	// Subscribe subscribes the player to events
	Subscribe(n ...string) Subscription
}

// PlaybackOptions describes the parameters to the playback operation
type PlaybackOptions struct {
	// Media URIs to play.
	Media []string `json:"media"`

	// For sounds, selects language for sound.
	Lang string `json:"lang,omitempty"`

	// Number of milliseconds to skip before playing. Only applies to the first URI if multiple media URIs are specified.
	OffsetMs int `json:"offsetms,omitempty"`

	// Number of milliseconds to skip for forward/reverse operations.
	SkipMs int `json:"skipms,omitempty"`
}

// PlaybackData represents the state of a playback
type PlaybackData struct {
	// Key is the cluster-unique identifier for this playback
	Key *Key `json:"key"`

	ID        string `json:"id"` // Unique ID for this playback session
	Language  string `json:"language,omitempty"`
	MediaURI  string `json:"media_uri"`  // URI for the media which is to be played
	State     string `json:"state"`      // State of the playback operation
	TargetURI string `json:"target_uri"` // URI of the channel or bridge on which the media should be played (follows format of 'type':'name')
}

// PlaybackHandle is the handle for performing playback operations
type PlaybackHandle struct {
	key      *Key
	p        Playback
	exec     func(pb *PlaybackHandle) error
	executed bool
}

// NewPlaybackHandle builds a handle to the playback id
func NewPlaybackHandle(key *Key, pb Playback, exec func(pb *PlaybackHandle) error) *PlaybackHandle {
	return &PlaybackHandle{
		key:  key,
		p:    pb,
		exec: exec,
	}
}

// ID returns the identifier for the playback
func (ph *PlaybackHandle) ID() string {
	return ph.key.ID
}

// Key returns the Key for the playback
func (ph *PlaybackHandle) Key() *Key {
	return ph.key
}

// Data gets the playback data
func (ph *PlaybackHandle) Data() (*PlaybackData, error) {
	return ph.p.Data(ph.key)
}

// Control performs the given operation
func (ph *PlaybackHandle) Control(op string) error {
	return ph.p.Control(ph.key, op)
}

// Stop stops the playback
func (ph *PlaybackHandle) Stop() error {
	return ph.p.Stop(ph.key)
}

// Subscribe subscribes the list of channel events
func (ph *PlaybackHandle) Subscribe(n ...string) Subscription {
	if ph == nil {
		return nil
	}

	return ph.p.Subscribe(ph.key, n...)
}

// Exec executes any staged operations
func (ph *PlaybackHandle) Exec() (err error) {
	if !ph.executed {
		ph.executed = true
		if ph.exec != nil {
			err = ph.exec(ph)
			ph.exec = nil
		}
	}

	return
}
