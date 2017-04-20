package native

import (
	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
)

// Config provides the ARI Configuration accessors for a native client
type Config struct {
	client *Client
}

// Get gets a lazy handle to a configuration object
func (c *Config) Get(key *ari.Key) ari.ConfigHandle {
	return NewConfigHandle(key, c)
}

// Data retrieves the state of a configuration object
func (c *Config) Data(key *ari.Key) (*ari.ConfigData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("config key not supplied")
	}

	class, kind, name, err := ari.ParseConfigID(key.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse configuration key")
	}

	data := &ari.ConfigData{
		Key:   c.client.stamp(key),
		Class: class,
		Type:  kind,
		Name:  name,
	}
	err = c.client.get("/asterisk/config/dynamic/"+key.ID, &data.Fields)
	if err != nil {
		return nil, dataGetError(err, "config", "%v", key.ID)
	}
	return data, nil
}

// Update updates the given configuration object
func (c *Config) Update(key *ari.Key, tuples []ari.ConfigTuple) (err error) {
	class, kind, name, err := ari.ParseConfigID(key.ID)
	if err != nil {
		return errors.Wrap(err, "failed to parse key")
	}
	return c.client.put("/asterisk/config/dynamic/"+class+"/"+kind+"/"+name, nil, &tuples)
}

// Delete deletes the configuration object
func (c *Config) Delete(key *ari.Key) error {
	class, kind, name, err := ari.ParseConfigID(key.ID)
	if err != nil {
		return errors.Wrap(err, "failed to parse key")
	}
	return c.client.del("/asterisk/config/dynamic/"+class+"/"+kind+"/"+name, nil, "")
}

// NewConfigHandle builds a new config handle
func NewConfigHandle(key *ari.Key, c *Config) ari.ConfigHandle {
	return &ConfigHandle{
		key: key,
		c:   c,
	}
}

// A ConfigHandle is a reference to a Config object
// on the asterisk service
type ConfigHandle struct {
	key *ari.Key

	c *Config
}

// ID returns the unique identifier for the config object
func (ch *ConfigHandle) ID() string {
	return ch.key.ID
}

// Data gets the current data for the config handle
func (ch *ConfigHandle) Data() (*ari.ConfigData, error) {
	return ch.c.Data(ch.key)
}

// Update creates or updates the given config tuples
func (ch *ConfigHandle) Update(tuples []ari.ConfigTuple) error {
	return ch.c.Update(ch.key, tuples)
}

// Delete deletes the dynamic configuration object
func (ch *ConfigHandle) Delete() error {
	return ch.c.Delete(ch.key)
}
