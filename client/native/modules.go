package native

import "github.com/CyCoreSystems/ari"

// Modules provides the ARI modules accessors for a native client
type Modules struct {
	client *Client
}

// Get obtains a lazy handle to an asterisk module
func (m *Modules) Get(name string) *ari.ModuleHandle {
	return ari.NewModuleHandle(name, m)
}

// List lists the modules and returns lists of handles
func (m *Modules) List() (hx []*ari.ModuleHandle, err error) {
	var modules = []struct {
		Name string `json:"name"`
	}{}

	err = m.client.get("/asterisk/modules", &modules)
	for _, i := range modules {
		hx = append(hx, m.Get(i.Name))
	}

	return
}

// Load loads the named asterisk module
func (m *Modules) Load(name string) (err error) {
	err = m.client.post("/asterisk/modules/"+name, nil, nil)
	return
}

// Reload reloads the named asterisk module
func (m *Modules) Reload(name string) (err error) {
	err = m.client.put("/asterisk/modules/"+name, nil, nil)
	return
}

// Unload unloads the named asterisk module
func (m *Modules) Unload(name string) (err error) {
	err = m.client.del("/asterisk/modules/"+name, nil, "")
	return
}

// Data retrieves the state of the named asterisk module
func (m *Modules) Data(name string) (md *ari.ModuleData, err error) {
	md = &ari.ModuleData{}
	err = m.client.get("/asterisk/modules/"+name, &md)
	if err != nil {
		md = nil
		err = dataGetError(err, "module", "%v", name)
	}
	return
}
