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

	go s.Listen()
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

		failed = err == nil || err.Error() != "Error getting channels"
		failed = failed || len(cx) != 0
		if failed {
			t.Errorf("nc.Channel.List() => '%v', '%v', expected '%v', '%v'", cx, err, "", "Error getting channels")
		}
	}
}
