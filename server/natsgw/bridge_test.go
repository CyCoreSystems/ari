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

func TestBridgeList(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().List().Return([]*ari.BridgeHandle{
		ari.NewBridgeHandle("b1", mockBridge),
		ari.NewBridgeHandle("b2", mockBridge),
	}, nil)

	mockBridge.EXPECT().List().Return([]*ari.BridgeHandle{}, errors.New("Error getting bridges"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		bx, err := natsClient.Bridge.List()

		failed = err != nil
		failed = failed || len(bx) != 2
		if failed {
			t.Errorf("nc.Bridge.List() => '%v', '%v', expected '%v', '%v'", bx, err, "b1,b2", nil)
		}
	}
	{
		bx, err := natsClient.Bridge.List()

		failed = err == nil || errors.Cause(err).Error() != "Error getting bridges"
		failed = failed || len(bx) != 0
		if failed {
			t.Errorf("nc.Bridge.List() => '%v', '%v', expected '%v', '%v'", bx, err, "", "Error getting bridges")
		}
	}

}

func TestBridgeData(t *testing.T) {

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

	var bridgeData ari.BridgeData

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().Data("b1").Return(bridgeData, nil)
	mockBridge.EXPECT().Data("b2").Return(ari.BridgeData{}, errors.New("Bridge not found"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		ret, err := natsClient.Bridge.Data("b1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Bridge.Data('%v') => '%v', '%v', expected '%v', '%v'", "b1",
				ret, err,
				bridgeData, nil)
		}
	}
	{

		ret, err := natsClient.Bridge.Data("b2")

		failed = err == nil || errors.Cause(err).Error() != "Bridge not found"
		if failed {
			t.Errorf("nc.Bridge.Data('%v') => '%v', '%v', expected '%v', '%v'", "b2",
				ret, err,
				bridgeData, "Bridge not found")
		}
	}
}

func TestBridgeAddChannel(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().AddChannel("b1", "c1").Return(nil)
	mockBridge.EXPECT().AddChannel("b2", "c2").Return(errors.New("Bridge not found"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		err := natsClient.Bridge.AddChannel("b1", "c1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Bridge.AddChannel('%v', '%v') => '%v', expected '%v'",
				"b1", "c1",
				err, nil)
		}
	}
	{

		err := natsClient.Bridge.AddChannel("b2", "c2")

		failed = err == nil || errors.Cause(err).Error() != "Bridge not found"
		if failed {
			t.Errorf("nc.Bridge.AddChannel('%v', '%v') => '%v', expected '%v'",
				"b2", "c2",
				err, "Bridge not found")
		}
	}
}

func TestBridgeRemoveChannel(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().RemoveChannel("b1", "c1").Return(nil)
	mockBridge.EXPECT().RemoveChannel("b2", "c2").Return(errors.New("Bridge not found"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		err := natsClient.Bridge.RemoveChannel("b1", "c1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Bridge.RemoveChannel('%v', '%v') => '%v', expected '%v'",
				"b1", "c1",
				err, nil)
		}
	}
	{

		err := natsClient.Bridge.RemoveChannel("b2", "c2")

		failed = err == nil || errors.Cause(err).Error() != "Bridge not found"
		if failed {
			t.Errorf("nc.Bridge.RemoveChannel('%v', '%v') => '%v', expected '%v'",
				"b2", "c2",
				err, "Bridge not found")
		}
	}
}

func TestBridgeDelete(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().Delete("b1").Return(nil)
	mockBridge.EXPECT().Delete("b2").Return(errors.New("Bridge not found"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		err := natsClient.Bridge.Delete("b1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Bridge.Delete('%v') => '%v', expected '%v'", "b1", err, nil)
		}
	}
	{

		err := natsClient.Bridge.Delete("b2")

		failed = err == nil || errors.Cause(err).Error() != "Bridge not found"
		if failed {
			t.Errorf("nc.Bridge.Delete('%v') => '%v', expected '%v'", "b2", err, "Bridge not found")
		}
	}
}

func TestBridgeCreate(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().Create("b1", "", "").Return(ari.NewBridgeHandle("b1", mockBridge), nil)
	mockBridge.EXPECT().Create("b2", "", "").Return(nil, errors.New("Error creating bridge"))

	cl := &ari.Client{
		Bridge: mockBridge,
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
		bh, err := natsClient.Bridge.Create("b1", "", "")

		failed = err != nil || bh == nil || bh.ID() != "b1"
		if failed {
			t.Errorf("nc.Bridge.Create('b1','','') => '%v', '%v', expected '%v', '%v'",
				bh, err,
				"bridgeHandle{b1}", nil)
		}
	}

	{
		bh, err := natsClient.Bridge.Create("b2", "", "")

		failed = bh != nil || err == nil || errors.Cause(err).Error() != "Error creating bridge"
		if failed {
			t.Errorf("nc.Bridge.Create('b2','','') => '%v', '%v', expected '%v', '%v'",
				bh, err,
				nil, "Error creating bridge")
		}
	}
}

func TestBridgeRecord(t *testing.T) {

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

	mockBridge := mock.NewMockBridge(ctrl)
	mockBridge.EXPECT().Record("b1", "name1", gomock.Any()).Return(
		ari.NewLiveRecordingHandle("name1", mockLiveRecording), nil)

	cl := &ari.Client{
		Bridge: mockBridge,
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
		lrh, err := natsClient.Bridge.Record("b1", "name1", nil)

		failed = err != nil || lrh == nil || lrh.ID() != "name1"
		if failed {
			t.Errorf("nc.Bridge.Record('b1','name',nil) => '%v', '%v', expected '%v', '%v'",
				lrh, err,
				"liveRecordingHandle{name1}", nil)
		}
	}
}
