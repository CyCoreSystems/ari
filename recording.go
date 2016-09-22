package ari

import "time"

// Recording is a namespace for the recording types
type Recording struct {
	Stored StoredRecording
	Live   LiveRecording
}

//NOTE: since record.Recorder needs RecordingOptions, we moved RecordingOptions into the ari package
// instead of placing it into record package. Since ari Handles will be implementing record.Recorder,
// this avoids a cyclical import as record already imports ari.

// RecordingOptions describes the set of options available when making a recording.
type RecordingOptions struct {
	// Format is the file format/encoding to which the recording should be stored.
	// This will usually be one of: slin, ulaw, alaw, wav, gsm.
	// If not specified, this will default to slin.
	Format string

	// MaxDuration is the maximum duration of the recording, after which the recording will
	// automatically stop.  If not set, there is no maximum.
	MaxDuration time.Duration

	// MaxSilence is the maximum duration of detected to be found before terminating the recording.
	MaxSilence time.Duration

	// Exists determines what should happen if the given recording already exists.
	// Valid values are: "fail", "overwrite", or "append".
	// If not specified, it will default to "fail"
	Exists string

	// Beep indicates whether a beep should be played to the recorded
	// party at the beginning of the recording.
	Beep bool

	// Terminate indicates whether the recording should be terminated on
	// receipt of a DTMF digit.
	// valid options are: "none", "any", "*", and "#"
	// If not specified, it will default to "none" (never terminate on DTMF).
	Terminate string
}
