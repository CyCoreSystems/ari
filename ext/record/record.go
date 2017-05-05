package record

import (
	"context"
	"fmt"
	"time"

	"sync"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// Options describes a set of recording options for a recording Session
type Options struct {
	// name is the name for the live recording
	name string

	options *ari.RecordingOptions
}

// Apply applies a set of options for the recording Session
func (o *Options) Apply(opts ...OptionFunc) {
	for _, f := range opts {
		f(o)
	}
}

// OptionFunc is a function which applies changes to an Options set
type OptionFunc func(*Options)

// Session desribes the interface to a generic recording session
type Session interface {
	// Done returns a channel which is closed when the session is complete
	Done() <-chan struct{}

	// Err waits for the session to complete, then returns any error encountered during its execution
	Err() error

	// Result waits for the session to complete, then returns the Result
	Result() (*Result, error)

	// Scrap terminates the recording session and throws away the recording.
	Scrap()

	// Stop stops the recording session
	Stop() *Result
}

// Result represents the result of a recording Session.  It provides an interface to disposition the recording.
type Result struct {
	h *ari.LiveRecordingHandle

	// DTMF holds any DTMF digits which are received during the recording session
	DTMF string

	// Duration indicates the duration of the recording
	Duration time.Duration

	// Error holds any error encountered during the recording session
	Error error

	// Status indicates the status of the recording
	Status Status
}

// Delete discards the recording
func (r *Result) Delete() {
	r.h.Scrap()
}

// Record starts a new recording Session
func Record(ctx context.Context, r ari.Recorder, opts ...OptionFunc) Session {
	s := newRecordingSession(opts...)

	var wg sync.WaitGroup
	wg.Add(1)
	go s.record(ctx, r, &wg)
	wg.Wait()

	return s
}

// New creates a new recording Session
func newRecordingSession(opts ...OptionFunc) *recordingSession {
	o := &Options{
		name:    uuid.NewV1().String(),
		options: new(ari.RecordingOptions),
	}

	o.Apply(opts...)

	return &recordingSession{
		cancel:  func() {},
		options: o,
		doneCh:  make(chan struct{}),
		status:  InProgress,
		res:     new(Result),
	}
}

type recordingSession struct {
	cancel context.CancelFunc

	doneCh chan struct{}

	options *Options

	status Status

	res *Result
}

func (s *recordingSession) Done() <-chan struct{} {
	return s.doneCh
}

func (s *recordingSession) Err() error {
	<-s.Done()
	return s.res.Error
}

func (s *recordingSession) Result() (*Result, error) {
	<-s.Done()
	return s.res, s.res.Error
}

func (s *recordingSession) Scrap() {
	s.res.h.Scrap()
}

func (s *recordingSession) Stop() *Result {
	// Signal stop
	s.res.h.Stop()

	// Wait for the stop to complete
	<-s.Done()

	// Return the result
	return s.res
}

func (s *recordingSession) record(ctx context.Context, r ari.Recorder, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	defer cancel()

	h, err := r.StageRecord(s.options.name, s.options.options)
	if err != nil {
		s.status = Failed
		s.res.Error = errors.Wrap(err, "failed to stage recording")
		wg.Done()
		return
	}

	dtmfSub := r.Subscribe(ari.Events.ChannelDtmfReceived)
	hangupSub := r.Subscribe(ari.Events.ChannelDestroyed, ari.Events.ChannelHangupRequest)
	startSub := h.Subscribe(ari.Events.RecordingStarted)
	failedSub := h.Subscribe(ari.Events.RecordingFailed)
	finishedSub := h.Subscribe(ari.Events.RecordingFinished)

	defer func() {
		hangupSub.Cancel()
		failedSub.Cancel()
		startSub.Cancel()
		finishedSub.Cancel()
		dtmfSub.Cancel()
	}()

	wg.Done()

	// Record the duration of the recording
	started := time.Now()
	defer func() {
		s.res.Duration = time.Since(started)
	}()

	if err := h.Exec(); err != nil {
		s.status = Failed
		s.res.Error = err
		return
	}

	go s.waitDtmf(ctx, dtmfSub)

	startTimer := time.NewTimer(RecordingStartTimeout)

	for {
		select {
		case <-ctx.Done():
			s.status = Canceled
			return
		case <-startTimer.C:
			s.status = Failed
			s.res.Error = timeoutErr{"Timeout waiting for recording to start"}
			return
		case e, ok := <-startSub.Events():
			if !ok {
				return
			}
			r := e.(*ari.RecordingStarted).Recording
			if r.Name == s.options.name {
				Logger.Debug("Recording started.")
				startTimer.Stop()
			}
		case e, ok := <-failedSub.Events():
			if !ok {
				return
			}
			r := e.(*ari.RecordingFailed).Recording
			if r.Name == s.options.name {
				s.status = Failed
				s.res.Error = fmt.Errorf("Recording failed: %s", r.Cause)
				return
			}
		case e, ok := <-finishedSub.Events():
			if !ok {
				return
			}
			r := e.(*ari.RecordingFinished).Recording
			if r.Name == s.options.name {
				Logger.Debug("Recording stopped")
				s.status = Finished
				s.res.Duration = time.Duration(r.Duration) * time.Second
				return
			}
		case <-hangupSub.Events():
			s.status = Hangup
			return
		}
	}
}

func (s *recordingSession) waitDtmf(ctx context.Context, dtmfSub ari.Subscription) {
	for {
		select {
		case e, ok := <-dtmfSub.Events():
			if !ok {
				return
			}
			v := e.(*ari.ChannelDtmfReceived)
			s.res.DTMF += v.Digit
		case <-ctx.Done():
			return
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
