package ari

// OriginateRequest is the basic structure for all channel creation methods
type OriginateRequest struct {
	Endpoint string `json:"endpoint"`           // Endpoint to use (tech/resource notation)
	Timeout  int    `json:"timeout,omitempty"`  // Dial Timeout in seconds (-1 = no limit)
	CallerID string `json:"callerId,omitempty"` // CallerID to set for outgoing call

	// One set of:
	Context   string `json:"context,omitempty"` // Drop the channel into the dialplan
	Extension string `json:"extension,omitempty"`
	Priority  int64  `json:"priority,omitempty"`
	// OR
	App     string `json:"app,omitempty"`     // Associate channel to Stasis (Ari) application
	AppArgs string `json:"appArgs,omitempty"` // Arguments to the application
	// End OneSetOf

	// Channel ID declarations
	ChannelID      string `json:"channelId,omitempty"`      // Optionally assign channel id
	OtherChannelID string `json:"otherChannelId,omitempty"` // Optionally assign second channel's id (only for local channels)

	Variables map[string]string `json:"variables,omitempty"` // Channel variables to set
}
