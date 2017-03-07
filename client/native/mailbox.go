package native

import (
	"strconv"

	"github.com/CyCoreSystems/ari"
)

type nativeMailbox struct {
	client *Client
}

func (m *nativeMailbox) Get(name string) *ari.MailboxHandle {
	return ari.NewMailboxHandle(name, m)
}

func (m *nativeMailbox) List() (mx []*ari.MailboxHandle, err error) {

	mailboxes := []struct {
		Name string `json:"name"`
	}{}

	err = m.client.conn.Get("/mailboxes", &mailboxes)
	for _, i := range mailboxes {
		mx = append(mx, m.Get(i.Name))
	}

	return
}

func (m *nativeMailbox) Data(name string) (md ari.MailboxData, err error) {
	err = m.client.conn.Get("/mailboxes/"+name, &md)
	return
}

func (m *nativeMailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	err = m.client.conn.Put("/mailboxes/"+name, nil, &req)
	return err
}

func (m *nativeMailbox) Delete(name string) (err error) {
	err = m.client.conn.Delete("/mailboxes/"+name, nil, "")
	return
}
