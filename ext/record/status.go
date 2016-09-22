package record

//go:generate stringer -type=Status status.go

// Status is the indicator for recording status
type Status int64

const (
	// InProgress indicates that a recording is still in progress
	InProgress Status = iota

	// Canceled indicates that a recording was canceled (by request)
	Canceled

	// Failed indicates that a recording failed
	Failed

	// Finished indicates that a recording finished normally
	Finished

	// Hangup indicates that a recording was ended due to hangup
	Hangup
)
