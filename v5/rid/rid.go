// Package rid provides unique resource IDs
//
// Deprecated: Use github.com/CyCoreSystems/ari-rid instead.  All code here was reproduced there.
//
// This package is no longer used internally and will be removed in a future release.
package rid

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

const (
	// Bridge indicates the resource ID is of a bridge
	Bridge = "br"

	// Channel indicates the resource ID is of a channel
	Channel = "ch"

	// Playback indicates the resource ID is of a playback
	Playback = "pb"

	// Recording indicates the resource ID is for a recording
	Recording = "rc"

	// Snoop indicates the resource ID is for a snoop session
	Snoop = "sn"
)

// New returns a new generic resource ID
func New(kind string) string {
	id := strings.ToLower(ulid.MustNew(ulid.Now(), rand.Reader).String())

	if kind != "" {
		if len(kind) > 2 {
			kind = kind[:2]
		}

		id += "-" + kind
	}

	return id
}

// Timestamp returns the timestamp stored within the resource ID
func Timestamp(id string) (ts time.Time, err error) {
	idx := strings.Index(id, "-")
	if idx > 0 {
		id = id[:idx]
	}

	uid, err := ulid.Parse(id)
	if err != nil {
		return
	}

	ms := int64(uid.Time())
	ts = time.Unix(ms/1000, (ms%1000)*1000000)

	return
}
