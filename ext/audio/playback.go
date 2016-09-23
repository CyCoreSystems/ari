package audio

import (
	"github.com/CyCoreSystems/ari"
	"golang.org/x/net/context"
)

// A Playback is a lifecycle managed audio object
type Playback struct {
	startCh chan struct{}
	stopCh  chan struct{}

	handle *ari.PlaybackHandle

	status Status
	err    error

	ctx    context.Context
	cancel context.CancelFunc
}

// Handle returns the ARI reference to the playback object
func (p *Playback) Handle() *ari.PlaybackHandle {
	return p.handle
}

// Started returns the channel that is closed when the playback has started
func (p *Playback) Started() <-chan struct{} {
	return p.startCh
}

// Stopped returns the channel that is closed when the playback has stopped
func (p *Playback) Stopped() <-chan struct{} {
	return p.stopCh
}

// Status returns the current status of the playback
func (p *Playback) Status() Status {
	return p.status
}

// Err returns any accumulated errors during playback
func (p *Playback) Err() error {
	return p.err
}

// Cancel stops the playback
func (p *Playback) Cancel() {
	p.cancel()
}
