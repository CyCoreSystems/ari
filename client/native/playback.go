package native

import "github.com/CyCoreSystems/ari"

// Playback provides the ARI Playback accessors for the native client
type Playback struct {
	client *Client
}

// Get gets a lazy handle for the given playback identifier
func (a *Playback) Get(key *ari.Key) (ph ari.PlaybackHandle) {
	ph = NewPlaybackHandle(key, a, nil)
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

// PlaybackHandle is the handle for performing playback operations
type PlaybackHandle struct {
	key      *ari.Key
	p        *Playback
	exec     func(pb *PlaybackHandle) error
	executed bool
}

// NewPlaybackHandle builds a handle to the playback id
func NewPlaybackHandle(key *ari.Key, pb *Playback, exec func(pb *PlaybackHandle) error) ari.PlaybackHandle {
	return &PlaybackHandle{
		key:  key,
		p:    pb,
		exec: exec,
	}
}

// ID returns the identifier for the playback
func (ph *PlaybackHandle) ID() string {
	return ph.key.ID
}

// Data gets the playback data
func (ph *PlaybackHandle) Data() (pd *ari.PlaybackData, err error) {
	pd, err = ph.p.Data(ph.key)
	return
}

// Control performs the given operation
func (ph *PlaybackHandle) Control(op string) (err error) {
	err = ph.p.Control(ph.key, op)
	return
}

// Stop stops the playback
func (ph *PlaybackHandle) Stop() (err error) {
	err = ph.p.Stop(ph.key)
	return
}

// Match returns true if the event matches the playback
func (ph *PlaybackHandle) Match(e ari.Event) bool {
	p, ok := e.(ari.PlaybackEvent)
	if !ok {
		return false
	}
	ids := p.GetPlaybackIDs()
	for _, i := range ids {
		if i == ph.ID() {
			return true
		}
	}
	return false
}

// Subscribe subscribes the list of channel events
func (ph *PlaybackHandle) Subscribe(n ...string) ari.Subscription {
	if ph == nil {
		return nil
	}
	return ph.p.Subscribe(ph.key, n...)
}

// Exec executes any staged operations
func (ph *PlaybackHandle) Exec() (err error) {
	if !ph.executed {
		ph.executed = true
		if ph.exec != nil {
			err = ph.exec(ph)
			ph.exec = nil
		}
	}
	return
}
