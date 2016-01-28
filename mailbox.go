package ari

import "strconv"

// Mailbox respresents the state of an Asterisk (voice) mailbox
type Mailbox struct {
	Name         string `json:"name"`
	New_messages int    `json:"new_messages"` // Number of new (unread) messages
	Old_messages int    `json:"old_messages"` // Number of old (read) messages
}

//List all mailboxes on asterisk server
//Equivalent to GET /mailboxes
func (c *Client) ListMailboxes() ([]Mailbox, error) {
	var m []Mailbox
	err := c.Get("/mailboxes", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Retrieve the current state of a specific mailbox
//Equivalent to GET /mailboxes/{mailboxName}
func (c *Client) GetMailbox(mailboxName string) (Mailbox, error) {
	var m Mailbox
	err := c.Get("/mailboxes/"+mailboxName, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Change the state of a mailbox. (Note - implicitly creates the mailbox).
//Equivalent to PUT /mailboxes/{mailboxName}
func (c *Client) ChangeMailboxState(mailboxName string, oldMessages int, newMessages int) error {
	req := map[string]string{
		"oldMessages": strconv.Itoa(oldMessages),
		"newMessages": strconv.Itoa(newMessages),
	}

	//send request
	err := c.Put("/mailboxes/"+mailboxName, nil, &req)
	return err
}

//Request structure for changing a mailbox state. Both arguments are required.

//Destroy a mailbox.
//Equivalent to DELETE /mailboxes/{mailboxName}
func (c *Client) DeleteMailbox(mailboxName string) error {
	err := c.Delete("/mailboxes/"+mailboxName, nil, nil)
	return err
}
