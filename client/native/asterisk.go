package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

var errOnlyUnsupported = errors.New("Only-restricted AsteriskInfo requests are not yet implemented")

type nativeAsterisk struct {
	conn    *Conn
	logging ari.Logging
	modules ari.Modules
	config  ari.Config
}

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
