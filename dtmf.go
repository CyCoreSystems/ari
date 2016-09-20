package ari

// DTMFSender is an object which can be send DTMF signals
type DTMFSender interface {
	SendDTMF(dtmf string)
}
