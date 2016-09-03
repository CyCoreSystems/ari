// Package audiouri provides conversions for common sounds to asterisk-supported audio URIs
package audiouri

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// SupportedPlaybackPrefixes is a list of valid prefixes
// for media URIs.
var SupportedPlaybackPrefixes = []string{
	"sound", "recording", "number", "digits", "characters", "tone",
}

func init() {
	sort.Strings(SupportedPlaybackPrefixes)
}

// WaitURI returns the set of media URIs for the
// given period of silence.
func WaitURI(t time.Duration) []string {
	q := []string{}
	for i := time.Duration(0); i <= t; i += time.Second {
		q = append(q, "sound:silence/1")
	}
	return q
}

// NumberURI returns the media URI to play
// the given number.
func NumberURI(number int) string {
	return fmt.Sprintf("number:%d", number)
}

// DigitsURI returns the set of media URIs to
// play the given digits.
func DigitsURI(digits string, hash string) []string {
	if strings.Contains(digits, "#") && hash != "" {
		hash = "sound:char/" + hash

		// Split the digits by each hash
		pieces := strings.Split(digits, "#")

		// Get the strings for each digit substring
		newDigits := []string{}
		for i, p := range pieces {
			newDigits = append(newDigits, DigitsURI(p, hash)...)
			if len(pieces) > i+1 {
				// If this wasn't the last piece, add a hash
				newDigits = append(newDigits, hash)
			}
		}
		return newDigits
	}

	// Handle '*' to 'star' conversion
	if strings.Contains(digits, "*") {
		star := "sound:char/star"

		// Split the digits by each *
		pieces := strings.Split(digits, "*")

		// Get the strings for each digit substring
		newDigits := []string{}
		for i, p := range pieces {
			newDigits = append(newDigits, DigitsURI(p, hash)...)
			if len(pieces) > i+1 {
				// If this wasn't the last piece, add a star
				newDigits = append(newDigits, star)
			}
		}
		return newDigits
	}

	// Otherwise, we can simply use the normal "digits:" URI
	if len(digits) < 1 {
		return []string{}
	}
	return []string{"digits:" + digits}
}

// DateTimeURI returns the set of media URIs for playing
// the current date and time.
func DateTimeURI(t time.Time) (ret []string) {
	ret = []string{}

	ret = append(ret,
		fmt.Sprintf("sound:digits/day-%d", t.Weekday()),
		fmt.Sprintf("sound:digits/mon-%d", t.Month()-1),
		NumberURI(t.Day()),
	)

	// Convert to 12-hour time
	pm := false
	hour := t.Hour()
	switch {
	case hour == 0:
		hour = 12
	case hour == 12:
		pm = true
	case hour > 12:
		hour -= 12
		pm = true
	default:
	}
	ret = append(ret, NumberURI(hour))

	// Humanize the minutes
	minute := t.Minute()
	switch {
	case minute == 0:
		ret = append(ret, "sound:digits/oclock")
	case minute < 10:
		ret = append(ret,
			"sound:digits/oh",
			NumberURI(minute),
		)
	default:
		ret = append(ret, NumberURI(minute))
	}

	// Add am/pm suffix
	if pm {
		ret = append(ret, "sound:digits/p-m")
	} else {
		ret = append(ret, "sound:digits/a-m")
	}

	// Add the year
	ret = append(ret, NumberURI(t.Year()))

	return
}

// DurationURI returns the set of media URIs for playing
// the given duration, in human terms (days, hours, minutes, seconds).
// If any of these terms are zero, they will not be spoken.
func DurationURI(dur time.Duration) (ret []string) {
	days := int(dur.Hours() / 24)
	hours := int(dur.Hours()) % 24
	minutes := int(dur.Minutes()) % 60
	seconds := int(dur.Seconds()) % 60

	ret = []string{}

	if days > 0 {
		ret = append(ret, NumberURI(days))
		if days > 1 {
			ret = append(ret, "sound:time/days")
		} else {
			ret = append(ret, "sound:time/day")
		}
	}

	if hours > 0 {
		ret = append(ret, NumberURI(hours))
		if hours > 1 {
			ret = append(ret, "sound:time/hours")
		} else {
			ret = append(ret, "sound:time/hour")
		}
	}

	if minutes > 0 {
		ret = append(ret, NumberURI(minutes))
		if minutes > 1 {
			ret = append(ret, "sound:time/minutes")
		} else {
			ret = append(ret, "sound:time/minute")
		}
	}

	if seconds > 0 {
		ret = append(ret, NumberURI(seconds))
		if seconds > 1 {
			ret = append(ret, "sound:time/seconds")
		} else {
			ret = append(ret, "sound:time/second")
		}
	}

	return
}

// RecordingURI returns the media URI for playing an
// Asterisk StoredRecording
func RecordingURI(name string) string {
	return "recording:" + name
}

// ToneURI returns the media URI for playing the
// given tone, which may be a system-defined indication
// or an explicit tone pattern construction.
func ToneURI(name string) string {
	return "tone:" + name
}

// checkAudioURI checks if the audio URI is formatted properly
func checkAudioURI(uri string) error {
	l := strings.Split(uri, ":")
	if len(l) != 2 {
		return fmt.Errorf("Audio URI %s is not formatted properly", uri)
	}

	if sort.SearchStrings(SupportedPlaybackPrefixes, l[0]) == len(SupportedPlaybackPrefixes) {
		return fmt.Errorf("Audio URI prefix %s not supported", l[0])
	}

	return nil
}
