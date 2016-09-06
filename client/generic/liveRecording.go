package generic

import "github.com/CyCoreSystems/ari"

type LiveRecording struct {
	Conn Conn
}

func (lr *LiveRecording) Get(name string) (h *ari.LiveRecordingHandle) {
	h = ari.NewLiveRecordingHandle(name, lr)
	return
}

func (lr *LiveRecording) Data(name string) (d ari.LiveRecordingData, err error) {
	err = lr.Conn.Get("/recordings/live/%s", []interface{}{name}, &d)
	return
}

func (lr *LiveRecording) Stop(name string) (err error) {
	err = lr.Conn.Post("/recordings/live/%s/stop", []interface{}{name}, nil, nil)
	return
}

func (lr *LiveRecording) Pause(name string) (err error) {
	err = lr.Conn.Post("/recordings/live/%s/pause", []interface{}{name}, nil, nil)
	return
}

func (lr *LiveRecording) Resume(name string) (err error) {
	err = lr.Conn.Delete("/recordings/live/%s/pause", []interface{}{name}, nil, "")
	return
}

func (lr *LiveRecording) Mute(name string) (err error) {
	err = lr.Conn.Post("/recordings/live/%s/mute", []interface{}{name}, nil, nil)
	return
}

func (lr *LiveRecording) Unmute(name string) (err error) {
	err = lr.Conn.Delete("/recordings/live/%s/mute", []interface{}{name}, nil, "")
	return
}

func (lr *LiveRecording) Delete(name string) (err error) {
	//NOTE: original code used 'stored' for this even though it's live
	err = lr.Conn.Delete("/recordings/stored/%s", []interface{}{name}, nil, "")
	return
}

func (lr *LiveRecording) Scrap(name string) (err error) {
	//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	err = lr.Conn.Delete("/recordings/live/%s", []interface{}{name}, nil, "")
	return
}
