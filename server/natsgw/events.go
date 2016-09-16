package natsgw

import (
	"fmt"

	v2 "github.com/CyCoreSystems/ari/v2"
)

func (srv *Server) events() {

	go func() {
		sub := srv.upstream.Bus.Subscribe(v2.ALL)
		defer sub.Cancel()

		for {
			select {
			case <-srv.ctx.Done():
				return
			case evt := <-sub.Events():
				//app := evt.GetApplication()
				t := evt.GetType()
				subj := fmt.Sprintf("ari.events.%s", t)
				srv.conn.Publish(subj, *evt.GetRaw())
			}
		}
	}()
}
