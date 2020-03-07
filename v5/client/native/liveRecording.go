package native

import (
	"errors"

	"github.com/CyCoreSystems/ari/v5"
)

// LiveRecording provides the ARI LiveRecording accessors for the native client
type LiveRecording struct {
	client *Client
}

// Get gets a lazy handle for the live recording name
func (lr *LiveRecording) Get(key *ari.Key) (h *ari.LiveRecordingHandle) {
	h = ari.NewLiveRecordingHandle(lr.client.stamp(key), lr, nil)
	return
}

// Data retrieves the state of the live recording
func (lr *LiveRecording) Data(key *ari.Key) (d *ari.LiveRecordingData, err error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("liveRecording key not supplied")
	}

	data := new(ari.LiveRecordingData)
	if err := lr.client.get("/recordings/live/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "liveRecording", "%v", key.ID)
	}

	data.Key = lr.client.stamp(key)

	return data, nil
}

// Stop stops the live recording.
//
// NOTE: if the recording is already stopped, this will return an error.
func (lr *LiveRecording) Stop(key *ari.Key) error {
	if key == nil || key.ID == "" {
		return errors.New("liveRecording key not supplied")
	}

	return lr.client.post("/recordings/live/"+key.ID+"/stop", nil, nil)
}

// Pause pauses the live recording (TODO: does it error if the live recording is already paused)
func (lr *LiveRecording) Pause(key *ari.Key) error {
	return lr.client.post("/recordings/live/"+key.ID+"/pause", nil, nil)
}

// Resume resumes the live recording (TODO: does it error if the live recording is already resumed)
func (lr *LiveRecording) Resume(key *ari.Key) error {
	return lr.client.del("/recordings/live/"+key.ID+"/pause", nil, "")
}

// Mute mutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Mute(key *ari.Key) error {
	return lr.client.post("/recordings/live/"+key.ID+"/mute", nil, nil)
}

// Unmute unmutes the live recording (TODO: does it error if the live recording is already muted)
func (lr *LiveRecording) Unmute(key *ari.Key) error {
	return lr.client.del("/recordings/live/"+key.ID+"/mute", nil, "")
}

// Scrap removes a live recording (TODO: describe difference between scrap and delete)
func (lr *LiveRecording) Scrap(key *ari.Key) error {
	return lr.client.del("/recordings/live/"+key.ID, nil, "")
}

// Stored returns the StoredRecording handle for the given LiveRecording
func (lr *LiveRecording) Stored(key *ari.Key) *ari.StoredRecordingHandle {
	return ari.NewStoredRecordingHandle(
		lr.client.stamp(key.New(ari.StoredRecordingKey, key.ID)),
		lr.client.StoredRecording(),
		nil,
	)
}

// Subscribe is a shim to enable recording handles to subscribe to their
// underlying bridge/channel for events.  It should not be called directly.
func (lr *LiveRecording) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	return lr.client.Bus().Subscribe(key, n...)
}
