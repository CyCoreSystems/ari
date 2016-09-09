package natsgw

import (
	"errors"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/nc"
)

func TestApplicationsSubscribe(t *testing.T) {

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
	}()

	<-time.After(4 * time.Second)

	// test clientiontruc

	cl := &ari.Client{
		Application: testApplication(0),
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	err = natsClient.Application.Subscribe("app1", "evt1")

	failed = err != nil
	if failed {
		t.Errorf("nc.Application.Subscribe(app1,evt1) => '%v', expected '%v'", err, "nil")
	}

	s.Close()
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
	}()

	<-time.After(4 * time.Second)

	// test clientiontruc

	cl := &ari.Client{
		Application: testApplication(0),
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()

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

	s.Close()
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
	}()

	<-time.After(4 * time.Second)

	// test clientiontruc

	cl := &ari.Client{
		Application: testApplicationListError(0),
	}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()

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

	s.Close()
}

type testApplication int

func (a testApplication) List() (ax []*ari.ApplicationHandle, err error) {
	ax = append(ax, ari.NewApplicationHandle("app1", a))
	ax = append(ax, ari.NewApplicationHandle("app2", a))
	return
}

func (a testApplication) Get(name string) *ari.ApplicationHandle {
	panic("not implemented")
}

func (a testApplication) Data(name string) (ari.ApplicationData, error) {
	panic("not implemented")
}

func (a testApplication) Subscribe(name string, eventSource string) error {
	if name != "app1" || eventSource != "evt1" {
		return errors.New("Not Found")
	}
	return nil
}

func (a testApplication) Unsubscribe(name string, eventSource string) error {
	if name != "app1" || eventSource != "evt1" {
		return errors.New("Not Found")
	}
	return nil
}

type testApplicationListError int

func (a testApplicationListError) List() (ax []*ari.ApplicationHandle, err error) {
	err = errors.New("Dummy Error")
	return
}

func (a testApplicationListError) Get(name string) *ari.ApplicationHandle {
	panic("not implemented")
}

func (a testApplicationListError) Data(name string) (ari.ApplicationData, error) {
	panic("not implemented")
}

func (a testApplicationListError) Subscribe(name string, eventSource string) error {
	panic("not implemented")
}

func (a testApplicationListError) Unsubscribe(name string, eventSource string) error {
	panic("not implemented")
}
