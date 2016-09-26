package natsgw

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/golang/mock/gomock"
)

func TestChannelListTest(t *testing.T) {

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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().List().Return([]*ari.ChannelHandle{
		ari.NewChannelHandle("c1", mockChannel),
		ari.NewChannelHandle("c2", mockChannel),
	}, nil)

	mockChannel.EXPECT().List().Return([]*ari.ChannelHandle{}, errors.New("Error getting channels"))

	mockSubscription := mock.NewMockSubscription(ctrl)
	mockSubscription.EXPECT().Cancel().AnyTimes() //cancel is in a defer, so it may not always be called
	mockSubscription.EXPECT().Events().Return(make(chan ari.Event))

	mockBus := mock.NewMockBus(ctrl)
	mockBus.EXPECT().Subscribe(ari.Events.All).Return(mockSubscription)

	cl := &ari.Client{
		Channel: mockChannel,
		Bus:     mockBus,
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
		cx, err := natsClient.Channel.List()

		failed = err != nil
		failed = failed || len(cx) != 2
		if failed {
			t.Errorf("nc.Channel.List() => '%v', '%v', expected '%v', '%v'", cx, err, "c1,c2", nil)
		}
	}
	{
		cx, err := natsClient.Channel.List()

		failed = err == nil || errors.Cause(err).Error() != "Error getting channels"
		failed = failed || len(cx) != 0
		if failed {
			t.Errorf("nc.Channel.List() => '%v', '%v', expected '%v', '%v'", cx, err, "", "Error getting channels")
		}
	}
}

func TestChannelAnswer(t *testing.T) {

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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().Answer("c1").Return(nil)

	cl := &ari.Client{
		Channel: mockChannel,
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
		err = natsClient.Channel.Answer("c1")

		failed = err != nil
		if failed {
			t.Errorf("nc.Channel.Answer() => '%v', expected '%v'", err, nil)
		}
	}
}

func TestChannelSendDTMF(t *testing.T) {

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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().SendDTMF("c1", "1234", gomock.Not(gomock.Nil())).Return(nil)

	cl := &ari.Client{
		Channel: mockChannel,
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
		err = natsClient.Channel.SendDTMF("c1", "1234", nil)

		failed = err != nil
		if failed {
			t.Errorf("nc.Channel.SendDTMF('c1', '1234', nil) => '%v', expected '%v'", err, nil)
		}
	}
}

func TestChannelContinue(t *testing.T) {

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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().Continue("c1", "1", "2", 3).Return(nil)

	cl := &ari.Client{
		Channel: mockChannel,
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
		err = natsClient.Channel.Continue("c1", "1", "2", 3)

		failed = err != nil
		if failed {
			t.Errorf("nc.Channel.Continue('c1', '1', '2', '3') => '%v', expected '%v'", err, nil)
		}
	}
}

func TestChannelDial(t *testing.T) {
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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().Dial("c1", "c2", 3*time.Second).Return(nil)
	mockChannel.EXPECT().Dial("c2", "c3", 3*time.Second).Return(errors.New("Failed to dial"))

	cl := &ari.Client{
		Channel: mockChannel,
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
		err = natsClient.Channel.Dial("c1", "c2", 3*time.Second)

		failed = err != nil
		if failed {
			t.Errorf("nc.Channel.Dial('c1', 'c2', 3s) => '%v', expected '%v'", err, nil)
		}
	}

	{
		err = natsClient.Channel.Dial("c2", "c3", 3*time.Second)

		failed = err == nil || errors.Cause(err).Error() != "Failed to dial"
		if failed {
			t.Errorf("nc.Channel.Dial('c1', 'c2', 3s) => '%v', expected '%v'", err, "Failed to dial")
		}
	}

}

func TestChannelSnoop(t *testing.T) {
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

	mockChannel := mock.NewMockChannel(ctrl)
	c2 := ari.NewChannelHandle("c1", mockChannel)

	mockChannel.EXPECT().Snoop("c1", "c2", "app1", gomock.Not(gomock.Nil())).Return(c2, nil)
	mockChannel.EXPECT().Snoop("c3", "c4", "app1", gomock.Not(gomock.Nil())).Return(nil, errors.New("Error snooping"))

	cl := &ari.Client{
		Channel: mockChannel,
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
		handle, err := natsClient.Channel.Snoop("c1", "c2", "app1", &ari.SnoopOptions{})

		failed = err != nil
		failed = failed || handle == nil || handle.ID() != "c2"
		if failed {
			t.Errorf("nc.Channel.Snoop('c1', 'c2', 'app1', {}) => '%v', '%v', expected '%v', '%v'",
				handle, err,
				"c2", nil)
		}
	}

	{
		handle, err := natsClient.Channel.Snoop("c3", "c4", "app1", &ari.SnoopOptions{})

		failed = err == nil || errors.Cause(err).Error() != "Error snooping"
		failed = failed || handle != nil
		if failed {
			t.Errorf("nc.Channel.Snoop('c3', 'c4', 'app1', {}) => '%v', '%v', expected '%v', '%v'",
				handle, err,
				"nil", "Error snooping")
		}
	}

}

func TestChannelRecord(t *testing.T) {

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

	mockChannel := mock.NewMockChannel(ctrl)
	mockChannel.EXPECT().Record("c1", "name1", gomock.Any()).Return(
		ari.NewLiveRecordingHandle("name1", mockLiveRecording), nil)

	cl := &ari.Client{
		Channel: mockChannel,
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
		lrh, err := natsClient.Channel.Record("c1", "name1", nil)

		failed = err != nil || lrh == nil || lrh.ID() != "name1"
		if failed {
			t.Errorf("nc.Channel.Record('c1','name',nil) => '%v', '%v', expected '%v', '%v'",
				lrh, err,
				"liveRecordingHandle{name1}", nil)
		}
	}
}
