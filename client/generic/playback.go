package generic

import "github.com/CyCoreSystems/ari"

type Playback struct {
	Conn Conn
}

func (a *Playback) Get(id string) (ph *ari.PlaybackHandle) {
	ph = ari.NewPlaybackHandle(id, a)
	return
}

// Data returns a playback's details.
// (Equivalent to GET /playbacks/{playbackID})
func (a *Playback) Data(id string) (p ari.PlaybackData, err error) {
	err = a.Conn.Get("/playbacks/%s", []interface{}{id}, &p)
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
	err = a.Conn.Post("/playbacks/%s/control", []interface{}{id}, nil, &req)
	return
}

// Stop stops a playback session.
// (Equivalent to DELETE /playbacks/{playbackID})
func (a *Playback) Stop(id string) (err error) {
	err = a.Conn.Delete("/playbacks/%s", []interface{}{id}, nil, "")
	return
}
