package ari

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// LiveRecording describes a recording which is in progress
type LiveRecording struct {
	Cause       string `json:"cause,omitempty"`            // If failed, the cause of the failure
	DurationSec int    `json:"duration,omitempty"`         // Length of recording in seconds
	Format      string `json:"format"`                     // Format of recording (wav, gsm, etc)
	Name        string `json:"name"`                       // (base) name for the recording
	SilenceSec  int    `json:"silence_duration,omitempty"` // If silence was detected in the recording, the duration in seconds of that silence (requires that maxSilenceSeconds be non-zero)
	State       string `json:"state"`                      // Current state of the recording
	TalkingSec  int    `json:"talking_duration,omitempty"` // Duration of talking, in seconds, that has been detected in the recording (requires that maxSilenceSeconds be non-zero)
	TargetURI   string `json:"target_uri"`                 // URI for the channel or bridge which is being recorded (TODO: figure out format for this)

	client *Client // Reference to the client which created or returned this LiveRecording

	doneChan chan struct{} // channel for indicating the the recording is stopped.

	status int // The status of the live recording
}

// Duration returns the duration of the live recording, if known
func (l *LiveRecording) Duration() (dur time.Duration) {
	if l.DurationSec > 0 {
		dur = time.Duration(l.DurationSec) * time.Second
	}
	return
}

// SilenceDuration returns the duration of the detected silence
// during the recording, if known
func (l *LiveRecording) SilenceDuration() (dur time.Duration) {
	if l.SilenceSec > 0 {
		dur = time.Duration(l.SilenceSec) * time.Second
	}
	return
}

// TalkingDuration returns the duration of the detected talking
// during the recording, if known
func (l *LiveRecording) TalkingDuration() (dur time.Duration) {
	if l.TalkingSec > 0 {
		dur = time.Duration(l.TalkingSec) * time.Second
	}
	return
}

var (
	// ExistsFail indicates that a recording should fail if the
	// given name already exists.
	ExistsFail = "fail"

	// ExistsOverwrite indicates that if a recording exists of
	// the same name, it should be overwritten.
	ExistsOverwrite = "overwrite"

	// ExistsAppend indicates that if a recording exists of the
	// same name, it should be appended to.
	ExistsAppend = "append"
)

var (
	// TerminateNever indicates that a recording should not be
	// ended on any DTMF tone.
	TerminateNever = "none"

	// TerminateAny indicates that a recording should be terminated
	// if any DTMF digit is received
	TerminateAny = "any"

	// TerminateStar indicates that a recording should be terminated
	// if a * DTMF character is received.
	TerminateStar = "*"

	// TerminateHash indicates that a recording should be terminated
	// if a # DTMF character is received.
	TerminateHash = "#"

	// TerminatePound indicates that a recording should be terminated
	// if a # DTMF character is received.
	TerminatePound = "#"
)

const (
	// RecordInProgress indicates that a recording is still in progress
	RecordInProgress = iota

	// RecordCanceled indicates that a recording was canceled (by request)
	RecordCanceled

	// RecordFailed indicates that a recording failed
	RecordFailed

	// RecordFinished indicates that a recording finished normally
	RecordFinished

	// RecordHangup indicates that a recording was ended due to hangup
	RecordHangup
)

// RecordingOptions describes the set of options available when making a recording.
type RecordingOptions struct {
	// Format is the file format/encoding to which the recording should be stored.
	// This will usually be one of: slin, ulaw, alaw, wav, gsm.
	// If not specified, this will default to slin.
	Format string

	// MaxDuration is the maximum duration of the recording, after which the recording will
	// automatically stop.  If not set, there is no maximum.
	MaxDuration time.Duration

	// MaxSilence is the maximum duration of detected to be found before terminating the recording.
	MaxSilence time.Duration

	// Exists determines what should happen if the given recording already exists.
	// Valid values are: "fail", "overwrite", or "append".
	// If not specified, it will default to "fail"
	Exists string

	// Beep indicates whether a beep should be played to the recorded
	// party at the beginning of the recording.
	Beep bool

	// Terminate indicates whether the recording should be terminated on
	// receipt of a DTMF digit.
	// valid options are: "none", "any", "*", and "#"
	// If not specified, it will default to "none" (never terminate on DTMF).
	Terminate string
}

// ToRequest converts a set of recording options to a
// record request.
func (o *RecordingOptions) ToRequest(name string) *RecordRequest {
	if o.Format == "" {
		o.Format = "slin"
	}
	return &RecordRequest{
		Name:               name,
		Format:             o.Format,
		MaxDurationSeconds: int(o.MaxDuration.Seconds()),
		MaxSilenceSeconds:  int(o.MaxSilence.Seconds()),
		IfExists:           o.Exists,
		Beep:               o.Beep,
		TerminateOn:        o.Terminate,
	}
}

// StoredRecording describes a past recording which may be played back (via GetStoredRecording)
type StoredRecording struct {
	Format string `json:"format"`
	Name   string `json:"name"`

	client *Client // Reference to the client which created or returned this StoredRecording
}

// A Recorder is anything which can "Record"
type Recorder interface {
	Record(string, *RecordingOptions) (*LiveRecording, error)
	GetClient() *Client
}

// A PlayerRecorder is anywhich which can Record and Play
type PlayerRecorder interface {
	Record(string, *RecordingOptions) (*LiveRecording, error)
	Play(string) (string, error)
	GetClient() *Client
}

//ListStoredRecordings lists all completed recordings
//Equivalent to GET /recordings/stored
func (c *Client) ListStoredRecordings() ([]StoredRecording, error) {
	var m []StoredRecording
	err := c.Get("/recordings/stored", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//GetStoredRecording returns a stored recording's details
//Equivalent to GET /recordings/stored/{recordingName}
func (c *Client) GetStoredRecording(recordingName string) (StoredRecording, error) {
	var m StoredRecording
	err := c.Get("/recordings/stored/"+recordingName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//GetLiveRecording returns a specific live recording
//Equivalent to GET /recordings/live/{recordingName}
func (c *Client) GetLiveRecording(recordingName string) (*LiveRecording, error) {
	var m LiveRecording
	err := c.Get("/recordings/live/"+recordingName, &m)
	return &m, err
}

// Record starts a recording on the given Recorder.
func Record(ctx context.Context, r Recorder, name string, opts *RecordingOptions) (*LiveRecording, error) {
	// Extract the ari client from the recorder
	c := r.GetClient()
	if c == nil {
		return nil, fmt.Errorf("Failed to find *ari.Client in Recorder")
	}

	// Listen for start, stop, and failed events
	startSub := c.Bus.Subscribe("RecordingStarted")
	defer startSub.Cancel()

	failedSub := c.Bus.Subscribe("RecordingFailed")
	defer failedSub.Cancel()

	finishedSub := c.Bus.Subscribe("RecordingFinished")
	defer finishedSub.Cancel()

	// Start recording
	rec, err := r.Record(name, opts)
	if err != nil {
		return nil, err
	}

	// TODO: we have no way to track hangups because we do
	// not have the affiliated channel ID.  We _may_ be able
	// to compare a ChannelHangupRequest event's channel with
	// the LiveRecording's TargetURI, but that will only work
	// for channels.

	// Wait for the recording to start
	startTimer := time.NewTimer(RecordingStartTimeout)
	for {
		select {
		case <-startTimer.C:
			rec.setStatus(RecordFailed)
			return rec, fmt.Errorf("Timed out waiting for recording to start.")
		case <-ctx.Done():
			rec.setStatus(RecordCanceled)
			return rec, fmt.Errorf("Recording canceled.")
		case e := <-startSub.C:
			r := e.(*RecordingStarted).Recording
			if r.Name == name {
				Logger.Debug("Recording started.")
				startTimer.Stop()
				rec = &r
				rec.client = c
			}
		case e := <-failedSub.C:
			r := e.(*RecordingFailed).Recording
			if r.Name == name {
				rec.setStatus(RecordFailed)
				return rec, fmt.Errorf("Recording failed: %s", r.Cause)
			}
		case e := <-finishedSub.C:
			r := e.(*RecordingFinished).Recording
			if r.Name == name {
				Logger.Debug("Recording stopped")
				rec = &r
				rec.client = c
				rec.setStatus(RecordFinished)
				return rec, nil
			}
		}
	}
}

//Copy current StoredRecording to a new name (retaining the existing copy)
func (s *StoredRecording) Copy(destination string) (StoredRecording, error) {
	var sRet StoredRecording
	if s.client == nil {
		return sRet, fmt.Errorf("No client found in StoredRecording")
	}
	return s.client.CopyStoredRecording(s.Name, destination)
}

func (l *LiveRecording) setStatus(status int) {
	l.status = status
}

// Status returns the status of the recording.
func (l *LiveRecording) Status() int {
	return l.status
}

//Stop and store current LiveRecording
func (l *LiveRecording) Stop() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.StopLiveRecording(l.Name)
}

//Pause current LiveRecording
func (l *LiveRecording) Pause() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.PauseLiveRecording(l.Name)
}

//Mute current LiveRecording
func (l *LiveRecording) Mute() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.MuteLiveRecording(l.Name)
}

//Delete current LiveRecording
func (l *LiveRecording) Delete() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.DeleteStoredRecording(l.Name)
}

// Scrap Stops and deletes the current LiveRecording
//TODO: reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
func (l *LiveRecording) Scrap() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.ScrapLiveRecording(l.Name)
}

// Resume unpauses the current LiveRecording
func (l *LiveRecording) Resume() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.ResumeLiveRecording(l.Name)
}

// Unmute current LiveRecording
func (l *LiveRecording) Unmute() error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.UnmuteLiveRecording(l.Name)
}

// CopyStoredRecording copies a stored recording to a new destination
// within the stored recordings tree.
//Equivalent to Post /recordings/stored/{recordingName}/copy
func (c *Client) CopyStoredRecording(recordingName string, destination string) (StoredRecording, error) {
	var m StoredRecording

	//Request structure to copy a stored recording. DestinationRecordingName is required.
	type request struct {
		DestinationRecordingName string `json:"destinationRecordingName"`
	}

	req := request{destination}

	//Make the request
	err := c.Post("/recordings/stored/"+recordingName+"/copy", &m, &req)
	//TODO add individual error handling

	if err != nil {
		return m, err
	}
	return m, nil
}

// StopLiveRecording stops and stores a live recording
//Equivalent to Post /recordings/live/{recordingName}/stop
func (c *Client) StopLiveRecording(recordingName string) error {
	err := c.Post("/recordings/live/"+recordingName+"/stop", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//PauseLiveRecording pauses a live recording
//Equivalent to Post /recordings/live/{recordingName}/pause
func (c *Client) PauseLiveRecording(recordingName string) error {

	//Since no request body is required nor return object
	//we just pass two nils.

	err := c.Post("/recordings/live/"+recordingName+"/pause", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//MuteLiveRecording mutes a live recording
//Equivalent to Post /recordings/live/{recordingName}/mute
func (c *Client) MuteLiveRecording(recordingName string) error {
	err := c.Post("/recordings/live/"+recordingName+"/mute", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//DeleteStoredRecording deletes a stored recording
//Equivalent to DELETE /recordings/stored/{recordingName}
func (c *Client) DeleteStoredRecording(recordingName string) error {
	err := c.Delete("/recordings/stored/"+recordingName, nil, "")
	return err
}

// ScrapLiveRecording stops a live recording and discard it
//Equivalent to DELETE /recordings/live/{recordingName}
//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
func (c *Client) ScrapLiveRecording(recordingName string) error {
	err := c.Delete("/recordings/live/"+recordingName, nil, "")
	return err
}

// ResumeLiveRecording resumes (unpauses) a live recording
//Equivalent to DELETE /recordings/live/{recordingName}/pause
func (c *Client) ResumeLiveRecording(recordingName string) error {
	err := c.Delete("/recordings/live/"+recordingName+"/pause", nil, "")
	return err
}

// UnmuteLiveRecording unmutes a live recording
//Equivalent to DELETE /recordings/live/{recordingName}/mute
func (c *Client) UnmuteLiveRecording(recordingName string) error {
	err := c.Delete("/recordings/live/"+recordingName+"/mute", nil, "")
	return err
}
