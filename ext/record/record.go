package record

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sync"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (

	// RecordingStartTimeout is the amount of time to wait for a recording to start
	// before declaring the recording to have failed.
	RecordingStartTimeout = 1 * time.Second

	// DefaultMaximumDuration is the default maximum amount of time a recording
	// should be allowed to continue before being terminated.
	DefaultMaximumDuration = 24 * time.Hour

	// DefaultMaximumSilence is the default maximum amount of time silence may be
	// detected before terminating the recording.
	DefaultMaximumSilence = 5 * time.Minute

	// ShutdownGracePeriod is the amount of time to allow a Stop transaction to
	// complete before shutting down the session anyway.
	ShutdownGracePeriod = 3 * time.Second
)

// Options describes a set of recording options for a recording Session
type Options struct {
	// name is the name for the live recording
	name string

	format string

	maxDuration time.Duration

	maxSilence time.Duration

	ifExists string

	beep bool

	terminateOn string
}

func defaultOptions() *Options {
	return &Options{
		beep:        false,
		format:      "wav",
		ifExists:    "fail",
		maxDuration: DefaultMaximumDuration,
		maxSilence:  DefaultMaximumSilence,
		name:        uuid.NewV1().String(),
		terminateOn: "none",
	}
}

func (o *Options) toRecordingOptions() *ari.RecordingOptions {
	return &ari.RecordingOptions{
		Beep:        o.beep,
		Format:      o.format,
		Exists:      o.ifExists,
		MaxDuration: o.maxDuration,
		MaxSilence:  o.maxSilence,
		Terminate:   o.terminateOn,
	}
}

// Apply applies a set of options for the recording Session
func (o *Options) Apply(opts ...OptionFunc) {
	for _, f := range opts {
		f(o)
	}
}

// OptionFunc is a function which applies changes to an Options set
type OptionFunc func(*Options)

// Beep indicates that a beep should be played to signal the start of recording
func Beep() OptionFunc {
	return func(o *Options) {
		o.beep = true
	}
}

// Format configures the file format to be used to store the recording
func Format(format string) OptionFunc {
	return func(o *Options) {
		o.format = format
	}
}

// IfExists configures the behaviour of the recording if the file to be
// recorded already exists.
//
// Valid options are:  "fail" (default), "overwrite", and "append".
func IfExists(action string) OptionFunc {
	return func(o *Options) {
		o.ifExists = action
	}
}

// MaxDuration sets the maximum duration to allow for the recording.  After
// this amount of time, the recording will be automatically Finished.
//
// A setting of 0 disables the limit.
func MaxDuration(max time.Duration) OptionFunc {
	return func(o *Options) {
		o.maxDuration = max
	}
}

// MaxSilence sets the amount of time a block of silence is allowed to become
// before the recording should be declared Finished.
//
// A setting of 0 disables silence detection.
func MaxSilence(max time.Duration) OptionFunc {
	return func(o *Options) {
		o.maxSilence = max
	}
}

// Name configures the recording to use the provided name
func Name(name string) OptionFunc {
	return func(o *Options) {
		o.name = name
	}
}

// TerminateOn configures the DTMF which, if received, will terminate the
// recording.
//
// Valid values are "none" (default), "any", "*", and "#".
func TerminateOn(dtmf string) OptionFunc {
	return func(o *Options) {
		o.terminateOn = dtmf
	}
}

// Session desribes the interface to a generic recording session
type Session interface {
	// Done returns a channel which is closed when the session is complete
	Done() <-chan struct{}

	// Err waits for the session to complete, then returns any error encountered
	// during its execution
	Err() error

	// Pause temporarily stops the recording session without ending the session
	Pause() error

	// Result waits for the session to complete, then returns the Result
	Result() (*Result, error)

	// Resume restarts a paused recording session
	Resume() error

	// Scrap terminates the recording session and throws away the recording.
	Scrap()

	// Stop stops the recording session
	Stop() *Result
}

// Result represents the result of a recording Session.  It provides an
// interface to disposition the recording.
type Result struct {
	h *ari.StoredRecordingHandle

	// Data holds the final data for the LiveRecording, if it was successful
	Data *ari.LiveRecordingData

	// DTMF holds any DTMF digits which are received during the recording session
	DTMF string

	// Duration indicates the duration of the recording
	Duration time.Duration

	// Error holds any error encountered during the recording session
	Error error

	// Hangup indicates that the Recorder disappeared (due to hangup or
	// destruction) during or after the recording.
	Hangup bool

	overwrite bool
}

// Delete discards the recording
func (r *Result) Delete() error {
	if r.h == nil {
		return errors.New("no stored recording handle available")
	}
	return r.h.Delete()
}

// Save stores the recording to the given name
func (r *Result) Save(name string) error {
	if name == "" {
		// no name indicates the default, which is where it already is
		return nil
	}

	if r.h == nil {
		return errors.New("no stored recording handle available")
	}

	// Copy the recording to the desired name
	destH, err := r.h.Copy(name)
	if err != nil {
		if !strings.Contains(err.Error(), "409 Conflict") || !r.overwrite {
			return errors.Wrapf(err, "failed to copy recording (%s)", r.h.ID())
		}

		// we are set to overwrite, so delete the previous recording
		Logger.Debug("overwriting previous recording")
		err = destH.Delete()
		if err != nil {
			return errors.Wrap(err, "failed to remove previous destination recording")
		}
		_, err = r.h.Copy(name)
		if err != nil {
			return errors.Wrap(err, "failed to copy recording")
		}
	}

	// Delete the original
	err = r.h.Delete()
	if err != nil {
		return errors.Wrap(err, "failed to remove temporary recording after copy")
	}

	return nil
}

// URI returns the AudioURI to play the recording
func (r *Result) URI() string {
	return "recording:" + r.h.ID()
}

// Record starts a new recording Session
func Record(ctx context.Context, r ari.Recorder, opts ...OptionFunc) Session {
	s := newRecordingSession(opts...)

	var wg sync.WaitGroup
	wg.Add(1)
	go s.record(ctx, r, &wg)
	wg.Wait()

	Logger.Debug("returned from internal recording start")
	return s
}

// New creates a new recording Session
func newRecordingSession(opts ...OptionFunc) *recordingSession {
	o := defaultOptions()
	o.Apply(opts...)

	s := &recordingSession{
		cancel:  func() {},
		options: o,
		doneCh:  make(chan struct{}),
		res:     new(Result),
	}

	// If the recording options declare that we should overwrite,
	// carry that over to the copy destination.
	if o.ifExists == "overwrite" {
		s.res.overwrite = true
	}

	return s
}

type recordingSession struct {
	h *ari.LiveRecordingHandle

	cancel context.CancelFunc

	doneCh chan struct{}

	mu sync.Mutex

	options *Options

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

func (s *recordingSession) Pause() error {
	return s.h.Pause()
}

func (s *recordingSession) Resume() error {
	return s.h.Resume()
}

func (s *recordingSession) Scrap() {
	s.h.Scrap()
}

func (s *recordingSession) Stop() *Result {
	// Signal stop
	s.res.Error = s.h.Stop()

	// If we successfully signaled a stop, Wait for the stop to complete
	if s.res.Error == nil {
		select {
		case <-s.Done():
		case <-time.After(ShutdownGracePeriod):
		}
	}

	// Shut down the session
	if s.cancel != nil {
		s.cancel()
	}

	// Return the result
	return s.res
}

// nolint: gocyclo
func (s *recordingSession) record(ctx context.Context, r ari.Recorder, wg *sync.WaitGroup) {
	defer close(s.doneCh)

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	defer cancel()

	var err error
	s.h, err = r.StageRecord(s.options.name, s.options.toRecordingOptions())
	if err != nil {
		s.res.Error = errors.Wrap(err, "failed to stage recording")
		wg.Done()
		return
	}

	// Store the eventual StoredRecording handle to the Result
	s.res.h = s.h.Stored()

	dtmfSub := r.Subscribe(ari.Events.ChannelDtmfReceived)
	hangupSub := r.Subscribe(ari.Events.ChannelDestroyed, ari.Events.ChannelHangupRequest, ari.Events.BridgeDestroyed)
	startSub := s.h.Subscribe(ari.Events.RecordingStarted)
	failedSub := s.h.Subscribe(ari.Events.RecordingFailed)
	finishedSub := s.h.Subscribe(ari.Events.RecordingFinished)

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
		Logger.Debug("recording duration", "duration", s.res.Duration)
	}()

	// Record any DTMF received during the recording
	go s.collectDtmf(ctx, dtmfSub)

	// Record hangup or destruction of our Recorder
	go s.watchHangup(ctx, hangupSub)

	// Start recording
	Logger.Debug("starting recording")
	if err := s.h.Exec(); err != nil {
		s.res.Error = err
		return
	}

	// Time the recording
	startTimer := time.NewTimer(RecordingStartTimeout)

	for {
		select {
		case <-ctx.Done():
			s.Stop()
			return
		case <-startTimer.C:
			Logger.Debug("timeout waiting to start recording")
			s.res.Error = timeoutErr{"Timeout waiting for recording to start"}
			return
		case _, ok := <-startSub.Events():
			if !ok {
				return
			}
			Logger.Debug("recording started")
			startTimer.Stop()
		case e, ok := <-failedSub.Events():
			if !ok {
				return
			}
			Logger.Debug("recording failed")
			r := e.(*ari.RecordingFailed).Recording
			s.res.Data = &r
			s.res.Error = fmt.Errorf("Recording failed: %s", r.Cause)
			return
		case e, ok := <-finishedSub.Events():
			if !ok {
				return
			}
			r := e.(*ari.RecordingFinished).Recording
			Logger.Debug("recording finished")
			s.res.Data = &r
			return
		}
	}
}

func (s *recordingSession) collectDtmf(ctx context.Context, dtmfSub ari.Subscription) {
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

func (s *recordingSession) watchHangup(ctx context.Context, hangupSub ari.Subscription) {
	select {
	case <-hangupSub.Events():
		s.res.Hangup = true
	case <-ctx.Done():
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
