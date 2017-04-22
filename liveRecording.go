package ari

// LiveRecording represents a communication path interacting with an Asterisk
// server for live recording resources
type LiveRecording interface {

	// Get gets the Recording by type
	Get(key *Key) *LiveRecordingHandle

	// Data gets the data for the live recording
	Data(key *Key) (*LiveRecordingData, error)

	// Stop stops the live recording
	Stop(key *Key) error

	// Pause pauses the live recording
	Pause(key *Key) error

	// Resume resumes the live recording
	Resume(key *Key) error

	// Mute mutes the live recording
	Mute(key *Key) error

	// Unmute unmutes the live recording
	Unmute(key *Key) error

	// Delete deletes the live recording
	Delete(key *Key) error

	// Scrap Stops and deletes the current LiveRecording
	//TODO: reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	Scrap(key *Key) error
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
func NewLiveRecordingHandle(key *Key, s LiveRecording, exec func() (err error)) *LiveRecordingHandle {
	return &LiveRecordingHandle{
		key:  key,
		s:    s,
		exec: exec,
	}
}

// A LiveRecordingHandle is a reference to a live recording that can be operated on
type LiveRecordingHandle struct {
	key      *Key
	s        LiveRecording
	exec     func() (err error)
	executed bool
}

// ID returns the identifier of the live recording
func (s *LiveRecordingHandle) ID() string {
	return s.key.ID
}

// Key returns the key of the live recording
func (s *LiveRecordingHandle) Key() *Key {
	return s.key
}

// Data gets the data for the live recording
func (s *LiveRecordingHandle) Data() (d *LiveRecordingData, err error) {
	d, err = s.s.Data(s.key)
	return
}

// Stop stops and saves the recording
func (s *LiveRecordingHandle) Stop() (err error) {
	err = s.s.Stop(s.key)
	return
}

// Scrap stops and deletes the recording
func (s *LiveRecordingHandle) Scrap() (err error) {
	err = s.s.Scrap(s.key)
	return
}

// Delete deletes the recording
func (s *LiveRecordingHandle) Delete() (err error) {
	err = s.s.Delete(s.key)
	return
}

// Resume resumes the recording
func (s *LiveRecordingHandle) Resume() (err error) {
	err = s.s.Resume(s.key)
	return
}

// Pause pauses the recording
func (s *LiveRecordingHandle) Pause() (err error) {
	err = s.s.Pause(s.key)
	return
}

// Mute mutes the recording
func (s *LiveRecordingHandle) Mute() (err error) {
	err = s.s.Mute(s.key)
	return
}

// Unmute mutes the recording
func (s *LiveRecordingHandle) Unmute() (err error) {
	err = s.s.Unmute(s.key)
	return
}

// Exec executes any staged operations attached to the `LiveRecordingHandle`
func (s *LiveRecordingHandle) Exec() (err error) {
	if !s.executed {
		s.executed = true
		if s.exec != nil {
			err = s.exec()
			s.exec = nil
		}
	}
	return
}
