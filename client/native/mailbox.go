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
func (m *Mailbox) Get(key *ari.Key) *ari.MailboxHandle {
	return ari.NewMailboxHandle(key, m)
}

// List lists the mailboxes and returns a list of handles
func (m *Mailbox) List(filter *ari.Key) (mx []*ari.Key, err error) {

	mailboxes := []struct {
		Name string `json:"name"`
	}{}
	if filter == nil {
		filter = ari.NodeKey(m.client.node, m.client.ApplicationName())
	}

	err = m.client.get("/mailboxes", &mailboxes)
	for _, i := range mailboxes {
		k := ari.NewKey(ari.MailboxKey, i.Name, ari.WithApp(m.client.ApplicationName()), ari.WithNode(m.client.node))
		if filter.Match(k) {
			mx = append(mx, k)
		}
	}

	return
}

// Data retrieves the state of the given mailbox
func (m *Mailbox) Data(key *ari.Key) (md *ari.MailboxData, err error) {
	md = &ari.MailboxData{}
	name := key.ID
	err = m.client.get("/mailboxes/"+name, &md)
	if err != nil {
		md = nil
		err = dataGetError(err, "mailbox", "%v", name)
	}
	return
}

// Update updates the new and old message counts of the mailbox
func (m *Mailbox) Update(key *ari.Key, oldMessages int, newMessages int) (err error) {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}
	name := key.ID
	err = m.client.put("/mailboxes/"+name, nil, &req)
	return err
}

// Delete deletes the mailbox
func (m *Mailbox) Delete(key *ari.Key) (err error) {
	name := key.ID
	err = m.client.del("/mailboxes/"+name, nil, "")
	return
}
