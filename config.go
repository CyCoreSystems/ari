package ari

import "fmt"

// Config represents a transport to the asterisk
// config ARI resource.
type Config interface {

	// Get gets the reference to a config object
	Get(configClass, objectType, id string) *ConfigHandle

	// Data gets the data for the config object
	Data(configClass, objectType, id string) (ConfigData, error)

	// Update creates or updates the given tuples
	Update(configClass, objectType, id string, tuples []ConfigTuple) error

	// Delete deletes the dynamic configuration object.
	Delete(configClass, objectType, id string) error
}

// NewConfigHandle builds a new config handle
func NewConfigHandle(configClass, objectType, id string, c Config) *ConfigHandle {
	return &ConfigHandle{
		configClass: configClass,
		objectType:  objectType,
		id:          id,
		c:           c,
	}
}

// A ConfigHandle is a reference to a Config object
// on the asterisk service
type ConfigHandle struct {
	configClass string
	objectType  string
	id          string

	c Config
}

// ID returns the unique identifier for the config object
func (ch *ConfigHandle) ID() string {
	return fmt.Sprintf("%v/%v/%v", ch.configClass, ch.objectType, ch.id)
}

// Data gets the current data for the config handle
func (ch *ConfigHandle) Data() (ConfigData, error) {
	return ch.c.Data(ch.configClass, ch.objectType, ch.id)
}

// Update creates or updates the given config tuples
func (ch *ConfigHandle) Update(tuples []ConfigTuple) error {
	return ch.c.Update(ch.configClass, ch.objectType, ch.id, tuples)
}

// Delete deletes the dynamic configuration object
func (ch *ConfigHandle) Delete() error {
	return ch.c.Delete(ch.configClass, ch.objectType, ch.id)
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
