package record

import "github.com/CyCoreSystems/ari"

// A Recording is a lifecycle managed audio recording
type Recording struct {
	Opts   *ari.RecordingOptions
	Handle *ari.LiveRecordingHandle
	Status Status
}
