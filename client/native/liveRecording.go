package native

import "github.com/CyCoreSystems/ari"

// LiveRecording provides the ARI LiveRecording accessors for the native client
type LiveRecording struct {
	client *Client
}

// Get gets a lazy handle for the live recording name
func (lr *LiveRecording) Get(name string) (h ari.LiveRecordingHandle) {
	h = NewLiveRecordingHandle(name, lr, nil)
	return
}

// Data retrieves the state of the live recording
func (lr *LiveRecording) Data(name string) (d *ari.LiveRecordingData, err error) {
	d = &ari.LiveRecordingData{}
	err = lr.client.get("/recordings/live/"+name, &d)
	if err != nil {
		d = nil
		err = dataGetError(err, "liveRecording", "%v", name)
		return
	}
	return
}

// Stop stops the live recording (TODO: does it error if the live recording is already stopped)
func (lr *LiveRecording) Stop(name string) (err error) {
	err = lr.client.post("/recordings/live/"+name+"/stop", nil, nil)
	return
}

// Pause pauses the live recording (TODO: does it error if the live recording is already paused)
func (lr *LiveRecording) Pause(name string) (err error) {
	err = lr.client.post("/recordings/live/"+name+"/pause", nil, nil)
	return
}

// Resume resumes the live recording (TODO: does it error if the live recording is already resumed)
func (lr *LiveRecording) Resume(name string) (err error) {
	err = lr.client.del("/recordings/live/"+name+"/pause", nil, "")
	return
}

// Mute mutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Mute(name string) (err error) {
	err = lr.client.post("/recordings/live/"+name+"/mute", nil, nil)
	return
}

// Unmute unmutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Unmute(name string) (err error) {
	err = lr.client.del("/recordings/live/"+name+"/mute", nil, "")
	return
}

// Delete deletes the live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Delete(name string) (err error) {
	//NOTE: original code used 'stored' for this even though it's live
	err = lr.client.del("/recordings/stored/"+name, nil, "")
	return
}

// Scrap removes a live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Scrap(name string) (err error) {
	//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	err = lr.client.del("/recordings/live/"+name, nil, "")
	return
}

// NewLiveRecordingHandle creates a new stored recording handle
func NewLiveRecordingHandle(name string, s *LiveRecording, exec func() (err error)) ari.LiveRecordingHandle {
	return &LiveRecordingHandle{
		name: name,
		s:    s,
		exec: exec,
	}
}

// A LiveRecordingHandle is a reference to a stored recording that can be operated on
type LiveRecordingHandle struct {
	name     string
	s        *LiveRecording
	exec     func() (err error)
	executed bool
}

// ID returns the identifier of the live recording
func (s *LiveRecordingHandle) ID() string {
	return s.name
}

// Data gets the data for the stored recording
func (s *LiveRecordingHandle) Data() (d *ari.LiveRecordingData, err error) {
	d, err = s.s.Data(s.name)
	return
}

// Stop stops and saves the recording
func (s *LiveRecordingHandle) Stop() (err error) {
	err = s.s.Stop(s.name)
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
func (s *LiveRecordingHandle) Match(e ari.Event) bool {
	r, ok := e.(ari.RecordingEvent)
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
