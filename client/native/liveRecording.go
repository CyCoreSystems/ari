package native

import "github.com/CyCoreSystems/ari"

// LiveRecording provides the ARI LiveRecording accessors for the native client
type LiveRecording struct {
	client *Client
}

// Get gets a lazy handle for the live recording name
func (lr *LiveRecording) Get(key *ari.Key) (h ari.LiveRecordingHandle) {
	h = NewLiveRecordingHandle(key, lr, nil)
	return
}

// Data retrieves the state of the live recording
func (lr *LiveRecording) Data(key *ari.Key) (d *ari.LiveRecordingData, err error) {
	d = &ari.LiveRecordingData{}
	name := key.ID
	err = lr.client.get("/recordings/live/"+name, &d)
	if err != nil {
		d = nil
		err = dataGetError(err, "liveRecording", "%v", name)
		return
	}
	return
}

// Stop stops the live recording (TODO: does it error if the live recording is already stopped)
func (lr *LiveRecording) Stop(key *ari.Key) (err error) {
	name := key.ID
	err = lr.client.post("/recordings/live/"+name+"/stop", nil, nil)
	return
}

// Pause pauses the live recording (TODO: does it error if the live recording is already paused)
func (lr *LiveRecording) Pause(key *ari.Key) (err error) {
	name := key.ID
	err = lr.client.post("/recordings/live/"+name+"/pause", nil, nil)
	return
}

// Resume resumes the live recording (TODO: does it error if the live recording is already resumed)
func (lr *LiveRecording) Resume(key *ari.Key) (err error) {
	name := key.ID
	err = lr.client.del("/recordings/live/"+name+"/pause", nil, "")
	return
}

// Mute mutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Mute(key *ari.Key) (err error) {
	name := key.ID
	err = lr.client.post("/recordings/live/"+name+"/mute", nil, nil)
	return
}

// Unmute unmutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Unmute(key *ari.Key) (err error) {
	name := key.ID
	err = lr.client.del("/recordings/live/"+name+"/mute", nil, "")
	return
}

// Delete deletes the live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Delete(key *ari.Key) (err error) {
	//NOTE: original code used 'stored' for this even though it's live
	name := key.ID
	err = lr.client.del("/recordings/stored/"+name, nil, "")
	return
}

// Scrap removes a live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Scrap(key *ari.Key) (err error) {
	//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	name := key.ID
	err = lr.client.del("/recordings/live/"+name, nil, "")
	return
}

// NewLiveRecordingHandle creates a new stored recording handle
func NewLiveRecordingHandle(key *ari.Key, s *LiveRecording, exec func() (err error)) ari.LiveRecordingHandle {
	return &LiveRecordingHandle{
		key:  key,
		s:    s,
		exec: exec,
	}
}

// A LiveRecordingHandle is a reference to a stored recording that can be operated on
type LiveRecordingHandle struct {
	key      *ari.Key
	s        *LiveRecording
	exec     func() (err error)
	executed bool
}

// ID returns the identifier of the live recording
func (s *LiveRecordingHandle) ID() string {
	return s.key.ID
}

// Data gets the data for the stored recording
func (s *LiveRecordingHandle) Data() (d *ari.LiveRecordingData, err error) {
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
