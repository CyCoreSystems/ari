package ari

// Direction describes an audio direction, as used by Mute, Snoop, and possibly others.  Valid values are "in", "out", and "both".
type Direction string

const (
	// DirectionIn indicates the direction flowing from the channel into Asterisk
	DirectionIn Direction = "in"

	// DirectionOut indicates the direction flowing from Asterisk to the channel
	DirectionOut Direction = "out"

	// DirectionBoth indicates both the directions flowing both inward to Asterisk and outward from Asterisk.
	DirectionBoth Direction = "both"
)
