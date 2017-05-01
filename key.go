package ari

import "fmt"

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

	// VariableKey is the key kind for the ARI Asterisk Variable resource
	VariableKey = "variable"
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

// Keys is a list of keys
type Keys []*Key

// Filter filters the key list using the given key type match
func (kx Keys) Filter(mx ...Matcher) (ret Keys) {
	for _, m := range mx {
		for _, k := range kx {
			if m.Match(k) {
				ret = append(ret, k)
			}
		}
	}
	return
}

// Without removes keys that match the given matcher
func (kx Keys) Without(m Matcher) (ret Keys) {
	for _, k := range kx {
		if !m.Match(k) {
			ret = append(ret, k)
		}
	}
	return
}

// First returns the first key from a list of keys.  It is safe to use on empty lists, in which case, it will return nil.
func (kx Keys) First() *Key {
	if len(kx) < 1 {
		return nil
	}
	return kx[0]
}

// Bridges returns just the bridge keys from a set of Keys
func (kx Keys) Bridges() Keys {
	return kx.Filter(NewKey(BridgeKey, ""))
}

// Channels returns just the channel keys from a set of Keys
func (kx Keys) Channels() Keys {
	return kx.Filter(NewKey(ChannelKey, ""))
}

// ID returns the key from a set of keys with ID matching the given ID.  If the
// key does not exist in the set, nil is returned.
func (kx Keys) ID(id string) *Key {
	return kx.Filter(NewKey("", id)).First()
}

// A Matcher provides the functionality for matching against a key.
type Matcher interface {
	Match(o *Key) bool
}

// MatchFunc is the functional type alias for providing functional `Matcher` implementations
type MatchFunc func(*Key) bool

// Match invokes the match function given the key
func (mf MatchFunc) Match(o *Key) bool {
	return mf(o)
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

// WithLocationOf copies the partial key fields Node, Application, Dialog from the reference key
func WithLocationOf(ref *Key) KeyOptionFunc {
	return func(key Key) Key {
		if ref != nil {
			key.Node = ref.Node
			key.Dialog = ref.Dialog
			key.App = ref.App
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

// ConfigID returns the configuration Key ID for the given configuration class, type/kind, and id.
func ConfigID(class, kind, id string) string {
	return fmt.Sprintf("%s/%s/%s", class, kind, id)
}

// EndpointID returns the endpoint Key ID for the given tech and resource
func EndpointID(tech, resource string) string {
	return fmt.Sprintf("%s/%s", tech, resource)
}

// DialogKey returns a key that is bound to the given dialog.
func DialogKey(dialog string) *Key {
	return NewKey("", "", WithDialog(dialog))
}

// NodeKey returns a key that is bound to the given application and node
func NodeKey(app, node string) *Key {
	return NewKey("", "", WithApp(app), WithNode(node))
}

// KindKey returns a key that is bound by a type only
func KindKey(kind string, opts ...KeyOptionFunc) *Key {
	return NewKey(kind, "", opts...)
}

// Match returns true if the given key matches the subject. Empty partial key fields are wildcards.
func (k *Key) Match(o *Key) bool {
	if k == o {
		return true
	}

	if k == nil || o == nil {
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
	if k.Kind != "" && o.Kind != "" && k.Kind != o.Kind {
		return false
	}
	if k.ID != "" && o.ID != "" && k.ID != o.ID {
		return false
	}

	return true
}

// New returns a new key with the location information from the source key.
// This includes the App, the Node, and the Dialog.  the `kind` and `id`
// parameters are optional.  If kind is empty, the resulting key will not be
// typed.  If id is empty, the key will not be unique.
func (k *Key) New(kind, id string) *Key {
	n := NodeKey(k.App, k.Node)
	n.Dialog = k.Dialog
	n.Kind = kind
	n.ID = id

	return n
}

func (k *Key) String() string {
	if k.ID != "" {
		return k.ID
	}

	if k.Dialog != "" {
		return "[" + k.Dialog + "]"
	}

	if k.Node != "" {
		return k.App + "@" + k.Node
	}

	return "emptyKey"
}
