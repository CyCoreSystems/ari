package audio

//go:generate stringer --type=Status status.go

// Status indicates the status of the prompt
// operation.
type Status int

const (
	// InProgress indicates that the audio is currently playing
	InProgress Status = iota

	// Finished indicates that the audio playback finished successfully
	Finished

	// Canceled indicates that the audio was canceled
	Canceled

	// Timeout indicates that audio playback timed out
	Timeout

	// Hangup indicates that the audio was interrupted by hangup
	Hangup

	// Failed indicates that the audio playback failed
	Failed
)
