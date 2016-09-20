package natsgw

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/nc"
)

func (srv *Server) channel() {

	srv.subscribe("ari.channels.all", func(_ string, _ []byte, reply Reply) {
		cx, err := srv.upstream.Channel.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var channels []string
		for _, channel := range cx {
			channels = append(channels, channel.ID())
		}

		reply(channels, nil)
	})

	srv.subscribe("ari.channels.create", func(subj string, data []byte, reply Reply) {
		var req ari.OriginateRequest

		if err := json.Unmarshal(data, &req); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		handle, err := srv.upstream.Channel.Create(req)

		if err != nil {
			reply(nil, err)
			return
		}

		reply(handle.ID(), nil)
	})

	srv.subscribe("ari.channels.data.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.data."):]
		d, err := srv.upstream.Channel.Data(name)
		reply(&d, err)
	})

	srv.subscribe("ari.channels.answer.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.answer."):]
		srv.log.Debug("answering channel", "subj", subj)
		err := srv.upstream.Channel.Answer(name)
		srv.log.Debug("answered channel", "subj", subj, "name", name, "error", err)

		reply(nil, err)
	})

	srv.subscribe("ari.channels.hangup.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.hangup."):]

		var reason string
		if err := json.Unmarshal(data, &reason); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Channel.Hangup(name, reason)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.ring.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.ring."):]
		err := srv.upstream.Channel.Ring(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.stopring.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.stopring."):]
		err := srv.upstream.Channel.StopRing(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.hold.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.hold."):]
		err := srv.upstream.Channel.Hold(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.stophold.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.stophold."):]
		err := srv.upstream.Channel.StopHold(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.mute.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.mute."):]

		var dir string
		if err := json.Unmarshal(data, &dir); err != nil {
			reply(nil, &decodingError{subj, err})
		}

		err := srv.upstream.Channel.Mute(name, dir)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.unmute.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.unmute."):]

		var dir string
		if err := json.Unmarshal(data, &dir); err != nil {
			reply(nil, &decodingError{subj, err})
		}

		err := srv.upstream.Channel.Unmute(name, dir)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.silence.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.silence."):]
		err := srv.upstream.Channel.Silence(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.stopsilence.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.stopsilence."):]
		err := srv.upstream.Channel.StopSilence(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.moh.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.moh."):]

		var music string
		if err := json.Unmarshal(data, &music); err != nil {
			reply(nil, &decodingError{subj, err})
		}

		err := srv.upstream.Channel.MOH(name, music)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.stopmoh.>", func(subj string, _ []byte, reply Reply) {
		name := subj[len("ari.channels.stopmoh."):]
		err := srv.upstream.Channel.StopMOH(name)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.play.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.play."):]

		var pr nc.PlayRequest
		if err := json.Unmarshal(data, &pr); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		_, err := srv.upstream.Channel.Play(name, pr.PlaybackID, pr.MediaURI)
		reply(nil, err)
	})

	srv.subscribe("ari.channels.dtmf.>", func(subj string, data []byte, reply Reply) {
		name := subj[len("ari.channels.dtmf."):]

		type request struct {
			Dtmf string           `json:"dtmf,omitempty"`
			Opts *ari.DTMFOptions `json:"options,omitempty"`
		}

		var req request
		if err := json.Unmarshal(data, &req); err != nil {
			reply(nil, &decodingError{subj, err})
			return
		}

		err := srv.upstream.Channel.SendDTMF(name, req.Dtmf, req.Opts)
		reply(nil, err)
	})

}
