package ari

import "time"

// DTMFOptions is the list of pptions for DTMF sending
type DTMFOptions struct {
	Before   time.Duration
	Between  time.Duration
	Duration time.Duration
	After    time.Duration
}

// DTMFSender is an object which can be send DTMF signals
type DTMFSender interface {
	SendDTMF(dtmf string, opts *DTMFOptions)
}
