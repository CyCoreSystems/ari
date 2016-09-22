package nc

import "github.com/CyCoreSystems/ari"

type natsApplication struct {
	conn *Conn
}

func (a *natsApplication) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

func (a *natsApplication) List() (ax []*ari.ApplicationHandle, err error) {
	var apps []string
	err = a.conn.readRequest("ari.applications.all", nil, &apps)
	for _, app := range apps {
		ax = append(ax, a.Get(app))
	}
	return
}

func (a *natsApplication) Data(name string) (d ari.ApplicationData, err error) {
	err = a.conn.readRequest("ari.applications.data."+name, nil, &d)
	return
}

func (a *natsApplication) Subscribe(name string, eventSource string) (err error) {
	err = a.conn.standardRequest("ari.applications.subscribe."+name, eventSource, nil)
	return
}

func (a *natsApplication) Unsubscribe(name string, eventSource string) (err error) {
	err = a.conn.standardRequest("ari.applications.unsubscribe."+name, eventSource, nil)
	return
}
