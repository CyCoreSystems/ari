package ari

import "net/url"

// Sound describes a media file which may be played back
type Sound struct {
	Formats []FormatLangPair `json:"formats"`
	Id      string           `json:"id"`
	Text    string           `json:"text,omitempty"`
}

// FormatLangPair describes the format and language of a sound file
type FormatLangPair struct {
	Format   string `json:"format"`
	Language string `json:"language"`
}

//ListSounds returns a list of (all) the available sounds
func (c *Client) ListSounds(filters map[string]string) ([]Sound, error) {
	var m []Sound
	err := c.Get("/sounds", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// ListSoundsFiltered lists sounds limited by the provided filters
// valid filters are "lang" and "format"
func (c *Client) ListSoundsFiltered(filters map[string]string) ([]Sound, error) {
	var m []Sound
	uri := "/sounds"
	if len(filters) > 0 {
		v := url.Values{}
		for key, val := range filters {
			v.Set(key, val)
		}
		uri += "?" + v.Encode()
	}
	err := c.Get(uri, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Get a sound's details
//Equivalent to GET /sounds/{soundId}
func (c *Client) GetSound(soundId string) (Sound, error) {
	var m Sound
	err := c.Get("/sounds/"+soundId, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
