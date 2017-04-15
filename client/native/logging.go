package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

// Logging provides the ARI Logging accessors for a native client
type Logging struct {
	client *Client
}

// Create creates a logging level
func (l *Logging) Create(name, level string) (err error) {
	type request struct {
		Configuration string `json:"configuration"`
	}
	req := request{level}
	err = l.client.post("/asterisk/logging/"+name, nil, &req)
	return
}

// Get returns a logging channel handle
func (l *Logging) Get(name string) ari.LogHandle {
	return &LogHandle{
		name: name,
	}
}

func (l *Logging) getLoggingChannels() ([]*ari.LogData, error) {
	var ld []*ari.LogData
	err := l.client.get("/asterisk/logging", &ld)
	return ld, err
}

// Data returns the data of a logging channel
func (l *Logging) Data(name string) (*ari.LogData, error) {
	ld, err := l.getLoggingChannels()
	if err != nil {
		return nil, err
	}

	for _, i := range ld {
		return i, nil
	}
	return nil, errors.New("not found")
}

// List lists the logging entities
func (l *Logging) List() ([]ari.LogHandle, error) {
	ld, err := l.getLoggingChannels()
	if err != nil {
		return nil, err
	}

	var ret []ari.LogHandle
	for _, i := range ld {
		ret = append(ret, l.Get(i.Name))
	}
	return ret, nil
}

// Rotate rotates the given log
func (l *Logging) Rotate(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to rotate unnamed channels")
		return
	}
	err = l.client.put("/asterisk/logging/"+name+"/rotate", nil, nil)
	return
}

// Delete deletes the named log
func (l *Logging) Delete(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to delete unnamed channels")
		return
	}
	err = l.client.del("/asterisk/logging/"+name, nil, "")
	return
}

// LogHandle provides an interface to manipulate a logging channel
type LogHandle struct {
	name string
	c    *Logging
}

// ID returns the ID (name) of the logging channel
func (l *LogHandle) ID() string {
	return l.name
}

// Data returns the data for the logging channel
func (l *LogHandle) Data() (*ari.LogData, error) {
	return l.c.Data(l.name)
}

// Rotate causes the logging channel's logfiles to be rotated
func (l *LogHandle) Rotate() error {
	return l.c.Rotate(l.name)
}

// Delete removes the logging channel from Asterisk
func (l *LogHandle) Delete() error {
	return l.c.Delete(l.name)
}
