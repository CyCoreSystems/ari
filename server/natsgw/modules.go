package natsgw

func (srv *Server) modules() {
	srv.subscribe("ari.modules.all", func(_ string, _ []byte, reply Reply) {
		mx, err := srv.upstream.Asterisk.Modules().List()
		if err != nil {
			reply(nil, err)
			return
		}

		var modules []string
		for _, m := range mx {
			modules = append(modules, m.ID())
		}

		reply(modules, nil)
	})

	srv.subscribe("ari.modules.data.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.modules.data."):]
		data, err := srv.upstream.Asterisk.Modules().Data(name)
		reply(data, err)
	})

	srv.subscribe("ari.modules.load.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.modules.load."):]
		err := srv.upstream.Asterisk.Modules().Load(name)
		reply(nil, err)
	})

	srv.subscribe("ari.modules.unload.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.modules.unload."):]
		err := srv.upstream.Asterisk.Modules().Unload(name)
		reply(nil, err)
	})

	srv.subscribe("ari.modules.reload.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.modules.reload."):]
		err := srv.upstream.Asterisk.Modules().Reload(name)
		reply(nil, err)
	})

}
