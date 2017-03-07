package native

import "github.com/CyCoreSystems/ari"

// LiveRecording provides the ARI LiveRecording accessors for the native client
type LiveRecording struct {
	client *Client
}

// Get gets a lazy handle for the live recording name
func (lr *LiveRecording) Get(name string) (h *ari.LiveRecordingHandle) {
	h = ari.NewLiveRecordingHandle(name, lr)
	return
}

// Data retrieves the state of the live recording
func (lr *LiveRecording) Data(name string) (d *ari.LiveRecordingData, err error) {
	d = &ari.LiveRecordingData{}
	err = lr.client.conn.Get("/recordings/live/"+name, &d)
	if err != nil {
		d = nil
		err = dataGetError(err, "liveRecording", "%v", name)
		return
	}
	return
}

// Stop stops the live recording (TODO: does it error if the live recording is already stopped)
func (lr *LiveRecording) Stop(name string) (err error) {
	err = lr.client.conn.Post("/recordings/live/"+name+"/stop", nil, nil)
	return
}

// Pause pauses the live recording (TODO: does it error if the live recording is already paused)
func (lr *LiveRecording) Pause(name string) (err error) {
	err = lr.client.conn.Post("/recordings/live/"+name+"/pause", nil, nil)
	return
}

// Resume resumes the live recording (TODO: does it error if the live recording is already resumed)
func (lr *LiveRecording) Resume(name string) (err error) {
	err = lr.client.conn.Delete("/recordings/live/"+name+"/pause", nil, "")
	return
}

// Mute mutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Mute(name string) (err error) {
	err = lr.client.conn.Post("/recordings/live/"+name+"/mute", nil, nil)
	return
}

// Unmute unmutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Unmute(name string) (err error) {
	err = lr.client.conn.Delete("/recordings/live/"+name+"/mute", nil, "")
	return
}

// Delete deletes the live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Delete(name string) (err error) {
	//NOTE: original code used 'stored' for this even though it's live
	err = lr.client.conn.Delete("/recordings/stored/"+name, nil, "")
	return
}

// Scrap removes a live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Scrap(name string) (err error) {
	//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	err = lr.client.conn.Delete("/recordings/live/"+name, nil, "")
	return
}
