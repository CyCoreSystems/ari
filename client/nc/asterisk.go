package nc

import (
	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

type natsAsterisk struct {
	conn *nats.Conn
}

type natsAsteriskVariables struct {
	conn *nats.Conn
}

func (a *natsAsterisk) ReloadModule(name string) (err error) {
	err = request(a.conn, "ari.asterisk.reload."+name, nil, nil)
	return
}

func (a *natsAsterisk) Info(only string) (ai *ari.AsteriskInfo, err error) {
	ai = &ari.AsteriskInfo{}
	err = request(a.conn, "ari.asterisk.info", only, ai)
	return
}

func (a *natsAsterisk) Variables() ari.Variables {
	return &natsAsteriskVariables{a.conn}
}

func (a *natsAsteriskVariables) Get(variable string) (ret string, err error) {
	err = request(a.conn, "ari.asterisk.variables.get."+variable, nil, &ret)
	return
}

func (a *natsAsteriskVariables) Set(variable string, value string) (err error) {
	err = request(a.conn, "ari.asterisk.variables.set."+variable, value, nil)
	return
}
