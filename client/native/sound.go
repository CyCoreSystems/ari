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
func (s *Sound) Get(name string) *ari.SoundHandle {
	return ari.NewSoundHandle(name, s)
}

// Data returns the details of a given ARI Sound
// Equivalent to GET /sounds/{name}
func (s *Sound) Data(name string) (sd *ari.SoundData, err error) {
	sd = &ari.SoundData{}
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
func (s *Sound) List(filters map[string]string) (sh []*ari.SoundHandle, err error) {

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

	err = s.client.get(uri, &sounds)

	// Store whatever we received, even if incomplete or error
	for _, i := range sounds {
		sh = append(sh, s.Get(i.Name))
	}

	return sh, err
}
