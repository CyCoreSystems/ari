package native

import (
	"strconv"

	"github.com/CyCoreSystems/ari"
)

type nativeMailbox struct {
	conn *Conn
}

func (m *nativeMailbox) Get(name string) *ari.MailboxHandle {
	return ari.NewMailboxHandle(name, m)
}

func (m *nativeMailbox) List() (mx []*ari.MailboxHandle, err error) {

	mailboxes := []struct {
		Name string `json:"name"`
	}{}

	err = Get(m.conn, "/mailboxes", &mailboxes)
	for _, i := range mailboxes {
		mx = append(mx, m.Get(i.Name))
	}

	return
}

func (m *nativeMailbox) Data(name string) (md ari.MailboxData, err error) {
	err = Get(m.conn, "/mailboxes/"+name, &md)
	return
}

func (m *nativeMailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	err = Put(m.conn, "/mailboxes/"+name, nil, &req)
	return err
}

func (m *nativeMailbox) Delete(name string) (err error) {
	err = Delete(m.conn, "/mailboxes/"+name, nil, "")
	return
}
