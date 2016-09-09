package generic

import "github.com/CyCoreSystems/ari"

type Endpoint struct {
	Conn Conn
}

func (e *Endpoint) Get(tech string, resource string) *ari.EndpointHandle {
	return ari.NewEndpointHandle(tech, resource, e)
}

func (e *Endpoint) List() (ex []*ari.EndpointHandle, err error) {
	endpoints := []struct {
		Tech     string `json:"technology"`
		Resource string `json:"resource"`
	}{}
	err = e.Conn.Get("/endpoints", nil, &endpoints)
	for _, i := range endpoints {
		ex = append(ex, e.Get(i.Tech, i.Resource))
	}

	return
}

func (e *Endpoint) ListByTech(tech string) (ex []*ari.EndpointHandle, err error) {
	err = e.Conn.Get("/endpoints/%s", []interface{}{tech}, &ex)
	return
}

func (e *Endpoint) Data(tech string, resource string) (ed ari.EndpointData, err error) {
	err = e.Conn.Get("/endpoints/%s/%s", []interface{}{tech, resource}, &ed)
	return
}
