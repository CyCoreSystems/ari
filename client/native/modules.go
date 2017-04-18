package native

import "github.com/CyCoreSystems/ari"

// Modules provides the ARI modules accessors for a native client
type Modules struct {
	client *Client
}

// Get obtains a lazy handle to an asterisk module
func (m *Modules) Get(key *ari.Key) ari.ModuleHandle {
	return NewModuleHandle(key, m)
}

// List lists the modules and returns lists of handles
func (m *Modules) List(filter *ari.Key) (hx []*ari.Key, err error) {
	var modules = []struct {
		Name string `json:"name"`
	}{}
	if filter == nil {
		filter = ari.NodeKey(m.client.ApplicationName(), m.client.node)
	}

	err = m.client.get("/asterisk/modules", &modules)
	for _, i := range modules {
		k := ari.NewKey(ari.ModuleKey, i.Name, ari.WithNode(m.client.node), ari.WithApp(m.client.ApplicationName()))
		if filter.Match(k) {
			hx = append(hx, k)
		}
	}

	return
}

// Load loads the named asterisk module
func (m *Modules) Load(key *ari.Key) (err error) {
	name := key.ID
	err = m.client.post("/asterisk/modules/"+name, nil, nil)
	return
}

// Reload reloads the named asterisk module
func (m *Modules) Reload(key *ari.Key) (err error) {
	name := key.ID
	err = m.client.put("/asterisk/modules/"+name, nil, nil)
	return
}

// Unload unloads the named asterisk module
func (m *Modules) Unload(key *ari.Key) (err error) {
	name := key.ID
	err = m.client.del("/asterisk/modules/"+name, nil, "")
	return
}

// Data retrieves the state of the named asterisk module
func (m *Modules) Data(key *ari.Key) (md *ari.ModuleData, err error) {
	md = &ari.ModuleData{}
	name := key.ID
	err = m.client.get("/asterisk/modules/"+name, &md)
	if err != nil {
		md = nil
		err = dataGetError(err, "module", "%v", name)
	}
	return
}

// ModuleHandle is the reference to an asterisk module
type ModuleHandle struct {
	key *ari.Key
	m   *Modules
}

// NewModuleHandle returns a new module handle
func NewModuleHandle(key *ari.Key, m *Modules) ari.ModuleHandle {
	return &ModuleHandle{key, m}
}

// ID returns the identifier for the module
func (mh *ModuleHandle) ID() string {
	return mh.key.ID
}

// Reload reloads the module
func (mh *ModuleHandle) Reload() error {
	return mh.m.Reload(mh.key)
}

// Unload unloads the module
func (mh *ModuleHandle) Unload() error {
	return mh.m.Unload(mh.key)
}

// Load loads the module
func (mh *ModuleHandle) Load() error {
	return mh.m.Load(mh.key)
}

// Data gets the module data
func (mh *ModuleHandle) Data() (*ari.ModuleData, error) {
	return mh.m.Data(mh.key)
}
