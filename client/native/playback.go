package native

import "github.com/CyCoreSystems/ari"

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
func (a *Playback) Data(key *ari.Key) (p *ari.PlaybackData, err error) {
	p = &ari.PlaybackData{}
	id := key.ID
	err = a.client.get("/playbacks/"+id, p)
	if err != nil {
		p = nil
		err = dataGetError(err, "playback", "%v", id)
	}
	return
}

// Control allows the user to manipulate an in-process playback.
// TODO: list available operations.
// (Equivalent to POST /playbacks/{playbackID}/control)
func (a *Playback) Control(key *ari.Key, op string) (err error) {

	//Request structure for controlling playback. Operation is required.
	type request struct {
		Operation string `json:"operation"`
	}
	id := key.ID
	req := request{op}
	err = a.client.post("/playbacks/"+id+"/control", nil, &req)
	return
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *Playback) Stop(key *ari.Key) (err error) {
	id := key.ID
	err = a.client.del("/playbacks/"+id, nil, "")
	return
}

// Subscribe listens for ARI events for the given playback entity
func (a *Playback) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	inSub := a.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		h := a.Get(key)

		for {
			select {
			case <-outSub.closedChan:
				return
			case e, ok := <-inSub.Events():
				if !ok {
					return
				}
				if h.Match(e) {
					outSub.events <- e
				}
			}
		}
	}()

	return outSub
}
