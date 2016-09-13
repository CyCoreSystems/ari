package natsgw

import (
	"encoding/json"
)

func (srv *Server) device() {
	srv.subscribe("ari.devices.all", func(_ string, _ []byte, reply Reply) {

		dx, err := srv.upstream.DeviceState.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var ret []string
		for _, device := range dx {
			ret = append(ret, device.ID())
		}

		reply(ret, nil)
	})

	srv.subscribe("ari.devices.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.devices.data."):]
		data, err := srv.upstream.DeviceState.Data(name)
		reply(data, err)
	})

	srv.subscribe("ari.devices.update.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.devices.update."):]

		var state string
		if err := json.Unmarshal(data, &state); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.DeviceState.Update(name, state)
		reply(nil, err)
	})

	srv.subscribe("ari.devices.delete.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.devices.delete."):]
		err := srv.upstream.DeviceState.Delete(name)
		reply(nil, err)
	})

}
