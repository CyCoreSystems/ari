package natsgw

import "encoding/json"

func (srv *Server) mailbox() {
	srv.subscribe("ari.mailboxes.all", func(_ string, _ []byte, reply Reply) {
		mx, err := srv.upstream.Mailbox.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var mailboxes []string
		for _, m := range mx {
			mailboxes = append(mailboxes, m.ID())
		}

		reply(mailboxes, nil)
	})

	srv.subscribe("ari.mailboxes.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.mailboxes.data."):]
		data, err := srv.upstream.Mailbox.Data(name)
		reply(data, err)
	})

	srv.subscribe("ari.mailboxes.update.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.mailboxes.delete."):]

		type req struct {
			Old int `json:"old"`
			New int `json:"new"`
		}

		var request req
		if err := json.Unmarshal(data, &request); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Mailbox.Update(name, request.Old, request.New)
		reply(nil, err)
	})

	srv.subscribe("ari.mailboxes.delete.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.mailboxes.delete."):]
		err := srv.upstream.Mailbox.Delete(name)
		reply(nil, err)
	})

}
