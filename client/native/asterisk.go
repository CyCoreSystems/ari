package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

var errOnlyUnsupported = errors.New("Only-restricted AsteriskInfo requests are not yet implemented")

type nativeAsterisk struct {
	conn *Conn
}

// Info returns various data about the Asterisk system
// Equivalent to GET /asterisk/info
func (a *nativeAsterisk) Info(only string) (*ari.AsteriskInfo, error) {
	var m ari.AsteriskInfo
	path := "/asterisk/info"

	// If we are passed an 'only' parameter
	// pass it on as the 'only' querystring parameter
	if only != "" {
		path += "?only=" + only
		return &m, errOnlyUnsupported
	}
	// TODO: handle "only" parameter
	// the problem is that responses with "only" do not
	// conform to the AsteriskInfo model; they just return
	// the subobjects requested
	// That means we should probably break this
	// method into multiple submethods

	err := Get(a.conn, path, &m)
	return &m, err
}

// GetVariable returns the value of the given global variable
// Equivalent to GET /asterisk/variable
func (a *nativeAsterisk) GetVariable(key string) (string, error) {
	type variable struct {
		Value string `json:"value"`
	}

	var m variable

	path := "/asterisk/variable?variable=" + key
	err := Get(a.conn, path, &m)
	if err != nil {
		return "", err
	}
	return m.Value, nil
}

// SetVariable sets a global channel variable
// (Equivalent to POST /asterisk/variable)
func (a *nativeAsterisk) SetVariable(key string, value string) error {
	path := "/asterisk/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{key, value}

	err := Post(a.conn, path, nil, &req)
	return err
}

// ReloadModule tells asterisk to load the given module
func (a *nativeAsterisk) ReloadModule(name string) error {
	err := Put(a.conn, "/asterisk/modules/"+name, nil, nil)
	return err
}
