package ari

// Modules is the communication path for interacting with the
// asterisk modules resource
type Modules interface {
	Get(name string) ModuleHandle

	List() ([]ModuleHandle, error)

	Load(name string) error

	Reload(name string) error

	Unload(name string) error

	Data(name string) (*ModuleData, error)
}

// ModuleData is the data for an asterisk module
type ModuleData struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SupportLevel string `json:"support_level"`
	UseCount     int    `json:"use_count"`
	Status       string `json:"status"`
}

// ModuleHandle is the reference to an asterisk module
type ModuleHandle interface {

	// ID returns the identifier for the module
	ID() string

	// Reload reloads the module
	Reload() error

	// Unload unloads the module
	Unload() error

	// Load loads the module
	Load() error

	// Data gets the module data
	Data() (*ModuleData, error)
}
