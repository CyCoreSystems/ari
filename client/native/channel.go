package native

import "fmt"

type nativeChannel struct {
	opts *Options
}

func (c *nativeChannel) Hangup(id, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}
	return Delete(c.opts, "/channels/"+id, nil, req)
}

// TODO:  implement ari.Channel interface
