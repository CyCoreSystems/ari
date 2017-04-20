package ari

// Mailbox is the communication path to an Asterisk server for
// operating on mailbox resources
type Mailbox interface {

	// Get gets a handle to the mailbox for further operations
	Get(key *Key) *MailboxHandle

	// List lists the mailboxes in asterisk
	List(filter *Key) ([]*Key, error)

	// Data gets the current state of the mailbox
	Data(key *Key) (*MailboxData, error)

	// Update updates the state of the mailbox, or creates if does not exist
	Update(key *Key, oldMessages int, newMessages int) error

	// Delete deletes the mailbox
	Delete(key *Key) error
}

// MailboxData respresents the state of an Asterisk (voice) mailbox
type MailboxData struct {
	// Key is the cluster-unique identifier for this mailbox
	Key *Key `json:"key"`

	Name        string `json:"name"`
	NewMessages int    `json:"new_messages"` // Number of new (unread) messages
	OldMessages int    `json:"old_messages"` // Number of old (read) messages
}

// A MailboxHandle is a handle to a mailbox instance attached to an
// ari transport
type MailboxHandle struct {
	key *Key
	m   Mailbox
}

// NewMailboxHandle creates a new mailbox handle given the name and mailbox transport
func NewMailboxHandle(key *Key, m Mailbox) *MailboxHandle {
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
func (mh *MailboxHandle) Data() (*MailboxData, error) {
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
