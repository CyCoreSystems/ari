package native

import "github.com/CyCoreSystems/ari"

type nativeStoredRecording struct {
	conn *Conn
}

func (sr *nativeStoredRecording) List() (sx []*ari.StoredRecordingHandle, err error) {
	var recs []struct {
		Name string `json:"name"`
	}

	err = Get(sr.conn, "/recordings/stored", &recs)
	for _, rec := range recs {
		sx = append(sx, sr.Get(rec.Name))
	}

	return
}

func (sr *nativeStoredRecording) Get(name string) (s *ari.StoredRecordingHandle) {
	s = ari.NewStoredRecordingHandle(name, sr)
	return
}

func (sr *nativeStoredRecording) Data(name string) (d ari.StoredRecordingData, err error) {
	err = Get(sr.conn, "/recordings/stored/"+name, &d)
	return
}

func (sr *nativeStoredRecording) Copy(name string, dest string) (h *ari.StoredRecordingHandle, err error) {

	var resp struct {
		Name string `json:"name"`
	}

	var request struct {
		Dest string `json:"destinationRecordingName"`
	}

	request.Dest = dest

	err = Post(sr.conn, "/recordings/stored/"+name+"/copy", &resp, &request)

	if err != nil {
		return nil, err
	}

	return sr.Get(resp.Name), nil
}

func (sr *nativeStoredRecording) Delete(name string) (err error) {
	err = Delete(sr.conn, "/recordings/stored/"+name+"", nil, "")
	return
}
