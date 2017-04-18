package ari

// Mailbox is the communication path to an Asterisk server for
// operating on mailbox resources
type Mailbox interface {

	// Get gets a handle to the mailbox for further operations
	Get(key *Key) MailboxHandle

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
	Name        string `json:"name"`
	NewMessages int    `json:"new_messages"` // Number of new (unread) messages
	OldMessages int    `json:"old_messages"` // Number of old (read) messages
}

// MailboxHandle is a wrapper for interacting with a particular mailbox
type MailboxHandle interface {

	// ID returns the identifier for the mailbox handle
	ID() string

	// Data gets the current state of the mailbox
	Data() (*MailboxData, error)

	// Update updates the state of the mailbox, or creates if does not exist
	Update(oldMessages int, newMessages int) error

	// Delete deletes the mailbox
	Delete() error
}
