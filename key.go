package ari

const (
	// ApplicationKey is the key kind for ARI Application resources.
	ApplicationKey = "application"

	// BridgeKey is the key kind for the ARI Bridge resources.
	BridgeKey = "bridge"

	// ChannelKey is the key kind for the ARI Channel resource
	ChannelKey = "channel"

	// DeviceStateKey is the key kind for the ARI DeviceState resource
	DeviceStateKey = "devicestate"

	// EndpointKey is the key kind for the ARI Endpoint resource
	EndpointKey = "endpoint"

	// LiveRecordingKey is the key kind for the ARI LiveRecording resource
	LiveRecordingKey = "liverecording"

	// LoggingKey is the key kind for the ARI Logging resource
	LoggingKey = "logging"

	// MailboxKey is the key kind for the ARI Mailbox resource
	MailboxKey = "mailbox"

	// ModuleKey is the key kind for the ARI Module resource
	ModuleKey = "module"

	// PlaybackKey is the key kind for the ARI Playback resource
	PlaybackKey = "playback"

	// SoundKey is the key kind for the ARI Sound resource
	SoundKey = "sound"

	// StoredRecordingKey is the key kind for the ARI StoredRecording resource
	StoredRecordingKey = "storedrecording"
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

	// App indiciates the ARI application that this key is bound to.
	App string `json:"app,omitempty"`
}

// KeyOptionFunc is a functional argument alias for providing options for ARI keys
type KeyOptionFunc func(Key) Key

// WithDialog sets the given dialog identifier on the key.
func WithDialog(dialog string) KeyOptionFunc {
	return func(key Key) Key {
		key.Dialog = dialog
		return key
	}
}

// WithNode sets the given node identifier on the key.
func WithNode(node string) KeyOptionFunc {
	return func(key Key) Key {
		key.Node = node
		return key
	}
}

// WithApp sets the given node identifier on the key.
func WithApp(app string) KeyOptionFunc {
	return func(key Key) Key {
		key.App = app
		return key
	}
}

// WithParent copies the partial key fields Node, Application, Dialog from the parent key
func WithParent(parent *Key) KeyOptionFunc {
	return func(key Key) Key {
		if parent != nil {
			key.Node = parent.Node
			key.Dialog = parent.Dialog
			key.App = parent.App
		}
		return key
	}
}

// NewKey builds a new key given the kind, identifier, and any optional arguments.
func NewKey(kind string, id string, opts ...KeyOptionFunc) *Key {
	k := Key{
		Kind: kind,
		ID:   id,
	}
	for _, o := range opts {
		k = o(k)
	}

	return &k
}

// AppKey returns a key that is bound to the given application.
func AppKey(app string) *Key {
	return NewKey("", "", WithApp(app))
}

// DialogKey returns a key that is bound to the given dialog.
func DialogKey(dialog string) *Key {
	return NewKey("", "", WithDialog(dialog))
}

// NodeKey returns a key that is bound to the given application and node
func NodeKey(app, node string) *Key {
	return NewKey("", "", WithApp(app), WithNode(node))
}

// Match returns true if the given key matches the subject. Empty partial key fields are wildcards.
func (k *Key) Match(o *Key) bool {
	if k == o {
		return true
	}

	if k.App != "" && o.App != "" && k.App != o.App {
		return false
	}
	if k.Dialog != "" && o.Dialog != "" && k.Dialog != o.Dialog {
		return false
	}
	if k.Node != "" && o.Node != "" && k.Node != o.Node {
		return false
	}

	if k.Kind == "" && k.ID != "" && k.ID != o.ID {
		return false
	}

	if k.ID == "" && k.Kind != "" && k.Kind != o.Kind {
		return false
	}

	return true
}
