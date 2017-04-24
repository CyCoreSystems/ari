package record

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// Record starts a recording on the given Recorder.
// TODO: simplify
// nolint:gocyclo
func Record(r Recorder, name string, opts *ari.RecordingOptions) (rec *Recording) {

	Logger.Debug("Starting record", "name", name, "opts", opts)

	rec = &Recording{}
	rec.doneCh = make(chan struct{})
	rec.status = InProgress

	ctx, cancel := context.WithCancel(context.Background())
	rec.cancel = cancel

	// Create recording handle
	h, err := r.StageRecord(name, opts)
	if err != nil {
		rec.err = err
		rec.status = Failed
		close(rec.doneCh)
		return
	}
	rec.handle = h

	// TODO: we have no way to track hangups because we do
	// not have the affiliated channel ID.  We _may_ be able
	// to compare a ChannelHangupRequest event's channel with
	// the LiveRecording's TargetURI, but that will only work
	// for channels.

	go func() {
		defer close(rec.doneCh)

		Logger.Debug("Grabbing subscriptions", "name", name, "opts", opts)

		// Listen for start, stop, and failed events
		startSub := h.Subscribe(ari.Events.RecordingStarted)
		defer startSub.Cancel()

		failedSub := h.Subscribe(ari.Events.RecordingFailed)
		defer failedSub.Cancel()

		finishedSub := h.Subscribe(ari.Events.RecordingFinished)
		defer finishedSub.Cancel()

		Logger.Debug("Starting recording", "name", name, "opts", opts)
		err := h.Exec()
		if err != nil {
			rec.status = Failed
			rec.err = errors.Wrap(err, "failed to start recording")
			return
		}

		// Wait for the recording to start
		Logger.Debug("Starting record event loop", "name", name, "opts", opts)
		startTimer := time.NewTimer(RecordingStartTimeout)
		for {
			select {
			case <-ctx.Done():
				rec.status = Canceled
				return
			case <-startTimer.C:
				rec.status = Failed
				rec.err = timeoutErr{"Timeout waiting for recording to start"}
				return
			case e := <-startSub.Events():
				r := e.(*ari.RecordingStarted).Recording
				if r.Name == name {
					Logger.Debug("Recording started.")
					startTimer.Stop()
				}
			case e := <-failedSub.Events():
				r := e.(*ari.RecordingFailed).Recording
				if r.Name == name {
					rec.status = Failed
					rec.err = fmt.Errorf("Recording failed: %s", r.Cause)
					return
				}
			case e := <-finishedSub.Events():
				r := e.(*ari.RecordingFinished).Recording
				if r.Name == name {
					Logger.Debug("Recording stopped")
					rec.status = Finished
					rec.data = &r
					return
				}
			}
		}
	}()

	return
}

type timeoutErr struct {
	msg string
}

func (err timeoutErr) Error() string {
	return err.msg
}

func (err timeoutErr) Timeout() bool {
	return true
}
