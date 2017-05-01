package ari

// Message is the first extension of the RawMessage type,
// containing only a Type
type Message struct {
	Type string `json:"type"`
}
