package nc

import "github.com/CyCoreSystems/ari"

type natsLiveRecording struct {
	conn *Conn
}

func (lr *natsLiveRecording) Get(name string) *ari.LiveRecordingHandle {
	return ari.NewLiveRecordingHandle(name, lr)
}

func (lr *natsLiveRecording) Data(name string) (lrd ari.LiveRecordingData, err error) {
	err = lr.conn.readRequest("ari.recording.live.data."+name, nil, &lrd)
	return
}

func (lr *natsLiveRecording) Stop(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.stop."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Pause(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.pause."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Resume(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.resume."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Mute(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.mute."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Unmute(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.unmute."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Delete(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.delete."+name, nil, nil)
	return
}

func (lr *natsLiveRecording) Scrap(name string) (err error) {
	err = lr.conn.standardRequest("ari.recording.live.scrap."+name, nil, nil)
	return
}
