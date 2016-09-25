package natsgw

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestLiveRecordingData(t *testing.T) {

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

	mockLiveRecording := mock.NewMockLiveRecording(ctrl)

	var liveRecordingData ari.LiveRecordingData
	liveRecordingData.Name = "lr1"

	gomock.InOrder(
		mockLiveRecording.EXPECT().Data("lr1").Return(liveRecordingData, nil),
		mockLiveRecording.EXPECT().Data("lr2").Return(ari.LiveRecordingData{}, errors.New("Failed to get live recording")),
	)

	cl := &ari.Client{
		Recording: &ari.Recording{
			Live: mockLiveRecording,
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
		lrd, err := natsClient.Recording.Live.Data("lr1")

		failed = err != nil
		failed = failed || lrd.ID() != "lr1"
		if failed {
			t.Errorf("nc.Recording.Live.Data('lr1') => '%v', '%v', expected '%v', '%v'",
				lrd, err, liveRecordingData, nil)
		}
	}
	{
		lrd, err := natsClient.Recording.Live.Data("lr2")

		failed = err == nil || errors.Cause(err).Error() != "Failed to get live recording"
		failed = failed || lrd.ID() != ""
		if failed {
			t.Errorf("nc.Recording.Live.Data('lr2') => '%v', '%v', expected '%v', '%v'",
				lrd, err, ari.LiveRecordingData{}, "Failed to get live recording")
		}
	}
}

func TestLiveRecordingActions(t *testing.T) {

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

	mockLiveRecording := mock.NewMockLiveRecording(ctrl)

	gomock.InOrder(
		mockLiveRecording.EXPECT().Stop("lr1").Return(nil),
		mockLiveRecording.EXPECT().Mute("lr1").Return(nil),
		mockLiveRecording.EXPECT().Unmute("lr1").Return(nil),
		mockLiveRecording.EXPECT().Pause("lr1").Return(nil),
		mockLiveRecording.EXPECT().Resume("lr1").Return(nil),
		mockLiveRecording.EXPECT().Scrap("lr1").Return(nil),
		mockLiveRecording.EXPECT().Delete("lr1").Return(nil),
	)

	cl := &ari.Client{
		Recording: &ari.Recording{
			Live: mockLiveRecording,
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
		err := natsClient.Recording.Live.Stop("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Stop('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Live.Mute("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Mute('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Live.Unmute("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Unmute('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}

	{
		err := natsClient.Recording.Live.Pause("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Pause('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Live.Resume("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Resume('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Live.Scrap("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Scrap('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}
	{
		err := natsClient.Recording.Live.Delete("lr1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Recording.Live.Delete('lr1') => '%v', expected '%v'",
				err, nil)
		}
	}

}
