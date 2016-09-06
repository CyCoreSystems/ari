package record

import (
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"
	v2 "github.com/CyCoreSystems/ari/v2"
	"github.com/pkg/errors"

	"golang.org/x/net/context"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// Record starts a recording on the given Recorder.
func Record(ctx context.Context, bus ari.Subscriber, r Recorder, name string, opts *ari.RecordingOptions) (rec *Recording, err error) {

	// Listen for start, stop, and failed events
	startSub := bus.Subscribe("RecordingStarted")
	defer startSub.Cancel()

	failedSub := bus.Subscribe("RecordingFailed")
	defer failedSub.Cancel()

	finishedSub := bus.Subscribe("RecordingFinished")
	defer finishedSub.Cancel()

	var handle *ari.LiveRecordingHandle

	rec = &Recording{
		Opts:   opts,
		Handle: handle,
		Status: InProgress,
	}

	// Start recording
	handle, err = r.Record(name, opts)
	if err != nil {
		rec.Status = Failed
		return
	}

	// TODO: we have no way to track hangups because we do
	// not have the affiliated channel ID.  We _may_ be able
	// to compare a ChannelHangupRequest event's channel with
	// the LiveRecording's TargetURI, but that will only work
	// for channels.

	// Wait for the recording to start
	startTimer := time.NewTimer(RecordingStartTimeout)
	for {
		select {
		case <-startTimer.C:
			rec.Status = Failed
			err = timeoutErr{"Timeout waiting for recording to start"}
			return
		case <-ctx.Done():
			rec.Status = Canceled
			err = errors.Wrap(ctx.Err(), "Recording canceled")
			return
		case e := <-startSub.C:
			r := e.(*v2.RecordingStarted).Recording
			if r.Name == name {
				Logger.Debug("Recording started.")
				startTimer.Stop()
			}
		case e := <-failedSub.C:
			r := e.(*v2.RecordingFailed).Recording
			if r.Name == name {
				rec.Status = Failed
				err = fmt.Errorf("Recording failed: %s", r.Cause)
				return
			}
		case e := <-finishedSub.C:
			r := e.(*v2.RecordingFinished).Recording
			if r.Name == name {
				Logger.Debug("Recording stopped")
				rec.Status = Finished
				return
			}
		}
	}
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
