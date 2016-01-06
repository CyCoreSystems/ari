package ari

import (
	"strconv"
	"time"
)

// ContactInfo describes a Contact
type ContactInfo struct {
	AOR       string `json:"aor"`                      // Address of Record for this contact
	Status    string `json:"contact_status"`           // Current status of the contact
	Roundtrip string `json:"roundtrip_usec,omitempty"` // Current round trip time, in microseconds
	URI       string `json:"uri"`                      // The location URI of the contact
}

// Latency returns the round-trip time for this contact
func (c *ContactInfo) Latency() time.Duration {
	if ts, err := strconv.Atoi(c.Roundtrip); err == nil {
		return time.Duration(ts) * time.Microsecond
	}
	return 0
}
