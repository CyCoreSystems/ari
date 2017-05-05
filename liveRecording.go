package ari

import "sync"

// LiveRecording represents a communication path interacting with an Asterisk
// server for live recording resources
type LiveRecording interface {

	// Get gets the Recording by type
	Get(key *Key) *LiveRecordingHandle

	// Data gets the data for the live recording
	Data(key *Key) (*LiveRecordingData, error)

	// Stop stops the live recording
	Stop(key *Key) (*StoredRecordingHandle, error)

	// Pause pauses the live recording
	Pause(key *Key) error

	// Resume resumes the live recording
	Resume(key *Key) error

	// Mute mutes the live recording
	Mute(key *Key) error

	// Unmute unmutes the live recording
	Unmute(key *Key) error

	// Scrap Stops and deletes the current LiveRecording
	Scrap(key *Key) error

	// Subscribe subscribes to events
	Subscribe(key *Key, n ...string) Subscription
}

// LiveRecordingData is the data for a live recording
type LiveRecordingData struct {
	// Key is the cluster-unique identifier for this live recording
	Key *Key `json:"key"`

	Cause     string      `json:"cause,omitempty"`            // If failed, the cause of the failure
	Duration  DurationSec `json:"duration,omitempty"`         // Length of recording in seconds
	Format    string      `json:"format"`                     // Format of recording (wav, gsm, etc)
	Name      string      `json:"name"`                       // (base) name for the recording
	Silence   DurationSec `json:"silence_duration,omitempty"` // If silence was detected in the recording, the duration in seconds of that silence (requires that maxSilenceSeconds be non-zero)
	State     string      `json:"state"`                      // Current state of the recording
	Talking   DurationSec `json:"talking_duration,omitempty"` // Duration of talking, in seconds, that has been detected in the recording (requires that maxSilenceSeconds be non-zero)
	TargetURI string      `json:"target_uri"`                 // URI for the channel or bridge which is being recorded (TODO: figure out format for this)
}

// ID returns the identifier of the live recording
func (s *LiveRecordingData) ID() string {
	return s.Name
}

// NewLiveRecordingHandle creates a new live recording handle
func NewLiveRecordingHandle(key *Key, r LiveRecording, exec func(*LiveRecordingHandle) (err error)) *LiveRecordingHandle {
	return &LiveRecordingHandle{
		key:  key,
		r:    r,
		exec: exec,
	}
}

// A LiveRecordingHandle is a reference to a live recording that can be operated on
type LiveRecordingHandle struct {
	key      *Key
	r        LiveRecording
	exec     func(*LiveRecordingHandle) (err error)
	executed bool

	mu sync.Mutex
}

// ID returns the identifier of the live recording
func (h *LiveRecordingHandle) ID() string {
	return h.key.ID
}

// Key returns the key of the live recording
func (h *LiveRecordingHandle) Key() *Key {
	return h.key
}

// Data gets the data for the live recording
func (h *LiveRecordingHandle) Data() (*LiveRecordingData, error) {
	return h.r.Data(h.key)
}

// Stop stops and saves the recording
func (h *LiveRecordingHandle) Stop() (*StoredRecordingHandle, error) {
	return h.r.Stop(h.key)
}

// Scrap stops and deletes the recording
func (h *LiveRecordingHandle) Scrap() error {
	return h.r.Scrap(h.key)
}

// Resume resumes the recording
func (h *LiveRecordingHandle) Resume() error {
	return h.r.Resume(h.key)
}

// Pause pauses the recording
func (h *LiveRecordingHandle) Pause() error {
	return h.r.Pause(h.key)
}

// Mute mutes the recording
func (h *LiveRecordingHandle) Mute() error {
	return h.r.Mute(h.key)
}

// Unmute mutes the recording
func (h *LiveRecordingHandle) Unmute() error {
	return h.r.Unmute(h.key)
}

// Exec executes any staged operations attached to the `LiveRecordingHandle`
func (h *LiveRecordingHandle) Exec() (err error) {
	h.mu.Lock()

	if !h.executed {
		h.executed = true
		if h.exec != nil {
			err = h.exec(h)
			h.exec = nil
		}
	}

	h.mu.Unlock()
	return
}

// Subscribe subscribes the recording handle's underlying recorder to
// the provided event types.
func (h *LiveRecordingHandle) Subscribe(n ...string) Subscription {
	return h.r.Subscribe(h.key, n...)
}
