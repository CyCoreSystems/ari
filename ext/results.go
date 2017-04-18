package ext

// Result describes the result of an operation
type Result struct {
	// Data is the received data (digits) during the prompt or queue
	Data string

	// Status is the status of the prompt play
	Status Status
}
