package testutils

import (
	"errors"
	"os/exec"
	"syscall"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/nc"

	"golang.org/x/net/context"
)

// LaunchNatsServer launches a nats server
func LaunchNatsServer() (*ari.Client, context.CancelFunc, error) {

	ctx, cancel := context.WithCancel(context.Background())

	bin, err := exec.LookPath("gnatsd")
	if err != nil {
		return nil, cancel, errors.New("Not found")
	}

	cmd := exec.Command(bin, "-p", "4333")
	if err := cmd.Start(); err != nil {
		return nil, cancel, err
	}

	go func() {
		<-ctx.Done()
		cmd.Process.Signal(syscall.SIGTERM)
	}()

	<-time.After(500 * time.Millisecond)

	cl, err := nc.New("nats://127.0.0.1:4333")
	if err != nil {
		cancel()
		return nil, cancel, nil
	}

	return cl, cancel, nil
}
