package ari

// Mailbox is the communication path to an Asterisk server for
// operating on mailbox resources
type Mailbox interface {

	// Get gets a handle to the mailbox for further operations
	Get(name string) *MailboxHandle

	// List lists the mailboxes in asterisk
	List() ([]*MailboxHandle, error)

	// Data gets the current state of the mailbox
	Data(name string) (MailboxData, error)

	// Update updates the state of the mailbox, or creates if does not exist
	Update(name string, oldMessages int, newMessages int) error

	// Delete deletes the mailbox
	Delete(name string) error
}

// MailboxData respresents the state of an Asterisk (voice) mailbox
type MailboxData struct {
	Name        string `json:"name"`
	NewMessages int    `json:"new_messages"` // Number of new (unread) messages
	OldMessages int    `json:"old_messages"` // Number of old (read) messages
}

// NewMailboxHandle creates a new mailbox handle given the name and mailbox transport
func NewMailboxHandle(name string, m Mailbox) *MailboxHandle {
	return &MailboxHandle{
		name: name,
		m:    m,
	}
}

// A MailboxHandle is a handle to a mailbox instance attached to an
// ari transport
type MailboxHandle struct {
	name string
	m    Mailbox
}

// ID returns the identifier for the mailbox handle
func (mh *MailboxHandle) ID() string {
	return mh.name
}

// Data gets the current state of the mailbox
func (mh *MailboxHandle) Data() (MailboxData, error) {
	return mh.m.Data(mh.name)
}

// Update updates the state of the mailbox, or creates if does not exist
func (mh *MailboxHandle) Update(oldMessages int, newMessages int) error {
	return mh.m.Update(mh.name, oldMessages, newMessages)
}

// Delete deletes the mailbox
func (mh *MailboxHandle) Delete() error {
	return mh.m.Delete(mh.name)
}
