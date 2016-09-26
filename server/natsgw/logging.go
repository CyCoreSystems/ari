package natsgw

import "encoding/json"

func (srv *Server) logging() {
	srv.subscribe("ari.logging.create.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.logging.create."):]

		var config string
		if err := json.Unmarshal(data, &config); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Asterisk.Logging().Create(name, config)
		reply(nil, err)
		return
	})

	srv.subscribe("ari.logging.all", func(_ string, _ []byte, reply Reply) {
		ld, err := srv.upstream.Asterisk.Logging().List()
		reply(ld, err)
	})

	srv.subscribe("ari.logging.delete.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.logging.delete."):]
		err := srv.upstream.Asterisk.Logging().Delete(name)
		reply(nil, err)
	})

	srv.subscribe("ari.logging.rotate.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.logging.rotate."):]
		err := srv.upstream.Asterisk.Logging().Rotate(name)
		reply(nil, err)
	})

}
