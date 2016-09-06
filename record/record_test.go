package record

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/internal/testutils"
	v2 "github.com/CyCoreSystems/ari/v2"
)

func TestRecordTimeout(t *testing.T) {

	bus := testutils.NewDelayedBus(1 * time.Millisecond)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	rec := Record(bus, recorder, "name1", nil)
	<-rec.Done()

	err := rec.Err()

	if !isTimeout(err) {
		t.Errorf("Expected timeout, got '%v'", err)
	}

	if err == nil || err.Error() != "Timeout waiting for recording to start" {
		t.Errorf("Expected timeout waiting for recording to start, got '%v'", err)
	}

	if rec.Status() != Failed {
		t.Errorf("Expected recording status to be Timeout, was '%v'", rec.Status())
	}

}

func TestRecord(t *testing.T) {

	bus := testutils.NewDelayedBus(1 * time.Millisecond)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	exp := bus.Expect("RecordingStarted")
	exp2 := bus.Expect("RecordingFinished")

	var rec *Recording
	var err error

	rec = Record(bus, recorder, "rc1", nil)

	select {
	case <-exp:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingStarted' subscription")
	}

	select {
	case <-exp2:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingFinished' subscription")
	}

	bus.Send(recordingStarted("rc1"))
	bus.Send(recordingFinished("rc1"))

	<-rec.Done()

	err = rec.Err()

	if err != nil {
		t.Errorf("Unexpected err: '%v'", err)
	}

	if rec.Status() != Finished {
		t.Errorf("Expected recording status to be Finished, was '%v'", rec.Status())
	}
}

func TestRecordCancel(t *testing.T) {

	bus := testutils.NewDelayedBus(1 * time.Millisecond)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	exp := bus.Expect("RecordingStarted")
	exp2 := bus.Expect("RecordingFinished")

	rec := Record(bus, recorder, "rc1", nil)

	select {
	case <-exp:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingStarted' subscription")
	}

	select {
	case <-exp2:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingFinished' subscription")
	}

	rec.Cancel()

	<-rec.Done()

	err := rec.Err()
	if err != nil {
		t.Errorf("Unexpected error: '%v'", err)
	}

	if rec.Status() != Canceled {
		t.Errorf("Expected recording status to be Canceled, was '%v'", rec.Status())
	}
}

func TestRecordFailOnRecord(t *testing.T) {

	bus := testutils.NewDelayedBus(1 * time.Millisecond)

	recorder := testutils.NewRecorder()
	recorder.Append(nil, errors.New("Dummy record error"))

	exp := bus.Expect("RecordingStarted")
	exp2 := bus.Expect("RecordingFinished")

	rec := Record(bus, recorder, "rc1", nil)

	select {
	case <-exp:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingStarted' subscription")
	}

	select {
	case <-exp2:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingFinished' subscription")
	}

	<-rec.Done()

	err := rec.Err()

	if err == nil || err.Error() != "Dummy record error" {
		t.Errorf("Expected error 'Dummy record error', got: '%v'", err)
	}

	if rec.Status() != Failed {
		t.Errorf("Expected recording status to be Failed, was '%v'", rec.Status())
	}
}

func TestRecordFailEvent(t *testing.T) {

	RecordingStartTimeout = 10 * time.Second

	bus := testutils.NewDelayedBus(1 * time.Millisecond)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	exp := bus.Expect("RecordingStarted")
	exp2 := bus.Expect("RecordingFinished")
	exp3 := bus.Expect("RecordingFailed")

	rec := Record(bus, recorder, "rc1", nil)

	select {
	case <-exp:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingStarted' subscription")
	}

	select {
	case <-exp2:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingFinished' subscription")
	}

	select {
	case <-exp3:
	case <-time.After(10 * time.Second):
		t.Errorf("Expected 'RecordingFailed' subscription")
	}

	bus.Send(recordingFailed("rc1"))

	<-rec.Done()

	err := rec.Err()

	if err == nil || err.Error() != "Recording failed: Dummy Failure Error" {
		t.Errorf("Expected error 'Recording failed: Dummy Failure Error', got: '%v'", err)
	}

	if rec.Status() != Failed {
		t.Errorf("Expected recording status to be Failed, was '%v'", rec.Status())
	}
}

type testRecording struct {
	id       string
	failData bool
}

func (tr *testRecording) Get(name string) *ari.LiveRecordingHandle {
	panic("not implemented")
}

func (tr *testRecording) Data(name string) (ari.LiveRecordingData, error) {
	panic("not implemented")
}

func (tr *testRecording) Stop(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Pause(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Resume(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Mute(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Unmute(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Delete(name string) error {
	panic("not implemented")
}

func (tr *testRecording) Scrap(name string) error {
	panic("not implemented")
}

func isTimeout(err error) bool {

	type timeout interface {
		Timeout() bool
	}

	te, ok := err.(timeout)
	return ok && te.Timeout()
}

var recordingStarted = func(id string) v2.Eventer {
	return &v2.RecordingStarted{
		Event: v2.Event{
			Message: v2.Message{
				Type: "RecordingStarted",
			},
		},
		Recording: v2.LiveRecording{
			Name: id,
		},
	}
}

var recordingFinished = func(id string) v2.Eventer {
	return &v2.RecordingFinished{
		Event: v2.Event{
			Message: v2.Message{
				Type: "RecordingFinished",
			},
		},
		Recording: v2.LiveRecording{
			Name: id,
		},
	}
}

var recordingFailed = func(id string) v2.Eventer {
	return &v2.RecordingFailed{
		Event: v2.Event{
			Message: v2.Message{
				Type: "RecordingFailed",
			},
		},
		Recording: v2.LiveRecording{
			Name:  id,
			Cause: "Dummy Failure Error",
		},
	}
}
