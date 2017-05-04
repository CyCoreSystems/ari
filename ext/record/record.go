package record

import (
	"context"
	"fmt"
	"time"

	"sync"

	"github.com/CyCoreSystems/ari"
	uuid "github.com/satori/go.uuid"
)

// RecordingStartTimeout is the amount of time to wait for a recording to start
// before declaring the recording to have failed.
var RecordingStartTimeout = 1 * time.Second

// Options describes a set of recording options for a recording Session
type Options struct {
	// name is the name for the live recording
	name string
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

// Save stores the recording to a Stored Recording, returning a handle to that stored recording.
func (r *Result) Save(name string, rec ari.Recording) (*ari.StoredRecordingHandle, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	// our live recording, once stopped, should be a stored recording.
	k := *r.h.Key() //copy
	k.Kind = ari.StoredRecordingKey

	handle, err := rec.Stored.Copy(&k, name)

	return handle, err
}

// Record starts a new recording Session
func Record(ctx context.Context, r ari.Recorder, opts ...OptionFunc) Session {
	s := newRecordingSession(opts...)

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	var wg sync.WaitGroup
	wg.Add(1)
	go s.record(ctx, r, &wg)
	wg.Wait()

	return s
}

// New creates a new recording Session
func newRecordingSession(opts ...OptionFunc) *recordingSession {
	o := &Options{
		name: uuid.NewV1().String(),
	}

	o.Apply(opts...)

	return &recordingSession{
		options: o,
		doneCh:  make(chan struct{}),
		status:  InProgress,
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
	select {
	case <-s.doneCh:
	}

	return s.res.Error
}

func (s *recordingSession) Result() (*Result, error) {
	select {
	case <-s.doneCh:
	}

	return s.res, s.res.Error
}

func (s *recordingSession) Scrap() {
	s.res.h.Scrap()
}

func (s *recordingSession) Stop() *Result {
	s.res.h.Stop()

	select {
	case <-s.doneCh:
	}

	return s.res
}

func (s *recordingSession) record(ctx context.Context, r ari.Recorder, wg *sync.WaitGroup) {

	s.res = &Result{}

	lhr, err := r.StageRecord(s.options.name, &ari.RecordingOptions{})
	if err != nil {
		s.status = Failed
		s.res.Error = err
		wg.Done()
		return
	}

	name := s.options.name

	dtmfSub := r.Subscribe(ari.Events.ChannelDtmfReceived)
	hangupSub := r.Subscribe(ari.Events.ChannelDestroyed, ari.Events.ChannelHangupRequest)
	startSub := lhr.Subscribe(ari.Events.RecordingStarted)
	failedSub := lhr.Subscribe(ari.Events.RecordingFailed)
	finishedSub := lhr.Subscribe(ari.Events.RecordingFinished)

	defer func() {
		hangupSub.Cancel()
		failedSub.Cancel()
		startSub.Cancel()
		finishedSub.Cancel()
		dtmfSub.Cancel()
	}()

	if err := lhr.Exec(); err != nil {
		s.status = Failed
		s.res.Error = err
		wg.Done()
		return
	}

	wg.Add(1)
	go s.waitDtmf(ctx, dtmfSub, wg)

	startTimer := time.NewTimer(RecordingStartTimeout)

	wg.Done()
	for {
		select {
		case <-ctx.Done():
			s.status = Canceled
			return
		case <-startTimer.C:
			s.status = Failed
			s.res.Error = timeoutErr{"Timeout waiting for recording to start"}
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
				s.status = Failed
				s.res.Error = fmt.Errorf("Recording failed: %s", r.Cause)
				return
			}
		case e := <-finishedSub.Events():
			r := e.(*ari.RecordingFinished).Recording
			if r.Name == name {
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

func (s *recordingSession) waitDtmf(ctx context.Context, dtmfSub ari.Subscription, wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case e, more := <-dtmfSub.Events():
			if !more {
				return
			}
			evt := e.(*ari.ChannelDtmfReceived)
			s.res.DTMF += evt.Digit
		case <-ctx.Done():
			return
		case <-s.doneCh:
			return
		}
	}
}

/*
	Logger.Debug("Starting record", "name", name, "opts", opts)

	// Create recording handle
	h, err := r.StageRecord(name, opts)
	if err != nil {
		rec.err = err
		rec.status = Failed
		close(rec.doneCh)
		return
	}
	rec.handle = h

	go func() {
		defer close(rec.doneCh)

		Logger.Debug("Grabbing subscriptions", "name", name, "opts", opts)

		hangupSub := r.Subscribe(ari.Events.ChannelDestroyed, ari.Events.ChannelHangupRequest)
		defer hangupSub.Cancel()

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
*/

type timeoutErr struct {
	msg string
}

func (err timeoutErr) Error() string {
	return err.msg
}

func (err timeoutErr) Timeout() bool {
	return true
}
