package native

import "github.com/CyCoreSystems/ari"

// StoredRecording provides the ARI StoredRecording accessors for the native client
type StoredRecording struct {
	client *Client
}

// List lists the current stored recordings and returns a list of handles
func (sr *StoredRecording) List() (sx []*ari.StoredRecordingHandle, err error) {
	var recs []struct {
		Name string `json:"name"`
	}

	err = sr.client.conn.Get("/recordings/stored", &recs)
	for _, rec := range recs {
		sx = append(sx, sr.Get(rec.Name))
	}

	return
}

// Get gets a lazy handle for the given stored recording name
func (sr *StoredRecording) Get(name string) (s *ari.StoredRecordingHandle) {
	s = ari.NewStoredRecordingHandle(name, sr)
	return
}

// Data retrieves the state of the stored recording
func (sr *StoredRecording) Data(name string) (d ari.StoredRecordingData, err error) {
	err = sr.client.conn.Get("/recordings/stored/"+name, &d)
	return
}

// Copy copies a stored recording and returns the new handle
func (sr *StoredRecording) Copy(name string, dest string) (h *ari.StoredRecordingHandle, err error) {

	var resp struct {
		Name string `json:"name"`
	}

	var request struct {
		Dest string `json:"destinationRecordingName"`
	}

	request.Dest = dest

	err = sr.client.conn.Post("/recordings/stored/"+name+"/copy", &resp, &request)

	if err != nil {
		return nil, err
	}

	return sr.Get(resp.Name), nil
}

// Delete deletes the stored recording
func (sr *StoredRecording) Delete(name string) (err error) {
	err = sr.client.conn.Delete("/recordings/stored/"+name+"", nil, "")
	return
}
