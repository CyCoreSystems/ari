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
