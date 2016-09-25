package natsgw

func (srv *Server) liveRecording() {
	srv.subscribe("ari.recording.live.data.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.data."):]
		lrd, err := srv.upstream.Recording.Live.Data(name)
		reply(lrd, err)
	})

	srv.subscribe("ari.recording.live.stop.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.stop."):]
		err := srv.upstream.Recording.Live.Stop(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.pause.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.pause."):]
		err := srv.upstream.Recording.Live.Pause(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.resume.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.resume."):]
		err := srv.upstream.Recording.Live.Resume(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.mute.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.mute."):]
		err := srv.upstream.Recording.Live.Mute(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.unmute.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.unmute."):]
		err := srv.upstream.Recording.Live.Unmute(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.delete.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.delete."):]
		err := srv.upstream.Recording.Live.Delete(name)
		reply(nil, err)
	})

	srv.subscribe("ari.recording.live.scrap.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.live.scrap."):]
		err := srv.upstream.Recording.Live.Scrap(name)
		reply(nil, err)
	})

}
