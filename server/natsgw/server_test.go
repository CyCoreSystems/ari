package natsgw

import (
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
)

func TestNewServer(t *testing.T) {
	{
		s, err := NewServer(nil, nil)

		failed := s != nil || err == nil
		if failed {
			t.Errorf("natsgw.NewServer(nil, nil) => {%v, %v}, expected {%v, %v}", s, err, "nil", "err")
		}

		s.Close()
	}

	{
		cl := &ari.Client{}
		s, err := NewServer(cl, nil)

		failed := s != nil || err == nil
		if failed {
			t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "nil", "err")
		}

		s.Close()
	}

	{

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

		// test client

		cl := &ari.Client{}
		s, err := NewServer(cl, &Options{
			URL: "nats://127.0.0.1:4333",
		})

		failed := s == nil || err != nil
		if failed {
			t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
		}

		s.Close()
	}

}

func TestListen(t *testing.T) {

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

	// test client

	cl := &ari.Client{}
	s, err := NewServer(cl, &Options{
		URL: "nats://127.0.0.1:4333",
	})

	failed := s == nil || err != nil
	if failed {
		t.Errorf("natsgw.NewServer(cl, nil) => {%v, %v}, expected {%v, %v}", s, err, "cl", "nil")
	}

	go s.Listen()

	<-time.After(4 * time.Second)

	s.Close()
}
