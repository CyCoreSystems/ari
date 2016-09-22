package natsgw

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
)

func (srv *Server) events() {

	if srv.upstream.Bus == nil {
		// useful for tests
		srv.log.Warn("No Upstream Bus in nats event forwarding")
		return
	}

	go func() {
		sub := srv.upstream.Bus.Subscribe(ari.Events.All)
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
