package native

import (
	"errors"

	"github.com/CyCoreSystems/ari/v6"
)

// StoredRecording provides the ARI StoredRecording accessors for the native client
type StoredRecording struct {
	client *Client
}

// List lists the current stored recordings and returns a list of handles
func (sr *StoredRecording) List(filter *ari.Key) (sx []*ari.Key, err error) {
	var recs []struct {
		Name string `json:"name"`
	}

	if filter == nil {
		filter = sr.client.stamp(ari.NewKey(ari.StoredRecordingKey, ""))
	}

	err = sr.client.get("/recordings/stored", &recs)

	for _, rec := range recs {
		k := sr.client.stamp(ari.NewKey(ari.StoredRecordingKey, rec.Name))
		if filter.Match(k) {
			sx = append(sx, k)
		}
	}

	return
}

// Get gets a lazy handle for the given stored recording name
func (sr *StoredRecording) Get(key *ari.Key) *ari.StoredRecordingHandle {
	return ari.NewStoredRecordingHandle(key, sr, nil)
}

// Data retrieves the state of the stored recording
func (sr *StoredRecording) Data(key *ari.Key) (*ari.StoredRecordingData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("storedRecording key not supplied")
	}

	data := new(ari.StoredRecordingData)
	if err := sr.client.get("/recordings/stored/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "storedRecording", "%v", key.ID)
	}

	data.Key = sr.client.stamp(key)

	return data, nil
}

// Copy copies a stored recording and returns the new handle
func (sr *StoredRecording) Copy(key *ari.Key, dest string) (*ari.StoredRecordingHandle, error) {
	h, err := sr.StageCopy(key, dest)
	if err != nil {
		// NOTE: return the handle even on failure so that it can be used to
		//   delete the existing stored recording, should the Copy fail.
		//   ARI provides no facility to force-copy a recording.
		return h, err
	}

	return h, h.Exec()
}

// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
func (sr *StoredRecording) StageCopy(key *ari.Key, dest string) (*ari.StoredRecordingHandle, error) {
	var resp struct {
		Name string `json:"name"`
	}

	req := struct {
		Dest string `json:"destinationRecordingName"`
	}{
		Dest: dest,
	}

	destKey := sr.client.stamp(ari.NewKey(ari.StoredRecordingKey, dest))

	return ari.NewStoredRecordingHandle(destKey, sr, func(h *ari.StoredRecordingHandle) error {
		return sr.client.post("/recordings/stored/"+key.ID+"/copy", &resp, &req)
	}), nil
}

// Delete deletes the stored recording
func (sr *StoredRecording) Delete(key *ari.Key) error {
	return sr.client.del("/recordings/stored/"+key.ID, nil, "")
}
