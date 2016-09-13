package natsgw

import "encoding/json"

func (srv *Server) application() {
	srv.subscribe("ari.applications.all", func(_ string, _ []byte, reply Reply) {
		ax, err := srv.upstream.Application.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var apps []string
		for _, a := range ax {
			apps = append(apps, a.ID())
		}

		reply(apps, nil)
	})

	srv.subscribe("ari.applications.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.applications.data."):]
		data, err := srv.upstream.Application.Data(name)
		reply(data, err)
	})

	srv.subscribe("ari.applications.subscribe.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.applications.subscribe."):]

		var eventSource string
		if err := json.Unmarshal(data, &eventSource); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Application.Subscribe(name, eventSource)
		reply(nil, err)
	})

	srv.subscribe("ari.applications.unsubscribe.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.applications.subscribe."):]

		var eventSource string
		if err := json.Unmarshal(data, &eventSource); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Application.Unsubscribe(name, eventSource)
		reply(nil, err)
	})
}
