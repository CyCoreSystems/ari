package natsgw

import "encoding/json"

func (srv *Server) asterisk() {

	srv.subscribe("ari.asterisk.reload.*",
		NamedHandler(len("ari.asterisk.reload."), func(name string, _ []byte, reply Reply) {
			err := srv.upstream.Asterisk.ReloadModule(name)
			reply(nil, err)
		}))

	srv.subscribe("ari.asterisk.info", func(_ string, _ []byte, reply Reply) {
		ai, err := srv.upstream.Asterisk.Info("")
		reply(ai, err)
	})

	srv.subscribe("ari.asterisk.variables.get.*",
		NamedHandler(len("ari.asterisk.variables.get."), func(name string, _ []byte, reply Reply) {
			val, err := srv.upstream.Asterisk.Variables().Get(name)
			reply(val, err)
		}))

	srv.subscribe("ari.asterisk.variables.set.*",
		NamedHandler(len("ari.asterisk.variables.set."), func(name string, data []byte, reply Reply) {
			var value string
			if err := json.Unmarshal(data, &value); err != nil {
				reply(nil, &decodingError{"ari.asterisk.variables.set." + name, err})
				return
			}

			err := srv.upstream.Asterisk.Variables().Set(name, value)
			reply(nil, err)
		}))

}
