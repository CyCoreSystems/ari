package nc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

func request(conn *nats.Conn, subj string, body interface{}, dest interface{}) (err error) {

	// prepare response channel

	replyID := uuid.NewV1().String()
	ch := make(chan *nats.Msg, 2)
	var sub *nats.Subscription
	sub, err = conn.ChanSubscribe(subj+":>", ch)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	// build json request

	data := []byte("{}")

	if data != nil {
		if data, err = json.Marshal(body); err != nil {
			return
		}
	}

	// send request

	if err = conn.PublishRequest(subj, replyID, data); err != nil {
		return
	}

	// listen for response

	var msg *nats.Msg
	select {
	case <-time.After(DefaultRequestTimeout):
		err = errors.New("timeout") //TODO: timeout error type
		return
	case msg = <-ch:
	}

	// handle error sent from server

	if msg.Subject[len(subj)+1:] == "err" {
		err = errors.New(string(msg.Data))
		return
	}

	// write into destination, if it isn't nil

	if dest != nil {
		err = json.Unmarshal(msg.Data, dest)
	}

	return
}
