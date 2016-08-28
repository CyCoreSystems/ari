package native

import "github.com/CyCoreSystems/ari"

type nativeEndpoint struct {
	conn *Conn
}

func (e *nativeEndpoint) Get(tech string, resource string) *ari.EndpointHandle {
	return ari.NewEndpointHandle(tech, resource, e)
}

func (e *nativeEndpoint) Data(tech string, resource string) (ed ari.EndpointData, err error) {
	err = Get(e.conn, "/endpoints/"+tech+"/"+resource, &ed)
	return
}
