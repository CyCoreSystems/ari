package natsgw

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari/client/nc"
)

func (srv *Server) bridge() {

	srv.subscribe("ari.bridges.all", func(_ string, _ []byte, reply Reply) {

		bx, err := srv.upstream.Bridge.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var bridges []string
		for _, bridge := range bx {
			bridges = append(bridges, bridge.ID())
		}

		reply(bridges, nil)
	})

	srv.subscribe("ari.bridges.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.bridges.data."):]
		bd, err := srv.upstream.Bridge.Data(name)
		reply(&bd, err)
		return
	})

	srv.subscribe("ari.bridges.addChannel.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.addChannel."):]

		var channelID string
		if err := json.Unmarshal(data, &channelID); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Bridge.AddChannel(name, channelID)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.removeChannel.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.removeChannel."):]

		var channelID string
		if err := json.Unmarshal(data, &channelID); err != nil {
			reply(nil, err)
			return
		}

		err := srv.upstream.Bridge.RemoveChannel(name, channelID)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.delete.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.bridges.delete."):]
		err := srv.upstream.Bridge.Delete(name)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.play.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.play."):]

		var pr nc.PlayRequest
		if err := json.Unmarshal(data, &pr); err != nil {
			reply(nil, err)
			return
		}

		_, err := srv.upstream.Bridge.Play(name, pr.PlaybackID, pr.MediaURI)
		reply(nil, err)
	})

}
