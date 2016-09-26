package nc

import "github.com/CyCoreSystems/ari"

type natsLogging struct {
	conn *Conn
}

func (l *natsLogging) Create(name string, level string) (err error) {
	err = l.conn.standardRequest("ari.logging.create."+name, level, nil)
	return
}

func (l *natsLogging) List() (ld []ari.LogData, err error) {
	err = l.conn.readRequest("ari.logging.all", nil, &ld)
	return
}

func (l *natsLogging) Rotate(name string) (err error) {
	err = l.conn.standardRequest("ari.logging.rotate."+name, nil, nil)
	return
}

func (l *natsLogging) Delete(name string) (err error) {
	err = l.conn.standardRequest("ari.logging.delete."+name, nil, nil)
	return
}
