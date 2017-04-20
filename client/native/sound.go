package native

import (
	"errors"
	"net/url"

	"github.com/CyCoreSystems/ari"
)

// Sound provides the ARI Sound accessors for the native client
type Sound struct {
	client *Client
}

// Get returns a managed handle to a SoundData
func (s *Sound) Get(key *ari.Key) *ari.SoundHandle {
	return ari.NewSoundHandle(key, s)
}

// Data returns the details of a given ARI Sound
// Equivalent to GET /sounds/{name}
func (s *Sound) Data(key *ari.Key) (*ari.SoundData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("sound key not supplied")
	}

	var data = new(ari.SoundData)
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
