package ari

// Modules is the communication path for interacting with the
// asterisk modules resource
type Modules interface {
	Get(name string) *ModuleHandle

	List() ([]*ModuleHandle, error)

	Load(name string) error

	Reload(name string) error

	Unload(name string) error

	Data(name string) (ModuleData, error)
}

// NewModuleHandle returns a new module handle
func NewModuleHandle(name string, m Modules) *ModuleHandle {
	return &ModuleHandle{name, m}
}

// ModuleHandle is the reference to an asterisk module
type ModuleHandle struct {
	name string
	m    Modules
}

// ID returns the identifier for the module
func (mh *ModuleHandle) ID() string {
	return mh.name
}

// Reload reloads the module
func (mh *ModuleHandle) Reload() error {
	return mh.m.Reload(mh.name)
}

// Unload unloads the module
func (mh *ModuleHandle) Unload() error {
	return mh.m.Unload(mh.name)
}

// Load loads the module
func (mh *ModuleHandle) Load() error {
	return mh.m.Load(mh.name)
}

// Data gets the module data
func (mh *ModuleHandle) Data() (ModuleData, error) {
	return mh.m.Data(mh.name)
}

// ModuleData is the data for an asterisk module
type ModuleData struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SupportLevel string `json:"support_level"`
	UseCount     int    `json:"use_count"`
	Status       string `json:"status"`
}
