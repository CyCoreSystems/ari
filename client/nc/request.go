package nc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

var request = standardRequest

func standardRequest(conn *nats.Conn, subj string, body interface{}, dest interface{}) (err error) {

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
	case <-time.After(DefaultRequestTimeout):
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
