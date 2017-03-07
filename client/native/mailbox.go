package native

import (
	"strconv"

	"github.com/CyCoreSystems/ari"
)

// Mailbox provides the ARI Mailbox accessors for the native client
type Mailbox struct {
	client *Client
}

// Get gets a lazy handle for the mailbox name
func (m *Mailbox) Get(name string) *ari.MailboxHandle {
	return ari.NewMailboxHandle(name, m)
}

// List lists the mailboxes and returns a list of handles
func (m *Mailbox) List() (mx []*ari.MailboxHandle, err error) {

	mailboxes := []struct {
		Name string `json:"name"`
	}{}

	err = m.client.conn.Get("/mailboxes", &mailboxes)
	for _, i := range mailboxes {
		mx = append(mx, m.Get(i.Name))
	}

	return
}

// Data retrieves the state of the given mailbox
func (m *Mailbox) Data(name string) (md ari.MailboxData, err error) {
	err = m.client.conn.Get("/mailboxes/"+name, &md)
	return
}

// Update updates the new and old message counts of the mailbox
func (m *Mailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	err = m.client.conn.Put("/mailboxes/"+name, nil, &req)
	return err
}

// Delete deletes the mailbox
func (m *Mailbox) Delete(name string) (err error) {
	err = m.client.conn.Delete("/mailboxes/"+name, nil, "")
	return
}
