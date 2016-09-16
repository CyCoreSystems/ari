package natsgw

import (
	"encoding/json"

	"github.com/nats-io/nats"
)

func (srv *Server) subscribe(endpoint string, h Handler) {

	cb := func(msg *nats.Msg) {

		reply := msg.Reply
		data := msg.Data
		subj := msg.Subject

		srv.conn.Publish(reply+".ok", []byte(""))

		h(subj, data, func(i interface{}, err error) {

			if err != nil {
				resp := []byte(err.Error())
				//			resp := append([]byte{nc.MagicRespErr}, resp...)
				if err2 := srv.conn.Publish(reply+".err", resp); err2 != nil {
					srv.log.Error("Error sending error reply", "reply", reply, "error", err2, "body", err)
				}
				return
			}

			resp := []byte("{}")
			if i != nil {
				resp, err = json.Marshal(i)
				if err != nil {
					srv.log.Error("Error building response reply", "reply", reply, "error", err)
					return
				}
			}

			//resp = append([]byte{nc.MagicRespOK}, resp...)
			if err = srv.conn.Publish(reply+".resp", resp); err != nil {
				srv.log.Error("Error sending response reply", "reply", reply, "error", err, "body", resp)
			}
		})
	}

	sub, err := srv.conn.Subscribe(endpoint, cb)
	if err != nil {
		srv.log.Error("Error starting subscription", "endpoint", endpoint, "error", err)
		srv.Close()
		return
	}

	go func() {
		<-srv.ctx.Done()
		sub.Unsubscribe()
	}()
}
