package native

import (
	"net/url"

	"github.com/CyCoreSystems/ari"
)

// Sound provides the ARI Sound accessors for the native client
type Sound struct {
	client *Client
}

// Get returns a managed handle to a SoundData
func (s *Sound) Get(key *ari.Key) ari.SoundHandle {
	return NewSoundHandle(key, s)
}

// Data returns the details of a given ARI Sound
// Equivalent to GET /sounds/{name}
func (s *Sound) Data(key *ari.Key) (sd *ari.SoundData, err error) {
	sd = &ari.SoundData{}
	name := key.ID
	err = s.client.get("/sounds/"+name, sd)
	if err != nil {
		sd = nil
		err = dataGetError(err, "sound", "%v", name)
	}
	return
}

// List returns available sounds limited by the provided filters.
// valid filters are "lang", "format", and nil (no filter)
// An empty filter returns all available sounds
func (s *Sound) List(filters map[string]string, keyFilter *ari.Key) (sh []*ari.Key, err error) {

	var sounds = []struct {
		Name string `json:"name"`
	}{}

	uri := "/sounds"
	if len(filters) > 0 {
		v := url.Values{}
		for key, val := range filters {
			v.Set(key, val)
		}
		uri += "?" + v.Encode()
	}

	if keyFilter == nil {
		keyFilter = ari.NodeKey(s.client.ApplicationName(), s.client.node)
	}

	err = s.client.get(uri, &sounds)

	// Store whatever we received, even if incomplete or error
	for _, i := range sounds {
		k := ari.NewKey(ari.SoundKey, i.Name, ari.WithApp(s.client.ApplicationName()), ari.WithNode(s.client.node))
		if keyFilter.Match(k) {
			sh = append(sh, k)
		}
	}

	return sh, err
}

// SoundHandle provides a wrapper to a Sound interface for
// operations on a specific Sound
type SoundHandle struct {
	key *ari.Key
	s   *Sound
}

// NewSoundHandle creates a new handle to the sound name
func NewSoundHandle(key *ari.Key, snd *Sound) ari.SoundHandle {
	return &SoundHandle{
		key: key,
		s:   snd,
	}
}

// ID returns the identifier for the sound
func (sh *SoundHandle) ID() string {
	return sh.key.ID
}

// Data retrieves the data for the Sound
func (sh *SoundHandle) Data() (sd *ari.SoundData, err error) {
	sd, err = sh.s.Data(sh.key)
	return sd, err
}
