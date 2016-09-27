package natsgw

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/client/nc"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestStoredRecordingList(t *testing.T) {

	//TODO: embed nats?

	bin, err := exec.LookPath("gnatsd")
	if err != nil {
		t.Skip("No gnatsd binary in PATH, skipping")
	}

	cmd := exec.Command(bin, "-p", "4333")
	if err := cmd.Start(); err != nil {
		t.Errorf("Unable to run gnatsd: '%v'", err)
		return
	}

	defer func() {
		cmd.Process.Signal(syscall.SIGTERM)
		cmd.Wait()
	}()

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoredRecording := mock.NewMockStoredRecording(ctrl)

	gomock.InOrder(
		mockStoredRecording.EXPECT().List().Return([]*ari.StoredRecordingHandle{
			ari.NewStoredRecordingHandle("sr1", mockStoredRecording),
			ari.NewStoredRecordingHandle("sr2", mockStoredRecording),
		}, nil),
		mockStoredRecording.EXPECT().List().Return([]*ari.StoredRecordingHandle{}, errors.New("Failed to list stored recordings")),
	)

	cl := &ari.Client{
		Recording: &ari.Recording{
			Stored: mockStoredRecording,
		},
	}

	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	s.Start()
	defer s.Close()

	natsClient, err := newNatsClient("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("newNatsClient(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{
		lx, err := natsClient.Recording.Stored.List()

		failed = err != nil
		failed = failed || len(lx) != 2
		if failed {
			t.Errorf("nc.Recording.Stored.List() => '%v', '%v', expected '%v', '%v'",
				lx, err, "[sr1,sr2]", nil)
		}
	}
	{
		lx, err := natsClient.Recording.Stored.List()

		failed = err == nil || errors.Cause(err).Error() != "Failed to list stored recordings"
		failed = failed || len(lx) != 0
		if failed {
			t.Errorf("nc.Recording.Stored.List() => '%v', '%v', expected '%v', '%v'",
				lx, err, "[]", "Failed to list stored recordings")
		}
	}
}

func TestStoredRecordingData(t *testing.T) {

	//TODO: embed nats?

	bin, err := exec.LookPath("gnatsd")
	if err != nil {
		t.Skip("No gnatsd binary in PATH, skipping")
	}

	cmd := exec.Command(bin, "-p", "4333")
	if err := cmd.Start(); err != nil {
		t.Errorf("Unable to run gnatsd: '%v'", err)
		return
	}

	defer func() {
		cmd.Process.Signal(syscall.SIGTERM)
		cmd.Wait()
	}()

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoredRecording := mock.NewMockStoredRecording(ctrl)

	var storedRecordingData ari.StoredRecordingData
	storedRecordingData.Name = "sr1"

	gomock.InOrder(
		mockStoredRecording.EXPECT().Data("sr1").Return(storedRecordingData, nil),
		mockStoredRecording.EXPECT().Data("sr2").Return(ari.StoredRecordingData{}, errors.New("Failed to get stored recording")),
	)

	cl := &ari.Client{
		Recording: &ari.Recording{
			Stored: mockStoredRecording,
		},
	}

	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	s.Start()
	defer s.Close()

	natsClient, err := newNatsClient("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("newNatsClient(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{
		srd, err := natsClient.Recording.Stored.Data("sr1")

		failed = err != nil
		failed = failed || srd.ID() != "sr1"
		if failed {
			t.Errorf("nc.Recording.Stored.Data('sr1') => '%v', '%v', expected '%v', '%v'",
				srd, err, storedRecordingData, nil)
		}
	}
	{
		srd, err := natsClient.Recording.Stored.Data("sr2")

		failed = err == nil || errors.Cause(err).Error() != "Failed to get stored recording"
		failed = failed || srd.ID() != ""
		if failed {
			t.Errorf("nc.Recording.Stored.Data('sr2') => '%v', '%v', expected '%v', '%v'",
				srd, err, ari.StoredRecordingData{}, "Failed to get stored recording")
		}
	}
}

func TestStoredRecordingActions(t *testing.T) {

	//TODO: embed nats?

	bin, err := exec.LookPath("gnatsd")
	if err != nil {
		t.Skip("No gnatsd binary in PATH, skipping")
	}

	cmd := exec.Command(bin, "-p", "4333")
	if err := cmd.Start(); err != nil {
		t.Errorf("Unable to run gnatsd: '%v'", err)
		return
	}

	defer func() {
		cmd.Process.Signal(syscall.SIGTERM)
		cmd.Wait()
	}()

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoredRecording := mock.NewMockStoredRecording(ctrl)

	gomock.InOrder(
		mockStoredRecording.EXPECT().Copy("sr1", "sr3").Return(ari.NewStoredRecordingHandle("sr3", mockStoredRecording), nil),
		mockStoredRecording.EXPECT().Delete("sr1").Return(nil),
	)

	cl := &ari.Client{
		Recording: &ari.Recording{
			Stored: mockStoredRecording,
		},
	}

	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	s.Start()
	defer s.Close()

	natsClient, err := newNatsClient("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("newNatsClient(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}
	{
		_, err := natsClient.Recording.Stored.Copy("sr1", "sr3")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Stored.Copy('sr1', 'sr3') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Stored.Delete("sr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Stored.Delete('sr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		type natsConn interface {
			NatsConnection() *nc.Conn
		}

		conn := natsClient.Recording.Stored.(natsConn).NatsConnection()
		msg, err := conn.RawRequest("ari.recording.stored.copy.sr1", []byte("asdf"))

		var decodingErrorMessage = "Remote Error in Endpoint 'ari.recording.stored.copy.sr1': Error decoding JSON body: invalid character 'a' looking for beginning of value"

		failed = msg != nil && err == nil || err.Error() != decodingErrorMessage
		if failed {
			t.Errorf("nc.Conn.RawRequest('ari.recording.stored.copy.sr1', 'asdf') => '%v', '%v', expected '%v', '%v'",
				msg, err,
				nil, decodingErrorMessage)
		}
	}

}
