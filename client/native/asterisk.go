package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

var errOnlyUnsupported = errors.New("Only-restricted AsteriskInfo requests are not yet implemented")

// Asterisk provides the ARI Asterisk accessors for a native client
type Asterisk struct {
	client *Client
}

// Info returns the intformation about the connected Asterisk system
func (a *Asterisk) Info(only string) (*ari.AsteriskInfo, error) {
	panic("not implemented")
}

// Variables provides the ARI Asterisk Variables accessors for a native client
func (a *Asterisk) Variables() ari.Variables {
	return &Variables{a.client}
}

// Variables provides the ARI Asterisk Variables accessors for a native client
type Variables struct {
	client *Client
}

// Get retrieves a global variable
func (v *Variables) Get(name string) (string, error) {
	var resp struct {
		Value string `json:"value"`
	}

	path := "/asterisk/variable?variable=" + name
	err := v.client.conn.Get(path, &resp)
	return resp.Value, nil
}

// Set sets a global variable
func (v *Variables) Set(name, val string) error {
	path := "/asterisk/variable"

	req := struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}{
		Variable: name,
		Value:    value,
	}

	return v.client.conn.Post(path, nil, &req)
}

// Logging provides the ARI Asterisk Logging accessors for a native client
func (a *Asterisk) Logging() ari.Logging {
	return &Logging{a.client}
}

// Modules provides the ARI Asterisk Modules accessors for a native client
func (a *Asterisk) Modules() ari.Modules {
	return &Modules{a.client}
}

// Config provides the ARI Asterisk Config accessors for a native client
func (a *Asterisk) Config() ari.Config {
	return &Config{a.client}
}

// ReloadModule requests a particular Asterisk module to be reloaded
func (a *Asterisk) ReloadModule(name string) error {
	panic("not implemented")
}

/*
	conn    *Conn
	logging ari.Logging
	modules ari.Modules
	config  ari.Config
}
*/

// Config returns the config resource
func (a *nativeAsterisk) Config() ari.Config {
	return a.config
}

// Modules returns the modules resource
func (a *nativeAsterisk) Modules() ari.Modules {
	return a.modules
}

// Logging returns the logging resource
func (a *nativeAsterisk) Logging() ari.Logging {
	return a.logging
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

// ReloadModule tells asterisk to load the given module
func (a *nativeAsterisk) ReloadModule(name string) error {
	return a.Modules().Reload(name)
}

type nativeAsteriskVariables struct {
	conn *Conn
}

// Variables returns the variables interface for the Asterisk server
func (a *nativeAsterisk) Variables() ari.Variables {
	return &nativeAsteriskVariables{a.conn}
}

// Get returns the value of the given global variable
// Equivalent to GET /asterisk/variable
func (a *nativeAsteriskVariables) Get(key string) (string, error) {
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

// Set sets a global channel variable
// (Equivalent to POST /asterisk/variable)
func (a *nativeAsteriskVariables) Set(key string, value string) error {
	path := "/asterisk/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{key, value}

	err := Post(a.conn, path, nil, &req)
	return err
}
