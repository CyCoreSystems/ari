package nc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

// Conn is the wrapper type for a nats connnection along some ARI specific options
type Conn struct {
	opts Options
	conn *nats.Conn
}

// a read request is a request that is a read operation. It is less about
// "read" operations and more about operations which are repeatable/idempotent.
func (c *Conn) readRequest(subj string, body interface{}, dest interface{}) (err error) {

	maxRetries := c.opts.ReadOperationRetryCount

	if maxRetries == 0 {
		maxRetries = 1
	}

	for i := 0; i != maxRetries; i++ {

		err = c.standardRequest(subj, body, dest)
		if err == nil {
			return
		}

		if t, ok := err.(temp); !ok || !t.Temporary() {
			return
		}
	}

	return
}

// a standard request
func (c *Conn) standardRequest(subj string, body interface{}, dest interface{}) (err error) {

	conn := c.conn

	// build json request

	data := []byte("{}")

	if data != nil {
		if data, err = json.Marshal(body); err != nil {
			return
		}
	}

	// prepare response channel

	var sub *nats.Subscription

	replyID := uuid.NewV1().String()
	ch := make(chan *nats.Msg, 2)
	sub, err = conn.ChanSubscribe(replyID+".>", ch)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	// send request

	if err = conn.PublishRequest(subj, replyID, data); err != nil {
		return
	}

	// listen for response

	var msg *nats.Msg
	select {
	case <-time.After(c.opts.RequestTimeout):
		err = timeoutErr("Timeout waiting for response")
		return
	case msg = <-ch:
	}

	// handle error sent from server

	msgType := msg.Subject[len(replyID)+1:]

	if msgType == "err" {
		err = errors.New(string(msg.Data))
		return
	}

	// write into destination, if it isn't nil

	if dest != nil {
		err = json.Unmarshal(msg.Data, dest)
	}

	return
}

type temp interface {
	Temporary() bool
}

type timeoutErr string

func (err timeoutErr) Error() string {
	return string(err)
}

func (err timeoutErr) Timeout() bool {
	return true
}

func (err timeoutErr) Temporary() bool {
	return true
}
