package ari

// LiveRecording represents a communication path interacting with an Asterisk
// server for live recording resources
type LiveRecording interface {

	// Get gets the Recording by type
	Get(name string) LiveRecordingHandle

	// Data gets the data for the live recording
	Data(name string) (*LiveRecordingData, error)

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

// A LiveRecordingHandle is a reference to a stored recording that can be operated on
type LiveRecordingHandle interface {
	// ID returns the identifier of the live recording
	ID() string

	// Data gets the data for the stored recording
	Data() (d *LiveRecordingData, err error)

	// Stop stops and saves the recording
	Stop() (err error)

	// Scrap stops and deletes the recording
	Scrap() (err error)

	// Delete deletes the recording
	Delete() (err error)

	// Resume resumes the recording
	Resume() (err error)

	// Pause pauses the recording
	Pause() (err error)

	// Mute mutes the recording
	Mute() (err error)

	// Unmute mutes the recording
	Unmute() (err error)

	// Match returns true if the event matches the bridge
	Match(e Event) bool

	// Exec executes any staged operations on the `LiveRecordingHandle`
	Exec() (err error)
}
