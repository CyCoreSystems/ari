package record

import (
	"github.com/CyCoreSystems/ari"
	"golang.org/x/net/context"
)

// A Recording is a lifecycle managed audio recording
type Recording struct {
	doneCh chan struct{}

	data   *ari.LiveRecordingData
	handle *ari.LiveRecordingHandle
	status Status
	err    error

	cancel context.CancelFunc
}

// Done returns a channel that is closed when the recording is done
func (r *Recording) Done() chan struct{} {
	return r.doneCh
}

// Err returns any errors in the recording
func (r *Recording) Err() error {
	return r.err
}

// Status returns the status of the recording
func (r *Recording) Status() Status {
	return r.status
}

// Cancel cancels the recording
func (r *Recording) Cancel() {
	r.cancel()
}

// Handle records the live recording handle
func (r *Recording) Handle() *ari.LiveRecordingHandle {
	return r.handle
}

// Data returns the live recording handle data saved at the end of recording event
func (r *Recording) Data() *ari.LiveRecordingData {
	return r.data
}
