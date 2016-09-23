package ari

// LiveRecording represents a communication path interacting with an Asterisk
// server for live recording resources
type LiveRecording interface {

	// Get gets the Recording by type
	Get(name string) *LiveRecordingHandle

	// Data gets the data for the live recording
	Data(name string) (LiveRecordingData, error)

	// Stop stops the live recording
	Stop(name string) error

	// Pause pauses the live recording
	Pause(name string) error

	// Resume resumes the live recording
	Resume(name string) error

	// Mute mutes the live recording
	Mute(name string) error

	// Unmute unmutes the live recording
	Unmute(name string) error

	// Delete deletes the live recording
	Delete(name string) error

	// Scrap Stops and deletes the current LiveRecording
	//TODO: reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	Scrap(name string) error
}

// LiveRecordingData is the data for a stored recording
type LiveRecordingData struct {
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

// NewLiveRecordingHandle creates a new stored recording handle
func NewLiveRecordingHandle(name string, s LiveRecording) *LiveRecordingHandle {
	return &LiveRecordingHandle{
		name: name,
		s:    s,
	}
}

// A LiveRecordingHandle is a reference to a stored recording that can be operated on
type LiveRecordingHandle struct {
	name string
	s    LiveRecording
}

// ID returns the identifier of the live recording
func (s *LiveRecordingHandle) ID() string {
	return s.name
}

// Data gets the data for the stored recording
func (s *LiveRecordingHandle) Data() (d LiveRecordingData, err error) {
	d, err = s.s.Data(s.name)
	return
}

// Scrap stops and deletes the recording
func (s *LiveRecordingHandle) Scrap() (err error) {
	err = s.s.Scrap(s.name)
	return
}

// Delete deletes the recording
func (s *LiveRecordingHandle) Delete() (err error) {
	err = s.s.Delete(s.name)
	return
}

// Resume resumes the recording
func (s *LiveRecordingHandle) Resume() (err error) {
	err = s.s.Resume(s.name)
	return
}

// Pause pauses the recording
func (s *LiveRecordingHandle) Pause() (err error) {
	err = s.s.Pause(s.name)
	return
}

// Mute mutes the recording
func (s *LiveRecordingHandle) Mute() (err error) {
	err = s.s.Mute(s.name)
	return
}

// Unmute mutes the recording
func (s *LiveRecordingHandle) Unmute() (err error) {
	err = s.s.Unmute(s.name)
	return
}

// Match returns true if the event matches the bridge
func (s *LiveRecordingHandle) Match(e Event) bool {
	r, ok := e.(RecordingEvent)
	if !ok {
		return false
	}
	rIDs := r.GetRecordingIDs()
	for _, i := range rIDs {
		if i == s.ID() {
			return true
		}
	}
	return false
}
