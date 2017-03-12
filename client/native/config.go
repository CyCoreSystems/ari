package native

import "github.com/CyCoreSystems/ari"

// Config provides the ARI Configuration accessors for a native client
type Config struct {
	client *Client
}

// Get gets a lazy handle to a configuration object
func (c *Config) Get(configClass string, objectType string, id string) *ari.ConfigHandle {
	return ari.NewConfigHandle(configClass, objectType, id, c)
}

// Data retrieves the state of a configuration object
func (c *Config) Data(configClass string, objectType string, id string) (cd *ari.ConfigData, err error) {
	cd = &ari.ConfigData{}
	cd.ID = id
	cd.Class = configClass
	cd.Type = objectType
	resourceID := configClass + "/" + objectType + "/" + id
	err = c.client.get("/asterisk/config/dynamic/"+resourceID, &cd.Fields)
	if err != nil {
		cd = nil
		err = dataGetError(err, "config", "%v", resourceID)
	}
	return
}

// Update updates the given configuration object
func (c *Config) Update(configClass string, objectType string, id string, tuples []ari.ConfigTuple) (err error) {
	err = c.client.put("/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+id, nil, &tuples)
	return
}

// Delete deletes the configuration object
func (c *Config) Delete(configClass string, objectType string, id string) (err error) {
	err = c.client.del("/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+id, nil, "")
	return
}
