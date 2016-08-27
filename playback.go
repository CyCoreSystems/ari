package ari

// Playback represents a communication path for interacting
// with an Asterisk server for playback resources
type Playback interface {

	// Data gets the playback data
	Data(id string) (PlaybackData, error)

	// Control performs the given operation on the current playback
	Control(id string, op string) error

	// Stop stops the playback
	Stop(id string) error
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
