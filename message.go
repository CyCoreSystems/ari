package ari

import (
	"encoding/json"
	"errors"
)

// Message is the first extension of the RawMessage type,
// containing only a Type
type Message struct {
	RawMessage
	Type string `json:"type"`
}

// NewMessage constructs a Message from a byte slice
func NewMessage(raw []byte) (*Message, error) {
	var m Message

	raw = append(raw, '\n')

	err := json.Unmarshal(raw, &m)
	if err != nil {
		return &m, err
	}

	// Set _raw to be our raw bytestream
	m._raw = &raw

	return &m, nil
}

// MessageRawer provides operations to get raw message data
type MessageRawer interface {
	SetRaw(*[]byte)
	GetRaw() *[]byte
}

// RawMessage contains the raw bytes
type RawMessage struct {
	_raw *[]byte // The raw message
}

// DecodeAs converts the current message to
// a new message type
func (m *RawMessage) DecodeAs(v MessageRawer) error {
	if v == nil {
		return errors.New("empty message")
	}
	// First, unmarshal raw into the new type
	err := json.Unmarshal(*m.GetRaw(), v)
	if err != nil {
		return err
	}

	// Set the new raw to the old raw
	v.SetRaw(m.GetRaw())

	return nil
}

// SetRaw sets the raw value of this RawMessage
func (m *RawMessage) SetRaw(raw *[]byte) {
	m._raw = raw
}

// GetRaw gets the raw value of this RawMessage
func (m *RawMessage) GetRaw() *[]byte {
	return m._raw
}
