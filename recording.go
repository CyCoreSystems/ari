package ari

import "fmt"

// LiveRecording describes a recording which is in progress
type LiveRecording struct {
	Cause            string  `json:"cause,omitempty"`            // If failed, the cause of the failure
	Duration         int     `json:"duration,omitempty"`         // Length of recording in seconds
	Format           string  `json:"format"`                     // Format of recording (wav, gsm, etc)
	Name             string  `json:"name"`                       // (base) name for the recording
	Silence_duration int     `json:"silence_duration,omitempty"` // If silence was detected in the recording, the duration in seconds of that silence (requires that maxSilenceSeconds be non-zero)
	State            string  `json:"state"`                      // Current state of the recording
	Talking_duration int     `json:"talking_duration,omitempty"` // Duration of talking, in seconds, that has been detected in the recording (requires that maxSilenceSeconds be non-zero)
	Target_uri       string  `json:"target_uri"`                 // URI for the channel or bridge which is being recorded (TODO: figure out format for this)
	client           *Client // Reference to the client which created or returned this LiveRecording
}

// StoredRecording describes a past recording which may be played back (via GetStoredRecording)
type StoredRecording struct {
	Format string `json:"format"`
	Name   string `json:"name"`

	client *Client // Reference to the client which created or returned this StoredRecording
}

//List all completed recordings
//Equivalent to GET /recordings/stored
func (c *Client) ListStoredRecordings() ([]StoredRecording, error) {
	var m []StoredRecording
	err := c.AriGet("/recordings/stored", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Get a stored recording's details
//Equivalent to GET /recordings/stored/{recordingName}
func (c *Client) GetStoredRecording(recordingName string) (StoredRecording, error) {
	var m StoredRecording
	err := c.AriGet("/recordings/stored/"+recordingName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Get a specific live recording
//Equivalent to GET /recordings/live/{recordingName}
func (c *Client) GetLiveRecording(recordingName string) (LiveRecording, error) {
	var m LiveRecording
	err := c.AriGet("/recordings/live/"+recordingName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Copy current StoredRecording
func (s *StoredRecording) Copy(recordingName string, destination string) (StoredRecording, error) {
	if s.client == nil {
		return fmt.Errorf("No client found in StoredRecording")
	}
	return s.client.CopyLiveRecording(recordingName, destination)
}

// No method for getting the current LiveRecording--you have it.

//Stop and store current LiveRecording
func (l *LiveRecording) Stop(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.StopLiveRecording(recordingName)
}

//Pause current LiveRecording
func (l *LiveRecording) Pause(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.PauseLiveRecording(recordingName)
}

//Mute current LiveRecording
func (l *LiveRecording) Mute(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.MuteLiveRecording(recordingName)
}

//Delete current LiveRecording
func (l *LiveRecording) Delete(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.DeleteStoredRecording(recordingName)
}

//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
//Stop and delete current LiveRecording
func (l *LiveRecording) Scrap(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.ScrapLiveRecording(recordingName)
}

//Unpause current LiveRecording
func (l *LiveRecording) Resume(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.ResumeLiveRecording(recordingName)
}

//Unmute current LiveRecording
func (l *LiveRecording) Unmute(recordingName string) error {
	if l.client == nil {
		return fmt.Errorf("No client found in LiveRecording")
	}
	return l.client.UnmuteLiveRecording(recordingName)
}

//Copy a stored recording
//Equivalent to Post /recordings/stored/{recordingName}/copy
func (c *Client) CopyStoredRecording(recordingName string, destination string) (StoredRecording, error) {
	var m StoredRecording

	//Request structure to copy a stored recording. DestinationRecordingName is required.
	type request struct {
		DestinationRecordingName string `json:"destinationRecordingName"`
	}

	req := request{destination}

	//Make the request
	err := c.AriPost("/recordings/stored/"+recordingName+"/copy", &m, &req)
	//TODO add individual error handling

	if err != nil {
		return m, err
	}
	return m, nil
}

//Stop and store a live recording
//Equivalent to Post /recordings/live/{recordingName}/stop
func (c *Client) StopLiveRecording(recordingName string) error {
	err := c.AriPost("/recordings/live/"+recordingName+"/stop", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Pause a live recording
//Equivalent to Post /recordings/live/{recordingName}/pause
func (c *Client) PauseLiveRecording(recordingName string) error {

	//Since no request body is required nor return object
	//we just pass two nils.

	err := c.AriPost("/recordings/live/"+recordingName+"/pause", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Mute a live recording
//Equivalent to Post /recordings/live/{recordingName}/mute
func (c *Client) MuteLiveRecording(recordingName string) error {
	err := c.AriPost("/recordings/live/"+recordingName+"/mute", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Delete a stored recording
//Equivalent to DELETE /recordings/stored/{recordingName}
func (c *Client) DeleteStoredRecording(recordingName string) error {
	err := c.AriDelete("/recordings/stored/"+recordingName, nil, nil)
	return err
}

//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
//Stop a live recording and discard it
//Equivalent to DELETE /recordings/live/{recordingName}
func (c *Client) ScrapLiveRecording(recordingName string) error {
	err := c.AriDelete("/recordings/live/"+recordingName, nil, nil)
	return err
}

//Unpause a live recording
//Equivalent to DELETE /recordings/live/{recordingName}/pause
func (c *Client) ResumeLiveRecording(recordingName string) error {
	err := c.AriDelete("/recordings/live/"+recordingName+"/pause", nil, nil)
	return err
}

//Unmute a live recording
//Equivalent to DELETE /recordings/live/{recordingName}/mute
func (c *Client) UnmuteLiveRecording(recordingName string) error {
	err := c.AriDelete("/recordings/live/"+recordingName+"/mute", nil, nil)
	return err
}
