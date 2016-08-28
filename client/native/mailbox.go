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

func (m *nativeMailbox) Data(name string) (md ari.MailboxData, err error) {
	err = Get(m.conn, "/mailboxes", &md)
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
