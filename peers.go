package ari

// Peer describes a remote peer that communicates with Asterisk
type Peer struct {
	Address string       `json:"address,omitempty"` // IP Address of the peer
	Cause   string       `json:"cause,omitempty"`   // Reason associated with change in peer_status
	Status  string       `json:"peer_status"`       // The current state of the peer
	Port    string       `json:"port,omitempty"`    // The port of the peer
	Time    AsteriskDate `json:"time,omitempty"`    // The last known time the peer was contacted
}
