package ari

//Listed below are the structures ubiquitously used for media playback or recording.

type PlayMediaRequest struct {
	Media    string `json:"media"`
	Lang     string `json:"lang,omitempty"`
	Offsetms int    `json:"offsetms,omitempty"`
	Skipms   int    `json:"skipms,omitempty"`

	//PlaybackId unnecessary if specifically naming it in query path
	PlaybackId string `json:"playbackId,omitempty"`
}

//Identical for recording both channels and bridges.
//options for IfExists are 'overwrite','fail' (the default), and 'append.' This represents what to do if the name of the recording already exists.

type RecordRequest struct {
	Name               string `json:"name"`
	Format             string `json:"format"`
	MaxDurationSeconds int    `json:"maxDurationSeconds,omitempty"`
	MaxSilenceSeconds  int    `json:"maxSilenceSeconds,omitempty"`
	IfExists           string `json:"ifExists,omitempty"`
	Beep               bool   `json:"beep,omitempty"`
	TerminateOn        string `json:"terminateOn,omitempty"`
}
