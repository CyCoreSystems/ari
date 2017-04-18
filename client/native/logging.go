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
func (l *Logging) Get(key *ari.Key) *ari.LogHandle {
	return ari.NewLogHandle(key, l)
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
