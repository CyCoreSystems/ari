package audio

import "github.com/CyCoreSystems/ari"

// A Player is an entity which can have an audio URI played
type Player interface {
	Play(string) (*ari.PlaybackHandle, error)
}
