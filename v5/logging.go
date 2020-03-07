package ari

// Logging represents a communication path to an
// Asterisk server for working with logging resources
type Logging interface {

	// Create creates a new log.  The levels are a comma-separated list of
	// logging levels on which this channel should operate.  The name of the
	// channel should be the key's ID.
	Create(key *Key, levels string) (*LogHandle, error)

	// Data retrives the data for a logging channel
	Data(key *Key) (*LogData, error)

	// Data retrives the data for a logging channel
	Get(key *Key) *LogHandle

	// List the logs
	List(filter *Key) ([]*Key, error)

	// Rotate rotates the log
	Rotate(key *Key) error

	// Delete deletes the log
	Delete(key *Key) error
}

// LogData represents the log data
type LogData struct {
	// Key is the cluster-unique identifier for this logging channel
	Key *Key `json:"key"`

	// Name is the name of the logging channel
	Name string `json:"channel"`

	// Levels is a comma-separated list of logging levels for this channel
	Levels string `json:"levels"`

	// Type indicates the type of logs for this channel
	Types string `json:"types"`

	// Status indicates whether this logging channel is enabled
	Status string `json:"status"`
}

// NewLogHandle builds a new log handle given the `Key` and `Logging`` client
func NewLogHandle(key *Key, l Logging) *LogHandle {
	return &LogHandle{
		key: key,
		c:   l,
	}
}

// LogHandle provides an interface to manipulate a logging channel
type LogHandle struct {
	key *Key
	c   Logging
}

// ID returns the ID (name) of the logging channel
func (l *LogHandle) ID() string {
	return l.key.ID
}

// Key returns the Key of the logging channel
func (l *LogHandle) Key() *Key {
	return l.key
}

// Data returns the data for the logging channel
func (l *LogHandle) Data() (*LogData, error) {
	return l.c.Data(l.key)
}

// Rotate causes the logging channel's logfiles to be rotated
func (l *LogHandle) Rotate() error {
	return l.c.Rotate(l.key)
}

// Delete removes the logging channel from Asterisk
func (l *LogHandle) Delete() error {
	return l.c.Delete(l.key)
}
