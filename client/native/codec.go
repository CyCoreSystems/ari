package native

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/CyCoreSystems/ari"

	"golang.org/x/net/websocket"
)

// Marshal is a no-op to implement websocket.Codec.  Asterisk
// websocket connections should never have the client send any data
func marshal(v interface{}) (data []byte, payloadType byte, err error) {
	return
}

// Unmarshal implements websocket.Codec
func unmarshal(data []byte, payloadType byte, v interface{}) error {
	data = append(data, '\n')

	e, ok := v.(*ari.Message)
	if !ok {
		return fmt.Errorf("Cannot cast receiver to a Message when it is of type %v", reflect.TypeOf(v))
	}

	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	// Store the raw data
	e.SetRaw(&data)

	return nil
}

// AsteriskCodec is a websocket Codec for Asterisk messages
var AsteriskCodec = websocket.Codec{
	Marshal:   marshal,
	Unmarshal: unmarshal,
}
