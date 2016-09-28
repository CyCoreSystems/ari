package natsgw

import (
	"encoding/json"
	"time"

	"github.com/CyCoreSystems/ari"
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

	srv.subscribe("ari.bridges.data.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.bridges.data."):]
		bd, err := srv.upstream.Bridge.Data(name)
		reply(&bd, err)
		return
	})

	srv.subscribe("ari.bridges.create", func(subj string, data []byte, reply Reply) {

		var req nc.CreateBridgeRequest
		if err := json.Unmarshal(data, &req); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		bh, err := srv.upstream.Bridge.Create(req.ID, req.Type, req.Name)

		if err != nil {
			reply(nil, err)
			return
		}

		reply(bh.ID(), err)
	})

	srv.subscribe("ari.bridges.addChannel.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.addChannel."):]

		var channelID string
		if err := json.Unmarshal(data, &channelID); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Bridge.AddChannel(name, channelID)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.removeChannel.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.removeChannel."):]

		var channelID string
		if err := json.Unmarshal(data, &channelID); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Bridge.RemoveChannel(name, channelID)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.delete.*", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.bridges.delete."):]
		err := srv.upstream.Bridge.Delete(name)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.play.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.play."):]

		var pr nc.PlayRequest
		if err := json.Unmarshal(data, &pr); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		_, err := srv.upstream.Bridge.Play(name, pr.PlaybackID, pr.MediaURI)
		reply(nil, err)
	})

	srv.subscribe("ari.bridges.record.*", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.bridges.record."):]

		var rr nc.RecordRequest
		if err := json.Unmarshal(data, &rr); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		var opts ari.RecordingOptions

		opts.Format = rr.Format
		opts.MaxDuration = time.Duration(rr.MaxDuration) * time.Second
		opts.MaxSilence = time.Duration(rr.MaxSilence) * time.Second
		opts.Exists = rr.IfExists
		opts.Beep = rr.Beep
		opts.Terminate = rr.TerminateOn

		_, err := srv.upstream.Bridge.Record(name, rr.Name, &opts)
		reply(nil, err)
	})
}
