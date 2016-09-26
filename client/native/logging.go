package native

import (
	"errors"

	"github.com/CyCoreSystems/ari"
)

type nativeLogging struct {
	conn *Conn
}

func (l *nativeLogging) Create(name, level string) (err error) {
	type request struct {
		Configuration string `json:"configuration"`
	}
	req := request{level}
	err = Post(l.conn, "/asterisk/logging/"+name, nil, &req)
	return
}

func (l *nativeLogging) List() (ld []ari.LogData, err error) {
	err = Get(l.conn, "/asterisk/logging", &ld)
	return
}

func (l *nativeLogging) Rotate(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to rotate unnamed channels")
		return
	}
	err = Put(l.conn, "/asterisk/logging/"+name+"/rotate", nil, nil)
	return
}

func (l *nativeLogging) Delete(name string) (err error) {
	if name == "" {
		err = errors.New("Not allowed to delete unnamed channels")
		return
	}
	err = Delete(l.conn, "/asterisk/logging/"+name, nil, "")
	return
}
