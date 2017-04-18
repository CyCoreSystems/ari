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
func (m *Mailbox) Get(key *ari.Key) ari.MailboxHandle {
	return NewMailboxHandle(key, m)
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

// A MailboxHandle is a handle to a mailbox instance attached to an
// ari transport
type MailboxHandle struct {
	key *ari.Key
	m   *Mailbox
}

// NewMailboxHandle creates a new mailbox handle given the name and mailbox transport
func NewMailboxHandle(key *ari.Key, m *Mailbox) *MailboxHandle {
	return &MailboxHandle{
		key: key,
		m:   m,
	}
}

// ID returns the identifier for the mailbox handle
func (mh *MailboxHandle) ID() string {
	return mh.key.ID
}

// Data gets the current state of the mailbox
func (mh *MailboxHandle) Data() (*ari.MailboxData, error) {
	return mh.m.Data(mh.key)
}

// Update updates the state of the mailbox, or creates if does not exist
func (mh *MailboxHandle) Update(oldMessages int, newMessages int) error {
	return mh.m.Update(mh.key, oldMessages, newMessages)
}

// Delete deletes the mailbox
func (mh *MailboxHandle) Delete() error {
	return mh.m.Delete(mh.key)
}
