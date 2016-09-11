package native

import (
	"github.com/CyCoreSystems/ari"
	//	"net/url"
)

type nativeEvent struct {
	conn *Conn
}

// Get returns a managed handle to an EventData
func (e *nativeEvent) Get(name string) *ari.EventHandle {
	return ari.NewEventHandle(name, e)
}

// Data returns the details of a given ARI Event
// Equivalent to GET /events/{name}
func (e *nativeEvent) Data(name string) (ed ari.EventData, err error) {
	err = Get(e.conn, "/events/"+name, &ed)
	return ed, err
}
