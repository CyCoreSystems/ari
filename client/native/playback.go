package native

import "github.com/CyCoreSystems/ari"

type nativePlayback struct {
	conn       *Conn
	subscriber ari.Subscriber
}

func (a *nativePlayback) Get(id string) (ph *ari.PlaybackHandle) {
	ph = ari.NewPlaybackHandle(id, a)
	return
}

// Data returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (a *nativePlayback) Data(id string) (p ari.PlaybackData, err error) {
	err = Get(a.conn, "/playbacks/"+id, &p)
	return
}

// Control allows the user to manipulate an in-process playback.
// TODO: list available operations.
// (Equivalent to POST /playbacks/{playbackID}/control)
func (a *nativePlayback) Control(id string, op string) (err error) {

	//Request structure for controlling playback. Operation is required.
	type request struct {
		Operation string `json:"operation"`
	}

	req := request{op}
	err = Post(a.conn, "/playbacks/"+id+"/control", nil, &req)
	return
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *nativePlayback) Stop(id string) (err error) {
	err = Delete(a.conn, "/playbacks/"+id, nil, "")
	return
}

func (a *nativePlayback) Subscribe(id string, n ...string) ari.Subscription {
	var ns nativeSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	playbackHandle := a.Get(id)

	go func() {
		sub := a.subscriber.Subscribe(n...)
		defer sub.Cancel()
		for {

			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt := <-sub.Events():
				if playbackHandle.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()

	return &ns
}
