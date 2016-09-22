package native

import "net/url"

type nativeTextMessage struct {
	conn *Conn
}

// Send sends a text message to an endpoint
func (t *nativeTextMessage) Send(from, tech, resource, body string, vars map[string]string) error {
	// Construct querystring values
	v := url.Values{}
	v.Set("from", from)
	v.Set("body", body)

	// vars must not be nil, or Ari will reject the request
	if vars == nil {
		vars = map[string]string{}
	}

	err := Post(t.conn, "/endpoints/"+tech+"/"+resource+"/sendMessage?"+v.Encode(), nil, &vars)
	return err
}

// SendByURI sends a text message to an endpoint by free-form URI (rather than tech/resource)
func (t *nativeTextMessage) SendByURI(from, to, body string, vars map[string]string) error {
	// Construct querystring values
	v := url.Values{}
	v.Set("from", from)
	v.Set("to", to)
	v.Set("body", body)

	// vars must not be nil, or Ari will reject the request
	if vars == nil {
		vars = map[string]string{}
	}

	err := Post(t.conn, "/endpoints/sendMessage?"+v.Encode(), nil, &vars)
	return err
}
