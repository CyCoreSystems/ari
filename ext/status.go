package ext

//go:generate stringer --type=Status status.go

// Status indicates the status of an operation
type Status int

const (
	// Incomplete indicates that the operation is incomplete.
	Incomplete Status = iota

	// Complete indicates that the operation has completed.
	Complete

	// Canceled indicates that the operation has been cancelled due to context cancellation.
	Canceled

	// Timeout indicates that the operation has timed out.
	Timeout

	// Invalid indicates that a match cannot be found from the digits received.
	// No more digits should be received.
	Invalid

	// Hangup indicates that the operation was interrupted due to a hangup.
	Hangup

	// DTMFInterrupt indicates that the operation was interrupted due to a DTMF.
	DTMFInterrupt

	// Failed indicates that the operation has failed.
	Failed
)
