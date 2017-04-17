package audio

import (
	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/ext"
)

// Queue executes the series of media URIs in sequence to the given player.
func Queue(ctx context.Context, p ari.Player, mediaURIs ...string) (st ext.Status, err error) {
	for _, media := range mediaURIs {
		st, err = Play(ctx, p, media)
		if err != nil {
			return
		}
		if st != ext.Complete {
			return
		}
	}

	return
}
