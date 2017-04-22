package ari

// Modules is the communication path for interacting with the
// asterisk modules resource
type Modules interface {
	Get(key *Key) *ModuleHandle

	List(filter *Key) ([]*Key, error)

	Load(key *Key) error

	Reload(key *Key) error

	Unload(key *Key) error

	Data(key *Key) (*ModuleData, error)
}

// ModuleData is the data for an asterisk module
type ModuleData struct {
	// Key is the cluster-unique identifier for this module
	Key *Key `json:"key"`

	Name         string `json:"name"`
	Description  string `json:"description"`
	SupportLevel string `json:"support_level"`
	UseCount     int    `json:"use_count"`
	Status       string `json:"status"`
}

// ModuleHandle is the reference to an asterisk module
type ModuleHandle struct {
	key *Key
	m   Modules
}

// NewModuleHandle returns a new module handle
func NewModuleHandle(key *Key, m Modules) *ModuleHandle {
	return &ModuleHandle{key, m}
}

// ID returns the identifier for the module
func (mh *ModuleHandle) ID() string {
	return mh.key.ID
}

// Key returns the key for the module
func (mh *ModuleHandle) Key() *Key {
	return mh.key
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
func (mh *ModuleHandle) Data() (*ModuleData, error) {
	return mh.m.Data(mh.key)
}
