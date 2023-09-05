package native

import (
	"github.com/rotisserie/eris"

	"github.com/CyCoreSystems/ari/v6"
)

// Logging provides the ARI Logging accessors for a native client
type Logging struct {
	client *Client
}

// Create creates a logging level
func (l *Logging) Create(key *ari.Key, levels string) (*ari.LogHandle, error) {
	req := struct {
		Levels string `json:"configuration"`
	}{
		Levels: levels,
	}

	err := l.client.post("/asterisk/logging/"+key.ID, nil, &req)
	if err != nil {
		return nil, err
	}

	return l.Get(key), nil
}

// Get returns a logging channel handle
func (l *Logging) Get(key *ari.Key) *ari.LogHandle {
	return ari.NewLogHandle(l.client.stamp(key), l)
}

func (l *Logging) getLoggingChannels() ([]*ari.LogData, error) {
	var ld []*ari.LogData
	err := l.client.get("/asterisk/logging", &ld)

	return ld, err
}

// Data returns the data of a logging channel
func (l *Logging) Data(key *ari.Key) (*ari.LogData, error) {
	if key == nil || key.ID == "" {
		return nil, eris.New("logging key not supplied")
	}

	logChannels, err := l.getLoggingChannels()
	if err != nil {
		return nil, eris.Wrap(err, "failed to get list of logging channels")
	}

	for _, i := range logChannels {
		if i.Name == key.ID {
			i.Key = l.client.stamp(key)
			return i, nil
		}
	}

	return nil, eris.New("not found")
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
func (l *Logging) Rotate(key *ari.Key) error {
	name := key.ID
	if name == "" {
		return eris.New("Not allowed to rotate unnamed channels")
	}

	return l.client.put("/asterisk/logging/"+name+"/rotate", nil, nil)
}

// Delete deletes the named log
func (l *Logging) Delete(key *ari.Key) error {
	name := key.ID
	if name == "" {
		return eris.New("Not allowed to delete unnamed channels")
	}

	return l.client.del("/asterisk/logging/"+name, nil, "")
}
