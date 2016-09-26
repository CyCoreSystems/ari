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

func TestLoggingCreate(t *testing.T) {

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

	mockLogging := mock.NewMockLogging(ctrl)
	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Logging().Times(2).Return(mockLogging)

	mockLogging.EXPECT().Create("name", "levels").Return(nil)
	mockLogging.EXPECT().Create("name2", "levels").Return(errors.New("Failed to create log"))

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
		err := natsClient.Asterisk.Logging().Create("name", "levels")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Logging().Create('name', 'levels') => '%v', expected '%v'", err, nil)
		}
	}

	{
		err := natsClient.Asterisk.Logging().Create("name2", "levels")

		failed = err == nil || errors.Cause(err).Error() != "Failed to create log"
		if failed {
			t.Errorf("nc.Asterisk.Logging().Create('name2', 'levels') => '%v', expected '%v'", err, "Failed to create log")
		}
	}

}

func TestLoggingList(t *testing.T) {

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

	mockLogging := mock.NewMockLogging(ctrl)
	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Logging().Return(mockLogging)

	mockLogging.EXPECT().List().Return([]ari.LogData{}, nil)

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
		ld, err := natsClient.Asterisk.Logging().List()

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Logging().List() => '%v', '%v', expected '%v', '%v'",
				ld, err, "[]", nil)
		}
	}

}

func TestLoggingRotate(t *testing.T) {

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

	mockLogging := mock.NewMockLogging(ctrl)
	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Logging().Return(mockLogging)

	mockLogging.EXPECT().Rotate("name1").Return(nil)

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
		err := natsClient.Asterisk.Logging().Rotate("name1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Logging().Rotate() => '%v', expected '%v'",
				err, nil)
		}
	}

}

func TestLoggingDelete(t *testing.T) {

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

	mockLogging := mock.NewMockLogging(ctrl)
	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Logging().Return(mockLogging)

	mockLogging.EXPECT().Delete("name1").Return(nil)

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
		err := natsClient.Asterisk.Logging().Delete("name1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Logging().Delete() => '%v', expected '%v'",
				err, nil)
		}
	}

}
