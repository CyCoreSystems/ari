package ari

import "encoding/json"

type MessageRawer interface {
	SetRaw(*[]byte)
	GetRaw() *[]byte
}

type RawMessage struct {
	__raw *[]byte `json:"-"` // The raw message
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
		Logger.Println("Failed to unmarshal new message", err.Error())
		return &m, err
	}

	// Set __raw to be our raw bytestream
	m.__raw = &raw

	return &m, nil
}
