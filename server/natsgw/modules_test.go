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

func TestModulesList(t *testing.T) {

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModules := mock.NewMockModules(ctrl)
	mockModules.EXPECT().List().Return([]*ari.ModuleHandle{ari.NewModuleHandle("mod1", mockModules), ari.NewModuleHandle("mod2", mockModules)}, nil)

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Modules().MinTimes(1).Return(mockModules)

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

	mods, err := natsClient.Asterisk.Modules().List()

	failed = len(mods) != 2 || err != nil
	if failed {
		t.Errorf("nc.Asterisk.Modules().List() => {%v, %v}, expected {%v, %v}", mods, err, "[mod1,mod2]", "nil")
	}

}

func TestModulesData(t *testing.T) {

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	modData := ari.ModuleData{Name: "mod1"}

	mockModules := mock.NewMockModules(ctrl)
	mockModules.EXPECT().Data("mod1").Return(modData, nil)
	mockModules.EXPECT().Data("mod2").Return(ari.ModuleData{}, errors.New("Failed to get module"))

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Modules().MinTimes(1).Return(mockModules)

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
		md, err := natsClient.Asterisk.Modules().Data("mod1")

		failed = err != nil
		failed = failed || md.Name != "mod1"
		if failed {
			t.Errorf("nc.Asterisk.Modules().Data('mod1') => {%v, %v}, expected {%v, %v}",
				md, err, "{name='mod1'}", "nil")
		}
	}

	{
		md, err := natsClient.Asterisk.Modules().Data("mod2")

		failed = err == nil || errors.Cause(err).Error() != "Failed to get module"
		failed = failed || md.Name != ""
		if failed {
			t.Errorf("nc.Asterisk.Modules().Data('mod2') => {%v, %v}, expected {%v, %v}",
				md, err, "{}", "Failed to get module")
		}
	}
}

func TestModulesActions(t *testing.T) {

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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModules := mock.NewMockModules(ctrl)
	mockModules.EXPECT().Reload("mod1").Return(nil)
	mockModules.EXPECT().Load("mod1").Return(nil)
	mockModules.EXPECT().Unload("mod1").Return(nil)

	mockAsterisk := mock.NewMockAsterisk(ctrl)
	mockAsterisk.EXPECT().Modules().AnyTimes().Return(mockModules)

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
		err := natsClient.Asterisk.Modules().Unload("mod1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Modules().Unload('mod1') => {%v}, expected {%v}",
				err, "nil")
		}
	}

	{
		err := natsClient.Asterisk.Modules().Reload("mod1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Modules().Reload('mod1') => {%v}, expected {%v}",
				err, "nil")
		}
	}

	{
		err := natsClient.Asterisk.Modules().Load("mod1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Asterisk.Modules().Load('mod1') => {%v}, expected {%v}",
				err, "nil")
		}
	}

}
