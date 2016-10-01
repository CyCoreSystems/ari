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

func TestConfigData(t *testing.T) {

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

	mockConfig := mock.NewMockConfig(ctrl)

	var configData ari.ConfigData
	configData.ID = "i1"
	configData.Class = "c1"
	configData.Type = "t1"
	configData.Fields = []ari.ConfigTuple{
		ari.ConfigTuple{"1", "2"},
		ari.ConfigTuple{"3", "4"},
	}

	gomock.InOrder(
		mockConfig.EXPECT().Data("c1", "t1", "i1").Return(configData, nil),
		mockConfig.EXPECT().Data("c2", "t1", "i1").Return(ari.ConfigData{}, errors.New("Failed to get config")),
	)

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Config().MinTimes(1).Return(mockConfig)

	cl := &ari.Client{
		Asterisk: mockAsterisk,
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
		cd, err := natsClient.Asterisk.Config().Data("c1", "t1", "i1")

		failed = err != nil
		failed = failed || cd.ID != "i1" || cd.Class != "c1" || cd.Type != "t1"
		failed = failed || len(cd.Fields) != 2
		if failed {
			t.Errorf("nc.Asterisk.Config().Data('c1', 't1', 'i1') => ('%v', '%v'), expected ('%v', '%v')",
				cd, err, configData, nil)
		}
	}
	{
		cd, err := natsClient.Asterisk.Config().Data("c2", "t1", "i1")

		failed = err == nil || errors.Cause(err).Error() != "Failed to get config"
		if failed {
			t.Errorf("nc.Asterisk.Config().Data('c2', 't1', 'i1') => ('%v', '%v'), expected ('%v', '%v')",
				cd, err, ari.ConfigData{}, "Failed to get config")
		}
	}
}

func TestConfigDelete(t *testing.T) {

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

	mockConfig := mock.NewMockConfig(ctrl)

	gomock.InOrder(
		mockConfig.EXPECT().Delete("c1", "t1", "i1").Return(nil),
		mockConfig.EXPECT().Delete("c2", "t1", "i1").Return(errors.New("Failed to delete config")),
	)

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Config().MinTimes(1).Return(mockConfig)

	cl := &ari.Client{
		Asterisk: mockAsterisk,
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
		err := natsClient.Asterisk.Config().Delete("c1", "t1", "i1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Config().Delete('c1', 't1', 'i1') => ('%v'), expected ('%v')",
				err, nil)
		}
	}
	{
		err := natsClient.Asterisk.Config().Delete("c2", "t1", "i1")

		failed = err == nil || errors.Cause(err).Error() != "Failed to delete config"
		if failed {
			t.Errorf("nc.Asterisk.Config().Delete ('c2', 't1', 'i1') => ('%v'), expected ('%v')",
				err, "Failed to delete config")
		}
	}
}

func TestConfigUpdate(t *testing.T) {

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

	mockConfig := mock.NewMockConfig(ctrl)

	fields := []ari.ConfigTuple{
		ari.ConfigTuple{"1", "2"},
		ari.ConfigTuple{"3", "4"},
	}

	gomock.InOrder(
		mockConfig.EXPECT().Update("c1", "t1", "i1", fields).Return(nil),
		mockConfig.EXPECT().Update("c2", "t1", "i1", fields).Return(errors.New("Failed to update config")),
	)

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Config().MinTimes(1).Return(mockConfig)

	cl := &ari.Client{
		Asterisk: mockAsterisk,
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
		err := natsClient.Asterisk.Config().Update("c1", "t1", "i1", fields)

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Config().Update('c1', 't1', 'i1', fl) => ('%v'), expected ('%v')",
				err, nil)
		}
	}
	{
		err := natsClient.Asterisk.Config().Update("c2", "t1", "i1", fields)

		failed = err == nil || errors.Cause(err).Error() != "Failed to update config"
		if failed {
			t.Errorf("nc.Asterisk.Config().Update('c2', 't1', 'i1', fl) => ('%v'), expected ('%v')",
				err, "Failed to update config")
		}
	}
}
