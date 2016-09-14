package natsgw

import (
	"errors"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/client/nc"
	"github.com/golang/mock/gomock"
)

var ServerWaitDelay = 500 * time.Millisecond

func TestApplicationsSubscribeUnsubscribe(t *testing.T) {

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

	mockApplication := mock.NewMockApplication(ctrl)
	mockApplication.EXPECT().Subscribe("app1", "evt1").Return(nil)
	mockApplication.EXPECT().Subscribe("app2", "evt1").Return(errors.New("Application not found"))
	mockApplication.EXPECT().Unsubscribe("app1", "evt1").Return(nil)
	mockApplication.EXPECT().Unsubscribe("app2", "evt1").Return(errors.New("Application not found"))

	cl := &ari.Client{
		Application: mockApplication,
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()
	defer s.Close()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{
		err = natsClient.Application.Subscribe("app1", "evt1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Application.Subscribe(app1,evt1) => '%v', expected '%v'", err, "nil")
		}
	}

	{
		err = natsClient.Application.Subscribe("app2", "evt1")

		failed = err == nil || err.Error() != "Application not found"
		if failed {
			t.Errorf("nc.Application.Subscribe(app2,evt1) => '%v', expected '%v'", err, "nil")
		}
	}

	{
		err = natsClient.Application.Unsubscribe("app1", "evt1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Application.Unsubscribe(app1,evt1) => '%v', expected '%v'", err, "nil")
		}
	}

	{
		err = natsClient.Application.Unsubscribe("app2", "evt1")

		failed = err == nil || err.Error() != "Application not found"
		if failed {
			t.Errorf("nc.Application.Unsubscribe(app2,evt1) => '%v', expected '%v'", err, "nil")
		}
	}

}
func TestApplicationsData(t *testing.T) {

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

	// test clientiontruc

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var appData ari.ApplicationData

	mockApplication := mock.NewMockApplication(ctrl)
	mockApplication.EXPECT().Data("app1").Return(appData, nil)
	mockApplication.EXPECT().Data("app2").Return(appData, errors.New("Application not found"))

	cl := &ari.Client{
		Application: mockApplication,
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()
	defer s.Close()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{
		appData, err := natsClient.Application.Data("app1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Application.Data(app1) => '%v', '%v', expected '%v', '%v'", appData, err, "appData", nil)
		}
	}

	{
		appData, err := natsClient.Application.Data("app2")

		failed = err == nil || err.Error() != "Application not found"
		if failed {
			t.Errorf("nc.Application.Data(app2) => '%v', '%v', expected '%v', '%v'", appData, err, "appData", "Application not found")
		}
	}

}

func TestApplicationsList(t *testing.T) {

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

	// test clientiontruc

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplication := mock.NewMockApplication(ctrl)
	mockApplication.EXPECT().List().Return([]*ari.ApplicationHandle{ari.NewApplicationHandle("app1", mockApplication), ari.NewApplicationHandle("app2", mockApplication)}, nil)

	cl := &ari.Client{
		Application: mockApplication,
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()
	defer s.Close()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	apps, err := natsClient.Application.List()

	failed = len(apps) != 2 || err != nil
	if failed {
		t.Errorf("nc.Application.List() => {%v, %v}, expected {%v, %v}", apps, err, "[app1,app2]", "nil")
	}

}

func TestApplicationsListError(t *testing.T) {

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

	mockApplication := mock.NewMockApplication(ctrl)
	mockApplication.EXPECT().List().Return(nil, errors.New("Error getting application list"))

	cl := &ari.Client{
		Application: mockApplication,
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()
	defer s.Close()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	apps, err := natsClient.Application.List()

	failed = len(apps) != 0 || err == nil
	if failed {
		t.Errorf("nc.Application.List() => {%v, %v}, expected {%v, %v}", apps, err, "[]", "err")
	}
}
