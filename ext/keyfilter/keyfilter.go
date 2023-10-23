package keyfilter

import "github.com/PolyAI-LDN/ari/v6"

// Kind filters a list of keys by a particular Kind
func Kind(kind string, in []*ari.Key) (out []*ari.Key) {
	for _, k := range in {
		if k.Kind == kind {
			out = append(out, k)
		}
	}

	return
}

// Applications returns the Application keys from the given list of keys
func Applications(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.ApplicationKey, in)
}

// Bridges returns the Bridge keys from the given list of keys
func Bridges(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.BridgeKey, in)
}

// Channels returns the Channel keys from the given list of keys
func Channels(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.ChannelKey, in)
}

// DeviceStates returns the DeviceState keys from the given list of keys
func DeviceStates(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.DeviceStateKey, in)
}

// Endpoints returns the Endpoint keys from the given list of keys
func Endpoints(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.EndpointKey, in)
}

// LiveRecordings returns the LiveRecording keys from the given list of keys
func LiveRecordings(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.LiveRecordingKey, in)
}

// Loggings returns the Logging keys from the given list of keys
func Loggings(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.LoggingKey, in)
}

// Mailboxes returns the Mailbox keys from the given list of keys
func Mailboxes(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.MailboxKey, in)
}

// Modules returns the Module keys from the given list of keys
func Modules(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.ModuleKey, in)
}

// Playbacks returns the Playback keys from the given list of keys
func Playbacks(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.PlaybackKey, in)
}

// Sounds returns the Sound keys from the given list of keys
func Sounds(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.SoundKey, in)
}

// StoredRecordings returns the StoredRecording keys from the given list of keys
func StoredRecordings(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.StoredRecordingKey, in)
}

// Variables returns the Variable keys from the given list of keys
func Variables(in []*ari.Key) (out []*ari.Key) {
	return Kind(ari.VariableKey, in)
}
