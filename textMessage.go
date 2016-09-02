package ari

// TextMessage needs some verbiage here
type TextMessage interface {

	// Send() sends a text message to an endpoint
	Send(from, tech, resource, body string, vars map[string]string) error

	// SendByURI sends a text message to an endpoint by free-form URI
	SendByURI(from, to, body string, vars map[string]string) error
}
