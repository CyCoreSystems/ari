package nc

import "github.com/CyCoreSystems/ari"

type natsSound struct {
	conn *Conn
}

func (s *natsSound) List(filters map[string]string) (sx []*ari.SoundHandle, err error) {

	if filters == nil {
		filters = make(map[string]string)
	}

	var sounds []string
	err = s.conn.readRequest("ari.sounds.all", &filters, &sounds)
	for _, sh := range sounds {
		sx = append(sx, s.Get(sh))
	}
	return
}

func (s *natsSound) Get(name string) *ari.SoundHandle {
	return ari.NewSoundHandle(name, s)
}

func (s *natsSound) Data(name string) (sd ari.SoundData, err error) {
	err = s.conn.readRequest("ari.sounds.data."+name, nil, &sd)
	return
}
