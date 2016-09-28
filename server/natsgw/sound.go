package natsgw

import (
	"encoding/json"
)

func (srv *Server) sound() {
	srv.subscribe("ari.sounds.all", func(subj string, data []byte, reply Reply) {

		var filters map[string]string
		if err := json.Unmarshal(data, &filters); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		if len(filters) == 0 {
			filters = nil // just send nil to upstream if empty. makes tests easier
		}

		sx, err := srv.upstream.Sound.List(filters)
		if err != nil {
			reply(nil, err)
			return
		}

		var sounds []string
		for _, sound := range sx {
			sounds = append(sounds, sound.ID())
		}

		reply(sounds, nil)
	})

	srv.subscribe("ari.sounds.data.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.sounds.data."):]
		sd, err := srv.upstream.Sound.Data(name)
		reply(&sd, err)
	})

}
