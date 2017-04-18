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
func (l *Logging) Create(key *ari.Key, level string) (err error) {
	type request struct {
		Configuration string `json:"configuration"`
	}
	req := request{level}
	name := key.ID
	err = l.client.post("/asterisk/logging/"+name, nil, &req)
	return
}

// Get returns a logging channel handle
func (l *Logging) Get(key *ari.Key) ari.LogHandle {
	return &LogHandle{
		key: key,
	}
}

func (l *Logging) getLoggingChannels() ([]*ari.LogData, error) {
	var ld []*ari.LogData
	err := l.client.get("/asterisk/logging", &ld)
	return ld, err
}

// Data returns the data of a logging channel
func (l *Logging) Data(key *ari.Key) (*ari.LogData, error) {
	ld, err := l.getLoggingChannels()
	if err != nil {
		return nil, err
	}

	for _, i := range ld {
		if i.Name == key.ID {
			return i, nil
		}
	}
	return nil, errors.New("not found")
}

// List lists the logging entities
func (l *Logging) List(filter *ari.Key) ([]*ari.Key, error) {
	ld, err := l.getLoggingChannels()
	if err != nil {
		return nil, err
	}
	if filter == nil {
		filter = ari.NodeKey(l.client.ApplicationName(), l.client.node)
	}

	var ret []*ari.Key
	for _, i := range ld {
		k := ari.NewKey(ari.LoggingKey, i.Name, ari.WithApp(l.client.ApplicationName()), ari.WithNode(l.client.node))
		if filter.Match(k) {
			ret = append(ret, k)
		}
	}
	return ret, nil
}

// Rotate rotates the given log
func (l *Logging) Rotate(key *ari.Key) (err error) {
	name := key.ID
	if name == "" {
		err = errors.New("Not allowed to rotate unnamed channels")
		return
	}
	err = l.client.put("/asterisk/logging/"+name+"/rotate", nil, nil)
	return
}

// Delete deletes the named log
func (l *Logging) Delete(key *ari.Key) (err error) {
	name := key.ID
	if name == "" {
		err = errors.New("Not allowed to delete unnamed channels")
		return
	}
	err = l.client.del("/asterisk/logging/"+name, nil, "")
	return
}

// LogHandle provides an interface to manipulate a logging channel
type LogHandle struct {
	key *ari.Key
	c   *Logging
}

// ID returns the ID (name) of the logging channel
func (l *LogHandle) ID() string {
	return l.key.ID
}

// Data returns the data for the logging channel
func (l *LogHandle) Data() (*ari.LogData, error) {
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
