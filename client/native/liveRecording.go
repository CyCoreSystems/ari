package native

import "github.com/CyCoreSystems/ari"

// LiveRecording provides the ARI LiveRecording accessors for the native client
type LiveRecording struct {
	client *Client
}

// Get gets a lazy handle for the live recording name
func (lr *LiveRecording) Get(key *ari.Key) (h *ari.LiveRecordingHandle) {
	h = ari.NewLiveRecordingHandle(key, lr, nil)
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
