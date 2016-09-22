package nc

import "github.com/CyCoreSystems/ari"

type natsAsterisk struct {
	conn *Conn
}

type natsAsteriskVariables struct {
	conn *Conn
}

// tests and other advanced utility functions can cast to an interface to get the NatsConnection object out
func (a *natsAsterisk) NatsConnection() *Conn {
	return a.conn
}

func (a *natsAsterisk) ReloadModule(name string) (err error) {
	err = a.conn.standardRequest("ari.asterisk.reload."+name, nil, nil)
	return
}

func (a *natsAsterisk) Info(only string) (*ari.AsteriskInfo, error) {
	ai := &ari.AsteriskInfo{}
	err := a.conn.readRequest("ari.asterisk.info", only, ai)
	if err != nil {
		return nil, err
	}
	return ai, nil
}

func (a *natsAsterisk) Variables() ari.Variables {
	return &natsAsteriskVariables{a.conn}
}

func (a *natsAsteriskVariables) Get(variable string) (ret string, err error) {
	err = a.conn.readRequest("ari.asterisk.variables.get."+variable, nil, &ret)
	return
}

func (a *natsAsteriskVariables) Set(variable string, value string) (err error) {
	err = a.conn.standardRequest("ari.asterisk.variables.set."+variable, value, nil)
	return
}
