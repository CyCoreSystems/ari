package natsgw

import "encoding/json"

func (srv *Server) storedRecording() {
	srv.subscribe("ari.recording.stored.all", func(_ string, _ []byte, reply Reply) {
		handles, err := srv.upstream.Recording.Stored.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var ret []string
		for _, h := range handles {
			ret = append(ret, h.ID())
		}

		reply(ret, nil)
	})

	srv.subscribe("ari.recording.stored.data.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.recording.stored.data."):]
		srd, err := srv.upstream.Recording.Stored.Data(name)
		reply(srd, err)
	})

	srv.subscribe("ari.recording.stored.copy.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.recording.stored.copy."):]

		var dest string
		if err := json.Unmarshal(data, &dest); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		srd, err := srv.upstream.Recording.Stored.Copy(name, dest)
		reply(srd, err)
	})

	srv.subscribe("ari.recording.stored.delete.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.recording.stored.delete."):]
		err := srv.upstream.Recording.Stored.Delete(name)
		reply(nil, err)
	})

}
