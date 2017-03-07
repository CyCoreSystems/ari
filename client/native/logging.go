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
	err = l.client.conn.Post("/asterisk/logging/"+name, nil, &req)
	return
}

// List lists the logging entities
func (l *Logging) List() (ld []ari.LogData, err error) {
	err = l.client.conn.Get("/asterisk/logging", &ld)
	return
}

// Rotate rotates the given log
func (l *Logging) Rotate(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to rotate unnamed channels")
		return
	}
	err = l.client.conn.Put("/asterisk/logging/"+name+"/rotate", nil, nil)
	return
}

// Delete deletes the named log
func (l *Logging) Delete(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to delete unnamed channels")
		return
	}
	err = l.client.conn.Delete("/asterisk/logging/"+name, nil, "")
	return
}
