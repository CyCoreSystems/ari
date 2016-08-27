package ari

// ChannelData is the data for a specific channel
type ChannelData struct {
	ID           string      `json:"id"`    // Unique id for this channel (same as for AMI)
	Name         string      `json:"name"`  // Name of this channel (tech/name-id format)
	State        string      `json:"state"` // State of the channel
	Accountcode  string      `json:"accountcode"`
	Caller       CallerID    `json:"caller"`    // CallerId of the calling endpoint
	Connected    CallerID    `json:"connected"` // CallerId of the connected line
	Creationtime DateTime    `json:"creationtime"`
	Dialplan     DialplanCEP `json:"dialplan"` // Current location in the dialplan
}
