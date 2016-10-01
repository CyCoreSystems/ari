package native

import "github.com/CyCoreSystems/ari"

type nativeConfig struct {
	conn *Conn
}

func (c *nativeConfig) Get(configClass string, objectType string, id string) *ari.ConfigHandle {
	return ari.NewConfigHandle(configClass, objectType, id, c)
}

func (c *nativeConfig) Data(configClass string, objectType string, id string) (cd ari.ConfigData, err error) {
	cd.ID = id
	cd.Class = configClass
	cd.Type = objectType
	err = Get(c.conn, "/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+id, &cd.Fields)
	return
}

func (c *nativeConfig) Update(configClass string, objectType string, id string, tuples []ari.ConfigTuple) (err error) {
	err = Put(c.conn, "/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+id, nil, &tuples)
	return
}

func (c *nativeConfig) Delete(configClass string, objectType string, id string) (err error) {
	err = Delete(c.conn, "/asterisk/config/dynamic/"+configClass+"/"+objectType+"/"+id, nil, "")
	return
}
