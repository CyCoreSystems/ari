package natsgw

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/mock"
	"github.com/CyCoreSystems/ari/client/nc"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestPlaybackData(t *testing.T) {

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

	var playbackData ari.PlaybackData
	var playbackErrorMessage = "Could not find playback"

	mockPlayback := mock.NewMockPlayback(ctrl)
	mockPlayback.EXPECT().Data("pb1").Return(playbackData, nil)
	mockPlayback.EXPECT().Data("pb2").Return(playbackData, errors.New(playbackErrorMessage))

	cl := &ari.Client{
		Playback: mockPlayback,
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
		ret, err := natsClient.Playback.Data("pb1")
		failed = err != nil
		if failed {
			t.Errorf("nc.Playback.Data('pb1') => ('%v','%v'), expected ('%v','%v')",
				ret, err,
				playbackData, nil)
		}
	}

	{
		ret, err := natsClient.Playback.Data("pb2")
		failed = err == nil || errors.Cause(err).Error() != playbackErrorMessage
		if failed {
			t.Errorf("nc.Playback.Data('pb2') => ('%v','%v'), expected ('%v','%v')",
				ret, err,
				playbackData, playbackErrorMessage)
		}
	}

}

func TestPlaybackControl(t *testing.T) {

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

	var playbackErrorMessage = "Could not find playback"

	mockPlayback := mock.NewMockPlayback(ctrl)
	mockPlayback.EXPECT().Control("pb1", "command").Return(nil)
	mockPlayback.EXPECT().Control("pb2", "command").Return(errors.New(playbackErrorMessage))

	cl := &ari.Client{
		Playback: mockPlayback,
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
		err := natsClient.Playback.Control("pb1", "command")
		failed = err != nil
		if failed {
			t.Errorf("nc.Playback.Control('pb1', 'command') => '%v', expected '%v'",
				err,
				nil)
		}
	}

	{
		err := natsClient.Playback.Control("pb2", "command")
		failed = err == nil || errors.Cause(err).Error() != playbackErrorMessage
		if failed {
			t.Errorf("nc.Playback.Control('pb2', 'command') => '%v', expected '%v'",
				err,
				playbackErrorMessage)
		}
	}

	{
		type natsConn interface {
			NatsConnection() *nc.Conn
		}

		conn := natsClient.Asterisk.(natsConn).NatsConnection()
		msg, err := conn.RawRequest("ari.playback.control.pb1", []byte("asdf"))

		var decodingErrorMessage = "Remote Error in Endpoint 'ari.playback.control.pb1': Error decoding JSON body: invalid character 'a' looking for beginning of value"

		failed = msg != nil && err == nil || err.Error() != decodingErrorMessage
		if failed {
			t.Errorf("nc.Conn.RawRequest('ari.playback.control.pb1', 'asdf') => '%v', '%v', expected '%v', '%v'",
				msg, err,
				nil, decodingErrorMessage)
		}
	}

}

func TestPlaybackStop(t *testing.T) {

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

	var playbackErrorMessage = "Could not find playback"

	mockPlayback := mock.NewMockPlayback(ctrl)
	mockPlayback.EXPECT().Stop("pb1").Return(nil)
	mockPlayback.EXPECT().Stop("pb2").Return(errors.New(playbackErrorMessage))

	cl := &ari.Client{
		Playback: mockPlayback,
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
		err := natsClient.Playback.Stop("pb1")
		failed = err != nil
		if failed {
			t.Errorf("nc.Playback.Stop('pb1') => '%v', expected '%v'",
				err,
				nil)
		}
	}

	{

		err := natsClient.Playback.Stop("pb2")
		failed = err == nil || errors.Cause(err).Error() != playbackErrorMessage
		if failed {
			t.Errorf("nc.Playback.Stop('pb2') => '%v', expected '%v'",
				err,
				playbackErrorMessage)
		}
	}

}
