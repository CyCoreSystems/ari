package ari

// Logging represents a communication path to an
// Asterisk server for working with logging resources
type Logging interface {

	// Create creates a new log
	Create(name string, level string) error

	// List the logs
	List() ([]LogData, error)

	// Rotate rotates the log
	Rotate(name string) error

	// Delete deletes the log
	Delete(name string) error
}

// LogData represents the log data
type LogData struct {
	Name          string `json:"channel"`
	Configuration string `json:"configuration"`
	Type          string `json:"type"`
	Status        string `json:"configuration"`
}
