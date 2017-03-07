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

/*
	conn    *Conn
	logging ari.Logging
	modules ari.Modules
	config  ari.Config
}
*/

// Info returns various data about the Asterisk system
// Equivalent to GET /asterisk/info
func (a *Asterisk) Info(only string) (*ari.AsteriskInfo, error) {
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

	err := a.client.conn.Get(path, &m)
	return &m, err
}

// ReloadModule tells asterisk to load the given module
func (a *Asterisk) ReloadModule(name string) error {
	return a.Modules().Reload(name)
}

// AsteriskVariables provides the ARI Variables accessors for server-level variables
type AsteriskVariables struct {
	client *Client
}

// Variables returns the variables interface for the Asterisk server
func (a *Asterisk) Variables() ari.Variables {
	return &AsteriskVariables{a.client}
}

// Get returns the value of the given global variable
// Equivalent to GET /asterisk/variable
func (a *AsteriskVariables) Get(key string) (string, error) {
	type variable struct {
		Value string `json:"value"`
	}

	var m variable

	path := "/asterisk/variable?variable=" + key
	err := a.client.conn.Get(path, &m)
	if err != nil {
		return "", err
	}
	return m.Value, nil
}

// Set sets a global channel variable
// (Equivalent to POST /asterisk/variable)
func (a *AsteriskVariables) Set(key string, value string) error {
	path := "/asterisk/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{key, value}

	err := a.client.conn.Post(path, nil, &req)
	return err
}
