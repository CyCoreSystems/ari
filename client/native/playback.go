package native

import "github.com/CyCoreSystems/ari"

// Playback provides the ARI Playback accessors for the native client
type Playback struct {
	client *Client
}

// Get gets a lazy handle for the given playback identifier
func (a *Playback) Get(id string) (ph ari.PlaybackHandle) {
	ph = NewPlaybackHandle(id, a, nil)
	return
}

// Data returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (a *Playback) Data(id string) (p *ari.PlaybackData, err error) {
	p = &ari.PlaybackData{}
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
func (a *Playback) Control(id string, op string) (err error) {

	//Request structure for controlling playback. Operation is required.
	type request struct {
		Operation string `json:"operation"`
	}

	req := request{op}
	err = a.client.post("/playbacks/"+id+"/control", nil, &req)
	return
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *Playback) Stop(id string) (err error) {
	err = a.client.del("/playbacks/"+id, nil, "")
	return
}

// Subscribe listens for ARI events for the given playback entity
func (a *Playback) Subscribe(id string, n ...string) ari.Subscription {
	inSub := a.client.Bus().Subscribe(n...)
	outSub := newSubscription()

	go func() {
		defer inSub.Cancel()

		h := a.Get(id)

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
	id       string
	p        *Playback
	exec     func(pb *PlaybackHandle) error
	executed bool
}

// NewPlaybackHandle builds a handle to the playback id
func NewPlaybackHandle(id string, pb *Playback, exec func(pb *PlaybackHandle) error) ari.PlaybackHandle {
	return &PlaybackHandle{
		id:   id,
		p:    pb,
		exec: exec,
	}
}

// ID returns the identifier for the playback
func (ph *PlaybackHandle) ID() string {
	return ph.id
}

// Data gets the playback data
func (ph *PlaybackHandle) Data() (pd *ari.PlaybackData, err error) {
	pd, err = ph.p.Data(ph.id)
	return
}

// Control performs the given operation
func (ph *PlaybackHandle) Control(op string) (err error) {
	err = ph.p.Control(ph.id, op)
	return
}

// Stop stops the playback
func (ph *PlaybackHandle) Stop() (err error) {
	err = ph.p.Stop(ph.id)
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
	return ph.p.Subscribe(ph.id, n...)
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
