package native

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
)

// Config provides the ARI Configuration accessors for a native client
type Config struct {
	client *Client
}

// Get gets a lazy handle to a configuration object
func (c *Config) Get(configClass string, objectType string, key *ari.Key) ari.ConfigHandle {
	return NewConfigHandle(configClass, objectType, key, c)
}

// Data retrieves the state of a configuration object
func (c *Config) Data(configClass string, objectType string, key *ari.Key) (cd *ari.ConfigData, err error) {
	cd = &ari.ConfigData{}
	cd.ID = key.ID
	cd.Class = configClass
	cd.Type = objectType
	resourceID := configClass + "/" + objectType + "/" + key.ID
	err = c.client.get("/asterisk/config/dynamic/"+resourceID, &cd.Fields)
	if err != nil {
		cd = nil
		err = dataGetError(err, "config", "%v", resourceID)
	}
	return
}

// Update updates the given configuration object
func (c *Config) Update(configClass string, objectType string, key *ari.Key, tuples []ari.ConfigTuple) (err error) {
	err = c.client.put("/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+key.ID, nil, &tuples)
	return
}

// Delete deletes the configuration object
func (c *Config) Delete(configClass string, objectType string, key *ari.Key) (err error) {
	err = c.client.del("/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+key.ID, nil, "")
	return
}

// NewConfigHandle builds a new config handle
func NewConfigHandle(configClass, objectType string, key *ari.Key, c *Config) ari.ConfigHandle {
	return &ConfigHandle{
		configClass: configClass,
		objectType:  objectType,
		key:         key,
		c:           c,
	}
}

// A ConfigHandle is a reference to a Config object
// on the asterisk service
type ConfigHandle struct {
	configClass string
	objectType  string
	key         *ari.Key

	c *Config
}

// ID returns the unique identifier for the config object
func (ch *ConfigHandle) ID() string {
	return fmt.Sprintf("%v/%v/%v", ch.configClass, ch.objectType, ch.key.ID)
}

// Data gets the current data for the config handle
func (ch *ConfigHandle) Data() (*ari.ConfigData, error) {
	return ch.c.Data(ch.configClass, ch.objectType, ch.key)
}

// Update creates or updates the given config tuples
func (ch *ConfigHandle) Update(tuples []ari.ConfigTuple) error {
	return ch.c.Update(ch.configClass, ch.objectType, ch.key, tuples)
}

// Delete deletes the dynamic configuration object
func (ch *ConfigHandle) Delete() error {
	return ch.c.Delete(ch.configClass, ch.objectType, ch.key)
}
