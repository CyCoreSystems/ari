package native

import (
	"github.com/AVOXI/ari"
	"github.com/pkg/errors"
)

// Config provides the ARI Configuration accessors for a native client
type Config struct {
	client *Client
}

// Get gets a lazy handle to a configuration object
func (c *Config) Get(key *ari.Key) *ari.ConfigHandle {
	return ari.NewConfigHandle(key, c)
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
