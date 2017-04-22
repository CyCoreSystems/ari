package ari

// Sound represents a communication path to
// the asterisk server for Sound resources
type Sound interface {

	// List returns available sounds limited by the provided filters.
	// Valid filters are "lang", "format", and nil (no filter)
	List(filters map[string]string, keyFilter *Key) ([]*Key, error)

	// Get returns a handle pointer to the sound for further interaction
	Get(key *Key) *SoundHandle

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

// SoundHandle provides a wrapper to a Sound interface for
// operations on a specific Sound
type SoundHandle struct {
	key *Key
	s   Sound
}

// NewSoundHandle creates a new handle to the sound name
func NewSoundHandle(key *Key, snd Sound) *SoundHandle {
	return &SoundHandle{
		key: key,
		s:   snd,
	}
}

// ID returns the identifier for the sound
func (sh *SoundHandle) ID() string {
	return sh.key.ID
}

// Key returns the Key for the sound
func (sh *SoundHandle) Key() *Key {
	return sh.key
}

// Data retrieves the data for the Sound
func (sh *SoundHandle) Data() (sd *SoundData, err error) {
	sd, err = sh.s.Data(sh.key)
	return sd, err
}
