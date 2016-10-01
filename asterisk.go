package ari

// Asterisk represents a communication path for
// the Asterisk server for system-level resources
type Asterisk interface {

	// Info gets data about the asterisk system
	Info(only string) (*AsteriskInfo, error)

	// Variables returns the global asterisk variables
	Variables() Variables

	// Logging returns the interface for working with asterisk logs
	Logging() Logging

	// Modules returns the interface for working with asterisk modules
	Modules() Modules

	// Config returns the interface for working with dynamic configuration
	Config() Config

	// ReloadModule tells asterisk to load the given module
	ReloadModule(name string) error
}

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
	SetID           SetID   `json:"setid"` // Effective user/group id under which Asterisk is running
}

// SetID describes a userid/groupid pair
type SetID struct {
	Group string `json:"group"` // group id (not name? why string?)
	User  string `json:"user"`  // user id (not name? why string?)
}

// StatusInfo describes the state of an Asterisk system
type StatusInfo struct {
	LastReloadTime DateTime `json:"last_reload_time"`
	StartupTime    DateTime `json:"startup_time"`
}

// SystemInfo describes information about the Asterisk system
type SystemInfo struct {
	EntityID string `json:"entity_id"`
	Version  string `json:"version"`
}
