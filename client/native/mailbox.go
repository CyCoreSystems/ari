package native

import (
	"errors"
	"strconv"

	"github.com/Amtelco-Software/ari/v6"
)

// Mailbox provides the ARI Mailbox accessors for the native client
type Mailbox struct {
	client *Client
}

// Get gets a lazy handle for the mailbox name
func (m *Mailbox) Get(key *ari.Key) *ari.MailboxHandle {
	return ari.NewMailboxHandle(m.client.stamp(key), m)
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
		k := m.client.stamp(ari.NewKey(ari.MailboxKey, i.Name))
		if filter.Match(k) {
			mx = append(mx, k)
		}
	}

	return
}

// Data retrieves the state of the given mailbox
func (m *Mailbox) Data(key *ari.Key) (*ari.MailboxData, error) {
	if key == nil || key.ID == "" {
		return nil, errors.New("mailbox key not supplied")
	}

	data := new(ari.MailboxData)
	if err := m.client.get("/mailboxes/"+key.ID, data); err != nil {
		return nil, dataGetError(err, "mailbox", "%v", key.ID)
	}

	data.Key = m.client.stamp(key)

	return data, nil
}

// Update updates the new and old message counts of the mailbox
func (m *Mailbox) Update(key *ari.Key, oldMessages int, newMessages int) error {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	return m.client.put("/mailboxes/"+key.ID, nil, &req)
}

// Delete deletes the mailbox
func (m *Mailbox) Delete(key *ari.Key) error {
	return m.client.del("/mailboxes/"+key.ID, nil, "")
}
