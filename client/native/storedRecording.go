package native

import "github.com/CyCoreSystems/ari"

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
		filter = ari.NodeKey(sr.client.ApplicationName(), sr.client.node)
	}

	err = sr.client.get("/recordings/stored", &recs)
	for _, rec := range recs {
		k := ari.NewKey(ari.StoredRecordingKey, rec.Name, ari.WithNode(sr.client.node), ari.WithApp(sr.client.ApplicationName()))
		if filter.Match(k) {
			sx = append(sx, k)
		}
	}

	return
}

// Get gets a lazy handle for the given stored recording name
func (sr *StoredRecording) Get(key *ari.Key) (s ari.StoredRecordingHandle) {
	s = NewStoredRecordingHandle(key, sr, nil)
	return
}

// Data retrieves the state of the stored recording
func (sr *StoredRecording) Data(key *ari.Key) (d *ari.StoredRecordingData, err error) {
	d = &ari.StoredRecordingData{}
	name := key.ID
	err = sr.client.get("/recordings/stored/"+name, d)
	if err != nil {
		d = nil
		err = dataGetError(err, "storedRecording", "%v", name)
		return
	}
	return
}

// Copy copies a stored recording and returns the new handle
func (sr *StoredRecording) Copy(key *ari.Key, dest string) (h ari.StoredRecordingHandle, err error) {
	h = sr.StageCopy(key, dest)
	err = h.Exec()
	return
}

// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
func (sr *StoredRecording) StageCopy(key *ari.Key, dest string) (h ari.StoredRecordingHandle) {

	var resp struct {
		Name string `json:"name"`
	}

	var request struct {
		Dest string `json:"destinationRecordingName"`
	}

	request.Dest = dest

	name := key.ID

	destKey := ari.NewKey(ari.StoredRecordingKey, dest, ari.WithNode(sr.client.node), ari.WithApp(sr.client.ApplicationName()))
	return NewStoredRecordingHandle(destKey, sr, func(h *StoredRecordingHandle) (err error) {
		err = sr.client.post("/recordings/stored/"+name+"/copy", &resp, &request)
		return
	})
}

// Delete deletes the stored recording
func (sr *StoredRecording) Delete(key *ari.Key) (err error) {
	name := key.ID
	err = sr.client.del("/recordings/stored/"+name+"", nil, "")
	return
}

// A StoredRecordingHandle is a reference to a stored recording that can be operated on
type StoredRecordingHandle struct {
	key      *ari.Key
	s        *StoredRecording
	exec     func(a *StoredRecordingHandle) error
	executed bool
}

// NewStoredRecordingHandle creates a new stored recording handle
func NewStoredRecordingHandle(key *ari.Key, s *StoredRecording, exec func(a *StoredRecordingHandle) error) ari.StoredRecordingHandle {
	return &StoredRecordingHandle{
		key:  key,
		s:    s,
		exec: exec,
	}
}

// ID returns the identifier for the stored recording
func (s *StoredRecordingHandle) ID() string {
	return s.key.ID
}

// Exec executes any staged operations
func (s *StoredRecordingHandle) Exec() (err error) {
	if !s.executed {
		s.executed = true
		if s.exec != nil {
			err = s.exec(s)
			s.exec = nil
		}
	}
	return
}

// Data gets the data for the stored recording
func (s *StoredRecordingHandle) Data() (d *ari.StoredRecordingData, err error) {
	d, err = s.s.Data(s.key)
	return
}

// Copy copies the stored recording
func (s *StoredRecordingHandle) Copy(dest string) (h ari.StoredRecordingHandle, err error) {
	h, err = s.s.Copy(s.key, dest)
	return
}

// StageCopy creates a `StoredRecordingHandle` with a `Copy` operation staged.
func (s *StoredRecordingHandle) StageCopy(dest string) (h ari.StoredRecordingHandle) {
	h = s.s.StageCopy(s.key, dest)
	return
}

// Delete deletes the recording
func (s *StoredRecordingHandle) Delete() (err error) {
	err = s.s.Delete(s.key)
	return
}
