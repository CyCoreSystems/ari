package record

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// Record starts a recording on the given Recorder.
func Record(bus ari.Subscriber, r Recorder, name string, opts *ari.RecordingOptions) (rec *Recording) {

	Logger.Debug("Starting record", "name", name, "opts", opts)

	rec = &Recording{}
	rec.doneCh = make(chan struct{})
	rec.status = InProgress

	ctx, cancel := context.WithCancel(context.Background())
	rec.cancel = cancel

	// TODO: we have no way to track hangups because we do
	// not have the affiliated channel ID.  We _may_ be able
	// to compare a ChannelHangupRequest event's channel with
	// the LiveRecording's TargetURI, but that will only work
	// for channels.

	go func() {

		Logger.Debug("Grabbing subscriptions", "name", name, "opts", opts)

		// Listen for start, stop, and failed events
		startSub := bus.Subscribe(ari.Events.RecordingStarted)
		defer startSub.Cancel()

		Logger.Debug("Grabbing subscriptions", "name", name, "opts", opts)

		failedSub := bus.Subscribe(ari.Events.RecordingFailed)
		defer failedSub.Cancel()

		Logger.Debug("Grabbing subscriptions", "name", name, "opts", opts)

		finishedSub := bus.Subscribe(ari.Events.RecordingFinished)
		defer finishedSub.Cancel()

		Logger.Debug("Calling record on recorder", "name", name, "opts", opts)

		// Start recording
		handle, err := r.Record(name, opts)
		if err != nil {
			rec.err = err
			rec.status = Failed
			close(rec.doneCh)
			return
		}

		defer close(rec.doneCh)

		rec.handle = handle

		Logger.Debug("Starting record event loop", "name", name, "opts", opts)

		// Wait for the recording to start
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
