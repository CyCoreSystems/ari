package ari

import (
	"encoding/json"
	"fmt"
	"reflect"

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

	e, ok := v.(Message)
	if !ok {
		return fmt.Errorf("Cannot cast receiver to a Message", "type", reflect.TypeOf(v))
	}

	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	// Store the raw data
	e.__raw = &data

	return nil
}

// AsteriskCode is a websocket Codec for Asterisk messages
var AsteriskCodec = websocket.Codec{
	Marshal:   marshal,
	Unmarshal: unmarshal,
}

type MessageRawer interface {
	SetRaw(*[]byte)
	GetRaw() *[]byte
}

type RawMessage struct {
	__raw *[]byte // The raw message
}

// DecodeAs converts the current message to
// a new message type
func (m *RawMessage) DecodeAs(v MessageRawer) error {
	// First, unmarshal raw into the new type
	err := json.Unmarshal(*m.GetRaw(), v)
	if err != nil {
		return err
	}

	// Set the new raw to the old raw
	v.SetRaw(m.GetRaw())

	return nil
}

// Set the __raw value of this RawMessage
func (m *RawMessage) SetRaw(raw *[]byte) {
	m.__raw = raw
}

// Get the __raw value of this RawMessage
func (m *RawMessage) GetRaw() *[]byte {
	return m.__raw
}

// Message is the first extension of the RawMessage type,
// containing only a Type
type Message struct {
	RawMessage
	Type string `json:"type"`
}

// Construct a Message from a byte slice
func NewMessage(raw []byte) (*Message, error) {
	var m Message

	raw = append(raw, '\n')

	err := json.Unmarshal(raw, &m)
	if err != nil {
		Logger.Error("Failed to unmarshal new message", err.Error())
		return &m, err
	}

	// Set __raw to be our raw bytestream
	m.__raw = &raw

	return &m, nil
}
