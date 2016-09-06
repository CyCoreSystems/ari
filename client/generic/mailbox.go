package generic

import (
	"strconv"

	"github.com/CyCoreSystems/ari"
)

type Mailbox struct {
	Conn Conn
}

func (m *Mailbox) Get(name string) *ari.MailboxHandle {
	return ari.NewMailboxHandle(name, m)
}

func (m *Mailbox) List() (mx []*ari.MailboxHandle, err error) {

	mailboxes := []struct {
		Name string `json:"name"`
	}{}

	err = m.Conn.Get("/mailboxes", nil, &mailboxes)
	for _, i := range mailboxes {
		mx = append(mx, m.Get(i.Name))
	}

	return
}

func (m *Mailbox) Data(name string) (md ari.MailboxData, err error) {
	err = m.Conn.Get("/mailboxes/%s", []interface{}{name}, &md)
	return
}

func (m *Mailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	err = m.Conn.Put("/mailboxes/%s", []interface{}{name}, nil, &req)
	return err
}

func (m *Mailbox) Delete(name string) (err error) {
	err = m.Conn.Delete("/mailboxes/%s", []interface{}{name}, nil, "")
	return
}
