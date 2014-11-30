package ari

import "fmt"

// AsteriskInfo describes a running asterisk system
type AsteriskInfo struct {
	BuildInfo  BuildInfo  `json:"build"`
	ConfigInfo ConfigInfo `json:"config"`
	StatusInfo StatusInfo `json:"status"`
	SystemInfo SystemInfo `json:"system"`
}

// BuildInfo describes information about how Asterisk was built
type BuildInfo struct {
	Date    string `json:"date"`
	Kernel  string `json:"kernel"`
	Machine string `json:"machine"`
	Options string `json:"options"`
	Os      string `json:"os"`
	User    string `json:"user"`
}

// ConfigInfo describes information about the Asterisk configuration
type ConfigInfo struct {
	DefaultLanguage string  `json:"default_language"`
	MaxChannels     int     `json:"max_channels,omitempty"` //omitempty denotes an optional field, meaning the field may not be present if no value is assigned.
	MaxLoad         float64 `json:"max_load,omitempty"`
	MaxOpenFiles    int     `json:"max_open_files,omitempty"`
	Name            string  `json:"name"`  // Asterisk system name
	SetId           SetId   `json:"setid"` // Effective user/group id under which Asterisk is running
}

// SetId describes a userid/groupid pair
type SetId struct {
	Group string `json:"group"` // group id (not name? why string?)
	User  string `json:"user"`  // user id (not name? why string?)
}

// StatusInfo describes the state of an Asterisk system
type StatusInfo struct {
	LastReloadTime AsteriskDate `json:"last_reload_time"`
	StartupTime    AsteriskDate `json:"startup_time"`
}

// SystemInfo describes information about the Asterisk system
type SystemInfo struct {
	EntityId string `json:"entity_id"`
	Version  string `json:"version"`
}

// GetAsteriskInfo returns various data about the Asterisk
// system
// Equivalent to GET /asterisk/info
func (c *Client) GetAsteriskInfo(only string) (*AsteriskInfo, error) {
	var m AsteriskInfo
	path := "/asterisk/info"

	// If we are passed an 'only' parameter
	// pass it on as the 'only' querystring parameter
	if only != "" {
		path += "?only=" + only
		return &m, fmt.Errorf("Only-restricted AsteriskInfo requests are not yet implemented")
	}
	// TODO: handle "only" parameter
	// the problem is that responses with "only" do not
	// conform to the AsteriskInfo model; they just return
	// the subobjects requested
	// That means we should probably break this
	// method into multiple submethods

	err := c.AriGet(path, &m)
	if err != nil {
		return &m, err
	}
	return &m, nil
}

// GetAsteriskVariable returns the value of the given global variable
// Equivalent to GET /asterisk/variable
func (c *Client) GetAsteriskVariable(variable string) (string, error) {
	var m Variable
	path := "/asterisk/variable?variable=" + variable
	err := c.AriGet(path, &m)
	if err != nil {
		return "", err
	}
	return m.Value, nil
}

//Equivalent to POST /asterisk/variable
func (c *Client) SetAsteriskVariable(variable string, value string) error {
	path := "/asterisk/variable"

	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}
	req := request{variable, value}

	err := c.AriPost(path, nil, &req)
	if err != nil {
		return err
	}
	return nil
}
