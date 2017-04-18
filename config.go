package ari

// Config represents a transport to the asterisk
// config ARI resource.
type Config interface {

	// Get gets the reference to a config object
	Get(configClass, objectType string, key *Key) ConfigHandle

	// Data gets the data for the config object
	Data(configClass, objectType string, key *Key) (*ConfigData, error)

	// Update creates or updates the given tuples
	Update(configClass, objectType string, key *Key, tuples []ConfigTuple) error

	// Delete deletes the dynamic configuration object.
	Delete(configClass, objectType string, key *Key) error
}

// ConfigData contains the data for a given configuration object
type ConfigData struct {
	ID     string
	Class  string
	Type   string
	Fields []ConfigTuple
}

// ConfigTuple is the key-value pair that defines a configuration entry
type ConfigTuple struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

// A ConfigHandle is a reference to a Config object
// on the asterisk service
type ConfigHandle interface {

	// ID returns the unique identifier for the config object
	ID() string

	// Data gets the current data for the config handle
	Data() (*ConfigData, error)

	// Update creates or updates the given config tuples
	Update(tuples []ConfigTuple) error

	// Delete deletes the dynamic configuration object
	Delete() error
}
