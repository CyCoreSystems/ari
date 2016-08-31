package ari

import (
	v2 "github.com/CyCoreSystems/ari/v2"
)

// Bus is an event bus for ARI events.  It receives and
// redistributes events based on a subscription model.
type Bus interface {
	Send(m *v2.Message)

	Subscribe(n ...string) *v2.Subscription
}
