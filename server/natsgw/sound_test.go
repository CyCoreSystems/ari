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

func TestSoundList(t *testing.T) {

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

	mockSound := mock.NewMockSound(ctrl)

	gomock.InOrder(
		mockSound.EXPECT().List(nil).Return([]*ari.SoundHandle{
			ari.NewSoundHandle("s1", mockSound),
			ari.NewSoundHandle("s2", mockSound),
		}, nil),

		mockSound.EXPECT().List(nil).Return([]*ari.SoundHandle{}, errors.New("Error getting sounds")),

		mockSound.EXPECT().List(gomock.Any()).Return(
			[]*ari.SoundHandle{ari.NewSoundHandle("s1", mockSound)}, nil),
	)

	cl := &ari.Client{
		Sound: mockSound,
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
		sx, err := natsClient.Sound.List(nil)

		failed = err != nil
		failed = failed || len(sx) != 2
		if failed {
			t.Errorf("nc.Sound.List(nil) => '%v', '%v', expected '%v', '%v'", sx, err, "s1,s2", nil)
		}
	}
	{
		sx, err := natsClient.Sound.List(nil)

		failed = err == nil || errors.Cause(err).Error() != "Error getting sounds"
		failed = failed || len(sx) != 0
		if failed {
			t.Errorf("nc.Sound.List(nil) => '%v', '%v', expected '%v', '%v'", sx, err, "", "Error getting sounds")
		}
	}
	{

		filters := map[string]string{
			"lang": "X",
		}

		sx, err := natsClient.Sound.List(filters)

		failed = err != nil
		failed = failed || len(sx) != 1
		if failed {
			t.Errorf("nc.Sound.List(lang = 'X') => '%v', '%v', expected '%v', '%v'", sx, err, "s1", nil)
		}
	}
}

func TestSoundData(t *testing.T) {

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

	mockSound := mock.NewMockSound(ctrl)

	var soundData ari.SoundData
	soundData.ID = "s1"
	soundData.Text = "hello"

	gomock.InOrder(
		mockSound.EXPECT().Data("s1").Return(soundData, nil),
		mockSound.EXPECT().Data("s2").Return(ari.SoundData{}, errors.New("Failed to get sound")),
	)

	cl := &ari.Client{
		Sound: mockSound,
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
		sd, err := natsClient.Sound.Data("s1")

		failed = err != nil
		failed = failed || sd.ID != "s1" || sd.Text != "hello"
		if failed {
			t.Errorf("nc.Sound.Data('s1') => '%v', '%v', expected '%v', '%v'",
				sd, err, soundData, nil)
		}
	}
	{
		sd, err := natsClient.Sound.Data("s2")

		failed = err == nil || errors.Cause(err).Error() != "Failed to get sound"
		failed = failed || sd.ID != ""
		if failed {
			t.Errorf("nc.Sound.Data('s2') => '%v', '%v', expected '%v', '%v'",
				sd, err, ari.SoundData{}, "Failed to get sound")
		}
	}
}
