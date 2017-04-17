package native

import "github.com/CyCoreSystems/ari"

// StoredRecording provides the ARI StoredRecording accessors for the native client
type StoredRecording struct {
	client *Client
}

// List lists the current stored recordings and returns a list of handles
func (sr *StoredRecording) List() (sx []ari.StoredRecordingHandle, err error) {
	var recs []struct {
		Name string `json:"name"`
	}

	err = sr.client.get("/recordings/stored", &recs)
	for _, rec := range recs {
		sx = append(sx, sr.Get(rec.Name))
	}

	return
}

// Get gets a lazy handle for the given stored recording name
func (sr *StoredRecording) Get(name string) (s ari.StoredRecordingHandle) {
	s = NewStoredRecordingHandle(name, sr, nil)
	return
}

// Data retrieves the state of the stored recording
func (sr *StoredRecording) Data(name string) (d *ari.StoredRecordingData, err error) {
	d = &ari.StoredRecordingData{}
	err = sr.client.get("/recordings/stored/"+name, d)
	if err != nil {
		d = nil
		err = dataGetError(err, "storedRecording", "%v", name)
		return
	}
	return
}

// Copy copies a stored recording and returns the new handle
func (sr *StoredRecording) Copy(name string, dest string) (h ari.StoredRecordingHandle, err error) {
	h = sr.StageCopy(name, dest)
	err = h.Exec()
	return
}

// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
func (sr *StoredRecording) StageCopy(name string, dest string) (h ari.StoredRecordingHandle) {

	var resp struct {
		Name string `json:"name"`
	}

	var request struct {
		Dest string `json:"destinationRecordingName"`
	}

	request.Dest = dest

	return NewStoredRecordingHandle(name, sr, func(h *StoredRecordingHandle) (err error) {
		err = sr.client.post("/recordings/stored/"+name+"/copy", &resp, &request)
		return
	})
}

// Delete deletes the stored recording
func (sr *StoredRecording) Delete(name string) (err error) {
	err = sr.client.del("/recordings/stored/"+name+"", nil, "")
	return
}

// A StoredRecordingHandle is a reference to a stored recording that can be operated on
type StoredRecordingHandle struct {
	name string
	s    *StoredRecording
	exec func(a *StoredRecordingHandle) error
}

// NewStoredRecordingHandle creates a new stored recording handle
func NewStoredRecordingHandle(name string, s *StoredRecording, exec func(a *StoredRecordingHandle) error) ari.StoredRecordingHandle {
	return &StoredRecordingHandle{
		name: name,
		s:    s,
		exec: exec,
	}
}

// ID returns the identifier for the stored recording
func (s *StoredRecordingHandle) ID() string {
	return s.name
}

// Exec executes any staged operations
func (s *StoredRecordingHandle) Exec() (err error) {
	if s.exec != nil {
		err = s.exec(s)
		s.exec = nil
	}
	return
}

// Data gets the data for the stored recording
func (s *StoredRecordingHandle) Data() (d *ari.StoredRecordingData, err error) {
	d, err = s.s.Data(s.name)
	return
}

// Copy copies the stored recording
func (s *StoredRecordingHandle) Copy(dest string) (h ari.StoredRecordingHandle, err error) {
	h, err = s.s.Copy(s.name, dest)
	return
}

// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
func (s *StoredRecordingHandle) StageCopy(dest string) (h ari.StoredRecordingHandle) {
	h = s.s.StageCopy(s.name, dest)
	return
}

// Delete deletes the recording
func (s *StoredRecordingHandle) Delete() (err error) {
	err = s.s.Delete(s.name)
	return
}
