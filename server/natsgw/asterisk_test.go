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

func TestAsteriskModuleReload(t *testing.T) {

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

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().ReloadModule("module1").Return(nil)
	mockAsterisk.EXPECT().ReloadModule("module2").Return(errors.New("Can't find module"))

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

	go s.Listen()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	err = natsClient.Asterisk.ReloadModule("module1")

	failed = err != nil
	if failed {
		t.Errorf("nc.Asterisk.ReloadModule(%s) => %v, expected %v", "module1", err, nil)
	}

	err = natsClient.Asterisk.ReloadModule("module2")

	failed = err == nil || err.Error() != "Can't find module"
	if failed {
		t.Errorf("nc.Asterisk.ReloadModule(%s) => %v, expected %v", "module2", err, "Can't find module")
	}

}

func TestAsteriskInfo(t *testing.T) {

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

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	info := ari.AsteriskInfo{
		BuildInfo: ari.BuildInfo{Date: "Date1", Kernel: "kernel1", Machine: "machine1", Options: "o"},
	}

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Info("").Return(&info, nil)
	mockAsterisk.EXPECT().Info("").Return(nil, errors.New("Err getting info"))

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

	go s.Listen()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{

		ret, err := natsClient.Asterisk.Info("")

		failed = err != nil
		failed = failed || ret == nil || info.BuildInfo.Kernel != "kernel1"
		if failed {
			t.Errorf("nc.Asterisk.Info(%s) => %v, %v, expected %v, %v", "", ret, err, info, nil)
		}

		ret, err = natsClient.Asterisk.Info("")

		failed = err == nil || err.Error() != "Err getting info"
		if failed {
			t.Errorf("nc.Asterisk.Info(%s) => %v, %v, expected %v, %v", "", ret, err, nil, "Err getting info")
		}

	}
}

func TestAsteriskVariableGet(t *testing.T) {

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

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVariables := mock.NewMockVariables(ctrl)
	mockVariables.EXPECT().Get("var1").Return("val1", nil)
	mockVariables.EXPECT().Get("var2").Return("", errors.New("Variable not found"))

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Variables().Return(mockVariables)
	mockAsterisk.EXPECT().Variables().Return(mockVariables)

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

	go s.Listen()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{

		ret, err := natsClient.Asterisk.Variables().Get("var1")

		failed = err != nil
		failed = failed || ret != "val1"
		if failed {
			t.Errorf("nc.Asterisk.Variables().Get(%s) => %v, %v, expected %v, %v", "var1", ret, err, "val1", nil)
		}

		ret, err = natsClient.Asterisk.Variables().Get("var2")

		failed = err == nil
		failed = failed || ret != ""
		if failed {
			t.Errorf("nc.Asterisk.Variables().Get(%s) => %v, %v, expected %v, %v", "", ret, err, "", "Variable not found")
		}

	}
}

func TestAsteriskVariableSet(t *testing.T) {

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

	<-time.After(ServerWaitDelay)

	// test client

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVariables := mock.NewMockVariables(ctrl)
	mockVariables.EXPECT().Set("var1", "val1").Return(nil)
	mockVariables.EXPECT().Set("var2", "val2").Return(errors.New("Malformed variable name"))

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Variables().Return(mockVariables)
	mockAsterisk.EXPECT().Variables().Return(mockVariables)

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

	go s.Listen()

	natsClient, err := nc.New("nats://127.0.0.1:4333")

	failed = natsClient == nil || err != nil
	if failed {
		t.Errorf("nc.New(url) => {%v, %v}, expected {%v, %v}", natsClient, err, "cl", "nil")
	}

	{

		err := natsClient.Asterisk.Variables().Set("var1", "val1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Variables().Set(%s, %s) => %v, expected %v", "var1", "val1", err, nil)
		}

		err = natsClient.Asterisk.Variables().Set("var2", "val2")

		failed = err == nil || err.Error() != "Malformed variable name"
		if failed {
			t.Errorf("nc.Asterisk.Variables().Set(%s, %s) => %v, expected %v", "var2", "val2", err, "Malformed variable name")
		}

	}
}
