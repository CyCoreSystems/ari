package prompt

import (
	"strings"

	"github.com/CyCoreSystems/ari/ext"
)

// MatchAny is a MatchFunc which returns Complete if the pattern
// contains any characters.
func MatchAny(pat string) (string, ext.Status) {
	if len(pat) > 0 {
		return pat, ext.Complete
	}
	return pat, ext.Incomplete
}

// MatchHash is a MatchFunc which returns Complete if the pattern
// contains a hash (#) character and Incomplete, otherwise.
func MatchHash(pat string) (string, ext.Status) {
	if strings.Contains(pat, "#") {
		return strings.Split(pat, "#")[0], ext.Complete
	}
	return pat, ext.Incomplete
}

// MatchTerminatorFunc is a MatchFunc which returns Complete if the pattern
// contains the character and Incomplete, otherwise.
func MatchTerminatorFunc(terminator string) func(string) (string, ext.Status) {
	return func(pat string) (string, ext.Status) {
		if strings.Contains(pat, terminator) {
			return strings.Split(pat, terminator)[0], ext.Complete
		}
		return pat, ext.Incomplete
	}
}

// MatchLenFunc returns a MatchFunc which returns
// Complete if the given number of digits are received
// and Incomplete otherwise.
func MatchLenFunc(length int) func(string) (string, ext.Status) {
	return func(pat string) (string, ext.Status) {
		if len(pat) >= length {
			return pat, ext.Complete
		}
		return pat, ext.Incomplete
	}
}

// MatchLenOrTerminatorFunc returns a MatchFunc which returns
// Complete if the given number of digits are received
// or the given terminal string is received.  Otherwise,
// it returns Incomplete.
func MatchLenOrTerminatorFunc(length int, terminator string) func(string) (string, ext.Status) {
	return func(pat string) (string, ext.Status) {
		if len(pat) >= length {
			return pat, ext.Complete
		}
		if strings.HasSuffix(pat, terminator) {
			return strings.TrimSuffix(pat, terminator), ext.Complete
		}
		return pat, ext.Incomplete
	}
}
