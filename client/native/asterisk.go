package native

import (
	"fmt"

	"github.com/rotisserie/eris"

	"github.com/PolyAI-LDN/ari/v6"
)

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
func (a *Asterisk) Info(key *ari.Key) (*ari.AsteriskInfo, error) {
	var m ari.AsteriskInfo

	return &m, eris.Wrap(
		a.client.get("/asterisk/info", &m),
		"failed to get asterisk info",
	)
}

// AsteriskVariables provides the ARI Variables accessors for server-level variables
type AsteriskVariables struct {
	client *Client
}

// Variables returns the variables interface for the Asterisk server
func (a *Asterisk) Variables() ari.AsteriskVariables {
	return &AsteriskVariables{a.client}
}

// Get returns the value of the given global variable
// Equivalent to GET /asterisk/variable
func (a *AsteriskVariables) Get(key *ari.Key) (string, error) {
	var m struct {
		Value string `json:"value"`
	}

	err := a.client.get(fmt.Sprintf("/asterisk/variable?variable=%s", key.ID), &m)
	if err != nil {
		return "", eris.Wrapf(err, "Error getting asterisk variable '%v'", key.ID)
	}

	return m.Value, nil
}

// Set sets a global channel variable
// (Equivalent to POST /asterisk/variable)
func (a *AsteriskVariables) Set(key *ari.Key, value string) (err error) {
	req := struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}{
		Variable: key.ID,
		Value:    value,
	}

	return eris.Wrapf(
		a.client.post("/asterisk/variable", nil, &req),
		"Error setting asterisk variable '%s' to '%s'", key.ID, value,
	)
}
