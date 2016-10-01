package natsgw

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/CyCoreSystems/ari"
)

func (srv *Server) config() {
	srv.subscribe("ari.asterisk.config.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.asterisk.config.data."):]

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		cd, err := srv.upstream.Asterisk.Config().Data(items[0], items[1], items[2])
		reply(&cd.Fields, err)
	})

	srv.subscribe("ari.asterisk.config.delete.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.asterisk.config.delete."):]

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		err := srv.upstream.Asterisk.Config().Delete(items[0], items[1], items[2])
		reply(nil, err)
	})

	srv.subscribe("ari.asterisk.config.update.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.asterisk.config.delete."):]

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		var fl []ari.ConfigTuple
		if err := json.Unmarshal(data, &fl); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Asterisk.Config().Update(items[0], items[1], items[2], fl)
		reply(nil, err)
	})

}
