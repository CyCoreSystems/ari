package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

// Playback provides the ARI Playback accessors for the native client
type Playback struct {
	client *Client
}

// Get gets a lazy handle for the given playback identifier
func (a *Playback) Get(key *ari.Key) (ph *ari.PlaybackHandle) {
	ph = ari.NewPlaybackHandle(key, a, nil)
	return
}

// Data returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (a *Playback) Data(key *ari.Key) (*ari.PlaybackData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("playback key not supplied")
	}

	var data = new(ari.PlaybackData)
	if err := a.client.get("/playbacks/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "playback", "%v", key.ID)
	}

	data.Key = a.client.stamp(key)
	return data, nil
}

// Control allows the user to manipulate an in-process playback.
// TODO: list available operations.
// (Equivalent to POST /playbacks/{playbackID}/control)
func (a *Playback) Control(key *ari.Key, op string) (err error) {
	req := struct {
		Operation string `json:"operation"`
	}{
		Operation: op,
	}
	return a.client.post("/playbacks/"+key.ID+"/control", nil, &req)
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *Playback) Stop(key *ari.Key) (err error) {
	return a.client.del("/playbacks/"+key.ID, nil, "")
}

// Subscribe listens for ARI events for the given playback entity
func (a *Playback) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	inSub := a.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		for {
			select {
			case <-outSub.closedChan:
				return
			case e, ok := <-inSub.Events():
				if !ok {
					return
				}
				for _, k := range e.Keys() {
					if k.Match(key) {
						outSub.events <- e
						break
					}
				}
			}
		}
	}()

	return outSub
}
