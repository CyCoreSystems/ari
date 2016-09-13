package natsgw

import "encoding/json"

func (srv *Server) asterisk() {

	srv.subscribe("ari.asterisk.reload.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.asterisk.reload."):]
		err := srv.upstream.Asterisk.ReloadModule(name)
		reply(nil, err)
	})

	srv.subscribe("ari.asterisk.info", func(_ string, _ []byte, reply Reply) {
		ai, err := srv.upstream.Asterisk.Info("")
		reply(ai, err)
	})

	srv.subscribe("ari.asterisk.variables.get.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.asterisk.variables.get."):]
		val, err := srv.upstream.Asterisk.Variables().Get(name)
		reply(val, err)
	})

	srv.subscribe("ari.asterisk.variables.set.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.asterisk.variables.set."):]

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Asterisk.Variables().Set(name, value)
		reply(nil, err)
	})

}
