package nc

import "github.com/CyCoreSystems/ari"

type natsStoredRecording struct {
	conn *Conn
}

// tests and other advanced utility functions can cast to an interface to get the NatsConnection object out
func (sr *natsStoredRecording) NatsConnection() *Conn {
	return sr.conn
}

func (sr *natsStoredRecording) List() (sx []*ari.StoredRecordingHandle, err error) {
	var recordings []string
	err = sr.conn.readRequest("ari.recording.stored.all", nil, &recordings)
	for _, r := range recordings {
		sx = append(sx, sr.Get(r))
	}

	return
}

func (sr *natsStoredRecording) Get(name string) *ari.StoredRecordingHandle {
	return ari.NewStoredRecordingHandle(name, sr)
}

func (sr *natsStoredRecording) Data(name string) (srd ari.StoredRecordingData, err error) {
	err = sr.conn.readRequest("ari.recording.stored.data."+name, nil, &srd)
	return
}

func (sr *natsStoredRecording) Copy(name string, dest string) (h *ari.StoredRecordingHandle, err error) {
	err = sr.conn.standardRequest("ari.recording.stored.copy."+name, &dest, nil)
	if err != nil {
		return
	}
	h = sr.Get(dest) //TODO: confirm dest is ID of the new copy. Should be.
	return
}

func (sr *natsStoredRecording) Delete(name string) (err error) {
	err = sr.conn.standardRequest("ari.recording.stored.delete."+name, nil, nil)
	return
}
