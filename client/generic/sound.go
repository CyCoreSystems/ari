package generic

import (
	"net/url"

	"github.com/CyCoreSystems/ari"
)

type Sound struct {
	Conn Conn
}

// Get returns a managed handle to a SoundData
func (s *Sound) Get(name string) *ari.SoundHandle {
	return ari.NewSoundHandle(name, s)
}

// Data returns the details of a given ARI Sound
// Equivalent to GET /sounds/{name}
func (s *Sound) Data(name string) (sd ari.SoundData, err error) {
	err = s.Conn.Get("/sounds/%s", []interface{}{name}, &sd)
	return sd, err
}

// List returns available sounds limited by the provided filters.
// valid filters are "lang", "format", and nil (no filter)
// An empty filter returns all available sounds
func (s *Sound) List(filters map[string]string) (sh []*ari.SoundHandle, err error) {

	var sounds = []struct {
		Name string `json:"name"`
	}{}

	filter := ""
	if len(filters) > 0 {
		v := url.Values{}
		for key, val := range filters {
			v.Set(key, val)
		}
		filter = "?" + v.Encode()
	}

	err = s.Conn.Get("/sounds%s", []interface{}{filter}, &sounds)

	// Store whatever we received, even if incomplete or error
	for _, i := range sounds {
		sh = append(sh, s.Get(i.Name))
	}

	return sh, err
}
