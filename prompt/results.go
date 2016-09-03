package prompt

// Status indicates the status of the prompt
// operation.
type Status int

const (
	// Incomplete indicates that there are not enough digits to determine a match.
	Incomplete Status = iota

	// Complete means that a match was found from the digits received.
	Complete

	// Invalid indicates that a match cannot be found from the digits received.
	// No more digits should be received.
	Invalid

	// Canceled indicates that matching was canceled
	Canceled

	// Timeout indicates that the matching failed due to timeout
	Timeout

	// Hangup indicates that a match could not be made due to
	// the channel being hung up.
	Hangup

	// Failed indicates that matching could not be completed
	// due to a failure (usually an error)
	Failed
)

// Result describes the result of a prompt operation
type Result struct {
	// Data is the received data (digits) during the prompt
	Data string

	// Status is the status of the prompt play
	Status Status
}
