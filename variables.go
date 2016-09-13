package ari

// Variables represents a set of variables attached to an entity (Asterisk Server, Channel, etc)
type Variables interface {

	// Get returns the value of the given variable
	Get(variable string) (string, error)

	// Set sets the variable
	Set(variable string, value string) error
}
