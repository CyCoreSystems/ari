package native

import "github.com/CyCoreSystems/ari"

type nativeEndpoint struct {
	conn *Conn
}

func (e *nativeEndpoint) Get(tech string, resource string) *ari.EndpointHandle {
	return ari.NewEndpointHandle(tech, resource, e)
}

func (e *nativeEndpoint) List() (ex []*ari.EndpointHandle, err error) {
	endpoints := []struct {
		Tech     string `json:"technology"`
		Resource string `json:"resource"`
	}{}
	err = Get(e.conn, "/endpoints", &endpoints)
	for _, i := range endpoints {
		ex = append(ex, e.Get(i.Tech, i.Resource))
	}

	return
}

func (e *nativeEndpoint) ListByTech(tech string) (ex []*ari.EndpointHandle, err error) {
	err = Get(e.conn, "/endpoints/"+tech, &ex)
	return
}

func (e *nativeEndpoint) Data(tech string, resource string) (ed ari.EndpointData, err error) {
	err = Get(e.conn, "/endpoints/"+tech+"/"+resource, &ed)
	return
}
