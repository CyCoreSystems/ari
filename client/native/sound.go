package native

import (
	"errors"
	"net/url"

	"github.com/PolyAI-LDN/ari/v6"
)

// Sound provides the ARI Sound accessors for the native client
type Sound struct {
	client *Client
}

// Data returns the details of a given ARI Sound
// Equivalent to GET /sounds/{name}
func (s *Sound) Data(key *ari.Key) (*ari.SoundData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("sound key not supplied")
	}

	data := new(ari.SoundData)
	if err := s.client.get("/sounds/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "sound", "%v", key.ID)
	}

	data.Key = s.client.stamp(key)

	return data, nil
}

// List returns available sounds limited by the provided filters.
// valid filters are "lang", "format", and nil (no filter)
// An empty filter returns all available sounds
func (s *Sound) List(filters map[string]string, keyFilter *ari.Key) (sh []*ari.Key, err error) {
	sounds := []struct {
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
		keyFilter = s.client.stamp(ari.NewKey(ari.SoundKey, ""))
	}

	err = s.client.get(uri, &sounds)
	if err != nil {
		return nil, err
	}

	// Store whatever we received, even if incomplete or error
	for _, i := range sounds {
		k := s.client.stamp(ari.NewKey(ari.SoundKey, i.Name))
		if keyFilter.Match(k) {
			sh = append(sh, k)
		}
	}

	return sh, err
}
