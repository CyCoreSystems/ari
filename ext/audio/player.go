package audio

import "github.com/CyCoreSystems/ari"

// A Player is an entity which can have an audio URI played and can have event subscriptions
type Player interface {
	ari.Subscriber

	// Play plays the audio using the given playback ID and media URI
	Play(string, string) (*ari.PlaybackHandle, error)
}
