package native

import (
	"errors"

	"github.com/CyCoreSystems/ari/v5"
)

// Modules provides the ARI modules accessors for a native client
type Modules struct {
	client *Client
}

// Get obtains a lazy handle to an asterisk module
func (m *Modules) Get(key *ari.Key) *ari.ModuleHandle {
	return ari.NewModuleHandle(m.client.stamp(key), m)
}

// List lists the modules and returns lists of handles
func (m *Modules) List(filter *ari.Key) (ret []*ari.Key, err error) {
	if filter == nil {
		filter = ari.NodeKey(m.client.appName, m.client.node)
	}

	modules := []struct {
		Name string `json:"name"`
	}{}

	err = m.client.get("/asterisk/modules", &modules)
	if err != nil {
		return nil, err
	}

	for _, i := range modules {
		k := m.client.stamp(ari.NewKey(ari.ModuleKey, i.Name))
		if filter.Match(k) {
			if filter.Dialog != "" {
				k.Dialog = filter.Dialog
			}

			ret = append(ret, k)
		}
	}

	return
}

// Load loads the named asterisk module
func (m *Modules) Load(key *ari.Key) error {
	return m.client.post("/asterisk/modules/"+key.ID, nil, nil)
}

// Reload reloads the named asterisk module
func (m *Modules) Reload(key *ari.Key) error {
	return m.client.put("/asterisk/modules/"+key.ID, nil, nil)
}

// Unload unloads the named asterisk module
func (m *Modules) Unload(key *ari.Key) error {
	return m.client.del("/asterisk/modules/"+key.ID, nil, "")
}

// Data retrieves the state of the named asterisk module
func (m *Modules) Data(key *ari.Key) (*ari.ModuleData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("module key not supplied")
	}

	data := new(ari.ModuleData)
	if err := m.client.get("/asterisk/modules/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "module", "%v", key.ID)
	}

	data.Key = m.client.stamp(key)

	return data, nil
}
