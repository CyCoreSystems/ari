package testutils

import "github.com/AVOXI/ari"

// A PlayerPair is the pair of results returned from a mock Play request
type PlayerPair struct {
	Handle *ari.PlaybackHandle
	Err    error
}

// Player is the test player that can be primed with sample data
type Player struct {
	Next    chan struct{}
	results []PlayerPair
}

// NewPlayer creates a new player
func NewPlayer() *Player {
	return &Player{
		Next: make(chan struct{}, 10),
	}
}

// Append appends the given Play results
func (p *Player) Append(h *ari.PlaybackHandle, err error) {
	p.results = append(p.results, PlayerPair{h, err})
}

// Play pops the top results and returns them, as well as triggering player.Next
func (p *Player) Play(mediaURI string) (h *ari.PlaybackHandle, err error) {
	h = p.results[0].Handle
	err = p.results[0].Err
	p.results = p.results[1:]

	p.Next <- struct{}{}

	return
}
