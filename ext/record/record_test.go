package record

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/internal/testutils"
	"github.com/golang/mock/gomock"
)

func TestRecordTimeout(t *testing.T) {
	RecordingStartTimeout = 100 * time.Millisecond

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)
	startedSub := mock.NewMockSubscription(ctrl)
	startedCh := make(chan ari.Event)
	startedSub.EXPECT().Events().Return(startedCh)
	startedSub.EXPECT().Cancel()

	finishedSub := mock.NewMockSubscription(ctrl)
	finishedCh := make(chan ari.Event)
	finishedSub.EXPECT().Events().Return(finishedCh)
	finishedSub.EXPECT().Cancel()

	failedSub := mock.NewMockSubscription(ctrl)
	failedCh := make(chan ari.Event)
	failedSub.EXPECT().Events().Return(failedCh)
	failedSub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.RecordingStarted).Return(startedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFinished).Return(finishedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFailed).Return(failedSub)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)
	startedSub := mock.NewMockSubscription(ctrl)
	startedCh := make(chan ari.Event)
	startedSub.EXPECT().Events().MinTimes(1).Return(startedCh)
	startedSub.EXPECT().Cancel()

	finishedSub := mock.NewMockSubscription(ctrl)
	finishedCh := make(chan ari.Event)
	finishedSub.EXPECT().Events().MinTimes(1).Return(finishedCh)
	finishedSub.EXPECT().Cancel()

	failedSub := mock.NewMockSubscription(ctrl)
	failedCh := make(chan ari.Event)
	failedSub.EXPECT().Events().MinTimes(1).Return(failedCh)
	failedSub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.RecordingStarted).Return(startedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFinished).Return(finishedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFailed).Return(failedSub)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	var rec *Recording
	var err error

	rec = Record(bus, recorder, "rc1", nil)

	startedCh <- recordingStarted("rc1")
	finishedCh <- recordingFinished("rc1")

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)
	startedSub := mock.NewMockSubscription(ctrl)
	startedCh := make(chan ari.Event)
	startedSub.EXPECT().Events().Return(startedCh)
	startedSub.EXPECT().Cancel()

	finishedSub := mock.NewMockSubscription(ctrl)
	finishedCh := make(chan ari.Event)
	finishedSub.EXPECT().Events().Return(finishedCh)
	finishedSub.EXPECT().Cancel()

	failedSub := mock.NewMockSubscription(ctrl)
	failedCh := make(chan ari.Event)
	failedSub.EXPECT().Events().Return(failedCh)
	failedSub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.RecordingStarted).Return(startedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFinished).Return(finishedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFailed).Return(failedSub)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	rec := Record(bus, recorder, "rc1", nil)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)
	startedSub := mock.NewMockSubscription(ctrl)
	startedSub.EXPECT().Cancel()

	finishedSub := mock.NewMockSubscription(ctrl)
	finishedSub.EXPECT().Cancel()

	failedSub := mock.NewMockSubscription(ctrl)
	failedSub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.RecordingStarted).Return(startedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFinished).Return(finishedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFailed).Return(failedSub)

	recorder := testutils.NewRecorder()
	recorder.Append(nil, errors.New("Dummy record error"))

	rec := Record(bus, recorder, "rc1", nil)

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bus := mock.NewMockBus(ctrl)

	recorder := testutils.NewRecorder()
	recorder.Append(ari.NewLiveRecordingHandle("rc1", &testRecording{"rc1", false}), nil)

	startedSub := mock.NewMockSubscription(ctrl)
	startedCh := make(chan ari.Event)
	startedSub.EXPECT().Events().Return(startedCh)
	startedSub.EXPECT().Cancel()

	finishedSub := mock.NewMockSubscription(ctrl)
	finishedCh := make(chan ari.Event)
	finishedSub.EXPECT().Events().Return(finishedCh)
	finishedSub.EXPECT().Cancel()

	failedSub := mock.NewMockSubscription(ctrl)
	failedCh := make(chan ari.Event)
	failedSub.EXPECT().Events().Return(failedCh)
	failedSub.EXPECT().Cancel()

	bus.EXPECT().Subscribe(ari.Events.RecordingStarted).Return(startedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFinished).Return(finishedSub)
	bus.EXPECT().Subscribe(ari.Events.RecordingFailed).Return(failedSub)

	rec := Record(bus, recorder, "rc1", nil)

	failedCh <- recordingFailed("rc1")

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

var recordingStarted = func(id string) ari.Event {
	return &ari.RecordingStarted{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "RecordingStarted",
			},
		},
		Recording: ari.LiveRecordingData{
			Name: id,
		},
	}
}

var recordingFinished = func(id string) ari.Event {
	return &ari.RecordingFinished{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "RecordingFinished",
			},
		},
		Recording: ari.LiveRecordingData{
			Name: id,
		},
	}
}

var recordingFailed = func(id string) ari.Event {
	return &ari.RecordingFailed{
		EventData: ari.EventData{
			Message: ari.Message{
				Type: "RecordingFailed",
			},
		},
		Recording: ari.LiveRecordingData{
			Name:  id,
			Cause: "Dummy Failure Error",
		},
	}
}
