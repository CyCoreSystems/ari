package testutils

import "github.com/PolyAI-LDN/ari/v6"

// A RecorderPair is the pair of results returned from a mock Record request
type RecorderPair struct {
	Handle *ari.LiveRecordingHandle
	Err    error
}

// Recorder is the test player that can be primed with sample data
type Recorder struct {
	Next    chan struct{}
	results []RecorderPair
}

// NewRecorder creates a new player
func NewRecorder() *Recorder {
	return &Recorder{
		Next: make(chan struct{}, 10),
	}
}

// Append appends the given Play results
func (r *Recorder) Append(h *ari.LiveRecordingHandle, err error) {
	r.results = append(r.results, RecorderPair{h, err})
}

// Record pops the top results and returns them, as well as triggering player.Next
func (r *Recorder) Record(name string, opts *ari.RecordingOptions) (h *ari.LiveRecordingHandle, err error) {
	h = r.results[0].Handle
	err = r.results[0].Err
	r.results = r.results[1:]

	r.Next <- struct{}{}

	return
}
