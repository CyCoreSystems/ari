package ari

// TextMessage needs some verbiage here
type TextMessage interface {

	// Send() sends a text message to an endpoint
	Send(from, tech, resource, body string, vars map[string]string) error

	// SendByURI sends a text message to an endpoint by free-form URI
	SendByURI(from, to, body string, vars map[string]string) error
}

// TextMessageData describes text message
type TextMessageData struct {
	Body      string                `json:"body"` // The body (text) of the message
	From      string                `json:"from"` // Technology-specific source URI
	To        string                `json:"to"`   // Technology-specific destination URI
	Variables []TextMessageVariable `json:"variables,omitempty"`
}

// TextMessageVariable describes a key-value pair (associated with a text message)
type TextMessageVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
