package ari

// Playback represents a communication path for interacting
// with an Asterisk server for playback resources
type Playback interface {

	// Get gets the handle to the given playbacl ID
	Get(id string) *PlaybackHandle

	// Data gets the playback data
	Data(id string) (PlaybackData, error)

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

// NewPlaybackHandle builds a handle to the playback id
func NewPlaybackHandle(id string, pb Playback) *PlaybackHandle {
	return &PlaybackHandle{
		id: id,
		p:  pb,
	}
}

// PlaybackHandle is the handle for performing playback operations
type PlaybackHandle struct {
	id string
	p  Playback
}

// ID returns the identifier for the playback
func (ph *PlaybackHandle) ID() string {
	return ph.id
}

// Data gets the playback data
func (ph *PlaybackHandle) Data() (pd PlaybackData, err error) {
	pd, err = ph.p.Data(ph.id)
	return
}

// Control performs the given operation
func (ph *PlaybackHandle) Control(op string) (err error) {
	err = ph.p.Control(ph.id, op)
	return
}

// Stop stops the playback
func (ph *PlaybackHandle) Stop() (err error) {
	err = ph.p.Stop(ph.id)
	return
}

// Match returns true if the event matches the playback
func (ph *PlaybackHandle) Match(e Event) bool {
	p, ok := e.(PlaybackEvent)
	if !ok {
		return false
	}
	ids := p.GetPlaybackIDs()
	for _, i := range ids {
		if i == ph.ID() {
			return true
		}
	}
	return false
}

// Subscribe subscribes the list of channel events
func (ph *PlaybackHandle) Subscribe(n ...string) Subscription {
	if ph == nil {
		return nil
	}
	return ph.p.Subscribe(ph.id, n...)
}
