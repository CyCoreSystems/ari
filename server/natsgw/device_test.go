package natsgw

import (
	"errors"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/golang/mock/gomock"
)

func TestDeviceStateUpdate(t *testing.T) {

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

	mockDeviceState := mock.NewMockDeviceState(ctrl)
	mockDeviceState.EXPECT().Update("ds1", "state1").Return(nil)

	cl := &ari.Client{
		DeviceState: mockDeviceState,
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
		err := natsClient.DeviceState.Update("ds1", "state1")

		failed = err != nil
		if failed {
			t.Errorf("nc.DeviceState.Update(ds1, state1) => '%v', expected '%v'", err, nil)
		}
	}

}

func TestDeviceStateData(t *testing.T) {

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

	var deviceData ari.DeviceStateData

	mockDeviceState := mock.NewMockDeviceState(ctrl)
	mockDeviceState.EXPECT().Data("ds1").Return(deviceData, nil)
	mockDeviceState.EXPECT().Data("ds2").Return(deviceData, errors.New("DeviceState not found"))

	cl := &ari.Client{
		DeviceState: mockDeviceState,
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
		dsData, err := natsClient.DeviceState.Data("ds1")

		failed = err != nil
		if failed {
			t.Errorf("nc.DeviceState.Data(ds1) => '%v', '%v', expected '%v', '%v'", dsData, err, "dsData", nil)
		}
	}

	{
		dsData, err := natsClient.DeviceState.Data("ds2")

		failed = err == nil || err.Error() != "DeviceState not found"
		if failed {
			t.Errorf("nc.DeviceState.Data(ds2) => '%v', '%v', expected '%v', '%v'", dsData, err, "dsData", "DeviceState not found")
		}
	}

}

func TestDeviceStateList(t *testing.T) {

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

	mockDeviceState := mock.NewMockDeviceState(ctrl)
	mockDeviceState.EXPECT().List().Return([]*ari.DeviceStateHandle{
		ari.NewDeviceStateHandle("ds1", mockDeviceState),
		ari.NewDeviceStateHandle("ds2", mockDeviceState),
	}, nil)

	cl := &ari.Client{
		DeviceState: mockDeviceState,
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

	devices, err := natsClient.DeviceState.List()

	failed = len(devices) != 2 || err != nil
	if failed {
		t.Errorf("nc.DeviceState.List() => {%v, %v}, expected {%v, %v}", devices, err, "[ds1,ds2]", "nil")
	}

}
