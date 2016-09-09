package ari

// Sound represents a communication path to
// the asterisk server for Sound resources
type Sound interface {

	// List returns available sounds limited by the provided filters.
	// Valid filters are "lang", "format", and nil (no filter)
	List(filters map[string]string) ([]*SoundHandle, error)

	// Get returns a handle pointer to the sound for further interaction
	Get(name string) *SoundHandle

	// Data returns the Sound's data
	Data(name string) (SoundData, error)
}

// SoundHandle provides a wrapper to a Sound interface for
// operations on a specific Sound
type SoundHandle struct {
	name string
	s    Sound
}

// ID returns the identifier for the sound
func (sh *SoundHandle) ID() string {
	return sh.name
}

// SoundData describes a media file which may be played back
type SoundData struct {
	Formats []FormatLangPair `json:"formats"`
	ID      string           `json:"id"`
	Text    string           `json:"text,omitempty"`
}

// FormatLangPair describes the format and language of a sound file
type FormatLangPair struct {
	Format   string `json:"format"`
	Language string `json:"language"`
}

// NewSoundHandle creates a new handle to the sound name
func NewSoundHandle(name string, snd Sound) *SoundHandle {
	return &SoundHandle{
		name: name,
		s:    snd,
	}
}

// Data retrieves the data for the Sound
func (sh *SoundHandle) Data() (sd SoundData, err error) {
	sd, err = sh.s.Data(sh.name)
	return sd, err
}
