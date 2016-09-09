package generic

import "github.com/CyCoreSystems/ari"

type StoredRecording struct {
	Conn Conn
}

func (sr *StoredRecording) List() (sx []*ari.StoredRecordingHandle, err error) {
	var recs []struct {
		Name string `json:"name"`
	}

	err = sr.Conn.Get("/recordings/stored", nil, &recs)
	for _, rec := range recs {
		sx = append(sx, sr.Get(rec.Name))
	}

	return
}

func (sr *StoredRecording) Get(name string) (s *ari.StoredRecordingHandle) {
	s = ari.NewStoredRecordingHandle(name, sr)
	return
}

func (sr *StoredRecording) Data(name string) (d ari.StoredRecordingData, err error) {
	err = sr.Conn.Get("/recordings/stored/%s", []interface{}{name}, &d)
	return
}

func (sr *StoredRecording) Copy(name string, dest string) (h *ari.StoredRecordingHandle, err error) {

	var resp struct {
		Name string `json:"name"`
	}

	var request struct {
		Dest string `json:"destinationRecordingName"`
	}

	request.Dest = dest

	err = sr.Conn.Post("/recordings/stored/%s/copy", []interface{}{name}, &resp, &request)

	if err != nil {
		return nil, err
	}

	return sr.Get(resp.Name), nil
}

func (sr *StoredRecording) Delete(name string) (err error) {
	err = sr.Conn.Delete("/recordings/stored/%s", []interface{}{name}, nil, "")
	return
}
