package native

import (
	"errors"

	"github.com/PolyAI-LDN/ari/v6"
)

// Playback provides the ARI Playback accessors for the native client
type Playback struct {
	client *Client
}

// Get gets a lazy handle for the given playback identifier
func (a *Playback) Get(key *ari.Key) *ari.PlaybackHandle {
	return ari.NewPlaybackHandle(a.client.stamp(key), a, nil)
}

// Data returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (a *Playback) Data(key *ari.Key) (*ari.PlaybackData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("playback key not supplied")
	}

	data := new(ari.PlaybackData)
	if err := a.client.get("/playbacks/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "playback", "%v", key.ID)
	}

	data.Key = a.client.stamp(key)

	return data, nil
}

// Control performs the given operation on the current playback.  Available operations are:
//   - restart
//   - pause
//   - unpause
//   - reverse
//   - forward
func (a *Playback) Control(key *ari.Key, op string) error {
	req := struct {
		Operation string `json:"operation"`
	}{
		Operation: op,
	}

	return a.client.post("/playbacks/"+key.ID+"/control", nil, &req)
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *Playback) Stop(key *ari.Key) error {
	return a.client.del("/playbacks/"+key.ID, nil, "")
}

// Subscribe listens for ARI events for the given playback entity
func (a *Playback) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	return a.client.Bus().Subscribe(key, n...)
}
