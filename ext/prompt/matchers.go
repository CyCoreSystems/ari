package prompt

import "strings"

// MatchAny is a MatchFunc which returns Complete if the pattern
// contains any characters.
func MatchAny(pat string) (string, Status) {
	if len(pat) > 0 {
		return pat, Complete
	}
	return pat, Incomplete
}

// MatchHash is a MatchFunc which returns Complete if the pattern
// contains a hash (#) character and Incomplete, otherwise.
func MatchHash(pat string) (string, Status) {
	if strings.Contains(pat, "#") {
		return strings.Split(pat, "#")[0], Complete
	}
	return pat, Incomplete
}

// MatchSetFunc returns a MatchFunc which will match any digit in the provided set of digits
func MatchSetFunc(set string) func(string) (string, Status) {
	return func(pat string) (string, Status) {
		if len(pat) < 1 {
			return pat, Incomplete
		}

		for _, c := range set {
			if strings.Contains(pat, string(c)) {
				return pat, Complete
			}
		}
		return pat, Invalid
	}
}

// MatchTerminatorFunc is a MatchFunc which returns Complete if the pattern
// contains the character and Incomplete, otherwise.
func MatchTerminatorFunc(terminator string) func(string) (string, Status) {
	return func(pat string) (string, Status) {
		if strings.Contains(pat, terminator) {
			return strings.Split(pat, terminator)[0], Complete
		}
		return pat, Incomplete
	}
}

// MatchLenFunc returns a MatchFunc which returns
// Complete if the given number of digits are received
// and Incomplete otherwise.
func MatchLenFunc(length int) func(string) (string, Status) {
	return func(pat string) (string, Status) {
		if len(pat) >= length {
			return pat, Complete
		}
		return pat, Incomplete
	}
}

// MatchLenOrTerminatorFunc returns a MatchFunc which returns
// Complete if the given number of digits are received
// or the given terminal string is received.  Otherwise,
// it returns Incomplete.
func MatchLenOrTerminatorFunc(length int, terminator string) func(string) (string, Status) {
	return func(pat string) (string, Status) {
		if len(pat) >= length {
			return pat, Complete
		}
		if strings.HasSuffix(pat, terminator) {
			return strings.TrimSuffix(pat, terminator), Complete
		}
		return pat, Incomplete
	}
}
