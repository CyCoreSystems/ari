package natsgw

import "encoding/json"

func (srv *Server) playback() {

	srv.subscribe("ari.playback.data.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.playback.data."):]
		d, err := srv.upstream.Playback.Data(name)
		reply(&d, err)
	})

	srv.subscribe("ari.playback.control.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.playback.control."):]

		var command string
		if err := json.Unmarshal(data, &command); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Playback.Control(name, command)
		reply(nil, err)
	})

	srv.subscribe("ari.playback.stop.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.playback.stop."):]
		err := srv.upstream.Playback.Stop(name)
		reply(nil, err)
	})

}
