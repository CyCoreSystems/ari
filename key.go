package ari

const (
	// ApplicationKey is the key kind for ARI Application resources.
	ApplicationKey = "application"

	// BridgeKey is the key kind for the ARI Bridge resources.
	BridgeKey = "bridge"

	// ChannelKey is the key kind for the ARI Channel resource
	ChannelKey = "channel"
)

// Key identifies a unique resource in the system
type Key struct {
	// Kind indicates the type of resource the Key points to.  e.g., "channel",
	// "bridge", etc.
	Kind string `json:"kind"`

	// ID indicates the unique identifier of the resource
	ID string `json:"id"`

	// Node indicates the unique identifier of the Asterisk node on which the
	// resource exists or will be created
	Node string `json:"node,omitempty"`

	// Dialog indicates a named scope of the resource, for receiving events
	Dialog string `json:"dialog,omitempty"`
}

// KeyOptionFunc is a functional argument alias for providing options for ARI keys
type KeyOptionFunc func(Key)

// WithDialog sets the given dialog identifier on the key.
func WithDialog(dialog string) KeyOptionFunc {
	return func(key Key) {
		key.Dialog = dialog
	}
}

// WithNode sets the given node identifier on the key.
func WithNode(node string) KeyOptionFunc {
	return func(key Key) {
		key.Node = node
	}
}

// NewKey builds a new key given the kind, identifier, and any optional arguments.
func NewKey(kind string, id string, opts ...KeyOptionFunc) *Key {
	k := Key{
		Kind: kind,
		ID:   id,
	}
	for _, o := range opts {
		o(k)
	}

	return &k
}
