package native

import "github.com/CyCoreSystems/ari"

type nativeLiveRecording struct {
	conn *Conn
}

func (lr *nativeLiveRecording) Get(name string) (h *ari.LiveRecordingHandle) {
	h = ari.NewLiveRecordingHandle(name, lr)
	return
}

func (lr *nativeLiveRecording) Data(name string) (d ari.LiveRecordingData, err error) {
	err = Get(lr.conn, "/recordings/live/"+name, &d)
	return
}

func (lr *nativeLiveRecording) Stop(name string) (err error) {
	err = Post(lr.conn, "/recordings/live/"+name+"/stop", nil, nil)
	return
}

func (lr *nativeLiveRecording) Pause(name string) (err error) {
	err = Post(lr.conn, "/recordings/live/"+name+"/pause", nil, nil)
	return
}

func (lr *nativeLiveRecording) Resume(name string) (err error) {
	err = Delete(lr.conn, "/recordings/live/"+name+"/pause", nil, "")
	return
}

func (lr *nativeLiveRecording) Mute(name string) (err error) {
	err = Post(lr.conn, "/recordings/live/"+name+"/mute", nil, nil)
	return
}

func (lr *nativeLiveRecording) Unmute(name string) (err error) {
	err = Delete(lr.conn, "/recordings/live/"+name+"/mute", nil, "")
	return
}

func (lr *nativeLiveRecording) Delete(name string) (err error) {
	//NOTE: original code used 'stored' for this even though it's live
	err = Delete(lr.conn, "/recordings/stored/"+name, nil, "")
	return
}

func (lr *nativeLiveRecording) Scrap(name string) (err error) {
	//TODO reproduce this error in isolation: does not delete. Cannot delete any recording produced by this.
	err = Delete(lr.conn, "/recordings/live/"+name, nil, "")
	return
}
