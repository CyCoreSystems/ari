package ari

// Logging represents a communication path to an
// Asterisk server for working with logging resources
type Logging interface {

	// Create creates a new log
	Create(name string, level string) error

	// Data retrives the data for a logging channel
	Data(name string) (*LogData, error)

	// Data retrives the data for a logging channel
	Get(name string) LogHandle

	// List the logs
	List() ([]LogHandle, error)

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
	Status        string `json:"status"`
}

// LogHandle provides a wrapper to a Logging channel to perform operations on that Logging channel.
type LogHandle interface {
	// ID returns the identifier for the logging channel
	ID() string

	// Data retrives the data for the logging channel
	Data() (*LogData, error)

	// Rotate rotates the log file for this channel
	Rotate() error

	// Delete removes this logging channel
	Delete() error
}
