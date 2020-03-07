package ari

// Sound represents a communication path to
// the asterisk server for Sound resources
type Sound interface {

	// List returns available sounds limited by the provided filters.
	// Valid filters are "lang", "format", and nil (no filter)
	List(filters map[string]string, keyFilter *Key) ([]*Key, error)

	// Data returns the Sound's data
	Data(key *Key) (*SoundData, error)
}

// SoundData describes a media file which may be played back
type SoundData struct {
	// Key is the cluster-unique identifier for this sound
	Key *Key `json:"key"`

	Formats []FormatLangPair `json:"formats"`
	ID      string           `json:"id"`
	Text    string           `json:"text,omitempty"`
}

// FormatLangPair describes the format and language of a sound file
type FormatLangPair struct {
	Format   string `json:"format"`
	Language string `json:"language"`
}
