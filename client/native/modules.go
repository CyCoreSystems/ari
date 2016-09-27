package native

import "github.com/CyCoreSystems/ari"

type nativeModules struct {
	conn *Conn
}

func (m *nativeModules) Get(name string) *ari.ModuleHandle {
	return ari.NewModuleHandle(name, m)
}

func (m *nativeModules) List() (hx []*ari.ModuleHandle, err error) {
	var modules = []struct {
		Name string `json:"name"`
	}{}

	err = Get(m.conn, "/asterisk/modules", &modules)
	for _, i := range modules {
		hx = append(hx, m.Get(i.Name))
	}

	return
}

func (m *nativeModules) Load(name string) (err error) {
	err = Post(m.conn, "/asterisk/modules/"+name, nil, nil)
	return
}

func (m *nativeModules) Reload(name string) (err error) {
	err = Put(m.conn, "/asterisk/modules/"+name, nil, nil)
	return
}

func (m *nativeModules) Unload(name string) (err error) {
	err = Delete(m.conn, "/asterisk/modules/"+name, nil, "")
	return
}

func (m *nativeModules) Data(name string) (md ari.ModuleData, err error) {
	err = Get(m.conn, "/asterisk/modules/"+name, &md)
	return
}
