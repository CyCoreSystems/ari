package ari

import "net/url"

// TextMessage describes text message
type TextMessage struct {
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

// SendMessage sends a text message to an endpoint
func (c *Client) SendMessage(from, tech, resource, body string, vars map[string]string) error {
	// Construct querystring values
	v := url.Values{}
	v.Set("from", from)
	v.Set("body", body)

	// vars must not be nil, or Ari will reject the request
	if vars == nil {
		vars = map[string]string{}
	}

	return c.Put("/endpoints/"+tech+"/"+resource+"/sendMessage?"+v.Encode(), nil, &vars)
}

// SendMessageByUri sends a text message to an endpoint by free-form URI (rather than tech/resource)
func (c *Client) SendMessageByUri(from, to, body string, vars map[string]string) error {
	// Construct querystring values
	v := url.Values{}
	v.Set("from", from)
	v.Set("to", to)
	v.Set("body", body)

	// vars must not be nil, or Ari will reject the request
	if vars == nil {
		vars = map[string]string{}
	}

	return c.Put("/endpoints/sendMessage?"+v.Encode(), nil, &vars)
}
