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

func TestMailboxList(t *testing.T) {

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

	mockMailbox := mock.NewMockMailbox(ctrl)
	mb1 := ari.NewMailboxHandle("mb1", mockMailbox)
	mb2 := ari.NewMailboxHandle("mb2", mockMailbox)

	mockMailbox.EXPECT().List().Return([]*ari.MailboxHandle{mb1, mb2}, nil)
	mockMailbox.EXPECT().List().Return([]*ari.MailboxHandle{}, errors.New("Failed getting mailbox list"))

	cl := &ari.Client{
		Mailbox: mockMailbox,
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
		l, err := natsClient.Mailbox.List()
		failed = err != nil
		failed = failed || len(l) != 2

		if failed {
			t.Errorf("nc.Mailbox.List() => '%v', '%v'; expected '%v', %v'",
				l, err,
				"[mb1, mb2]", nil)
		}
	}

	{
		l, err := natsClient.Mailbox.List()
		failed = err == nil || err.Error() != "Failed getting mailbox list"
		if failed {
			t.Errorf("nc.Mailbox.List() => '%v', '%v'; expected '%v', %v'",
				l, err,
				"[]", "Failed getting mailbox list")
		}
	}

}

func TestMailboxData(t *testing.T) {

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

	var mailboxData ari.MailboxData

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMailbox := mock.NewMockMailbox(ctrl)

	mockMailbox.EXPECT().Data("mb1").Return(mailboxData, nil)
	mockMailbox.EXPECT().Data("mb2").Return(mailboxData, errors.New("Failed to get mailbox data"))

	cl := &ari.Client{
		Mailbox: mockMailbox,
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
		d, err := natsClient.Mailbox.Data("mb1")
		failed = err != nil
		if failed {
			t.Errorf("nc.Mailbox.Data('%v') => '%v', '%v'; expected '%v', %v'",
				"mb1",
				d, err,
				mailboxData, nil)
		}
	}

	{
		d, err := natsClient.Mailbox.Data("mb2")
		failed = err == nil || err.Error() != "Failed to get mailbox data"
		if failed {
			t.Errorf("nc.Mailbox.Data('%v') => '%v', '%v'; expected '%v', %v'",
				"mb2",
				d, err,
				mailboxData, "Failed to get mailbox data")
		}
	}

}

func TestMailboxUpdate(t *testing.T) {

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

	mockMailbox := mock.NewMockMailbox(ctrl)

	mockMailbox.EXPECT().Update("mb1", 1, 2).Return(nil)

	cl := &ari.Client{
		Mailbox: mockMailbox,
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
		err := natsClient.Mailbox.Update("mb1", 1, 2)
		failed = err != nil
		if failed {
			t.Errorf("nc.Mailbox.Update('%v', 1, 2) => '%v'; expected '%v'",
				"mb1",
				err,
				nil)
		}
	}

}

func TestMailboxDelete(t *testing.T) {

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

	mockMailbox := mock.NewMockMailbox(ctrl)
	mockMailbox.EXPECT().Delete("mb1").Return(nil)

	cl := &ari.Client{
		Mailbox: mockMailbox,
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
		err := natsClient.Mailbox.Delete("mb1")
		failed = err != nil
		if failed {
			t.Errorf("nc.Mailbox.Delete('%v') => '%v'; expected '%v'",
				"mb1",
				err,
				nil)
		}
	}

}
