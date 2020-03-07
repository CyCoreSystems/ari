package ari

// OriginateRequest defines the parameters for the creation of a new Asterisk channel
type OriginateRequest struct {

	// Endpoint is the name of the Asterisk resource to be used to create the
	// channel.  The format is tech/resource.
	//
	// Examples:
	//
	//   - PJSIP/george
	//
	//   - Local/party@mycontext
	//
	//   - DAHDI/8005558282
	Endpoint string `json:"endpoint"`

	// Timeout specifies the number of seconds to wait for the channel to be
	// answered before giving up.  Note that this is REQUIRED and the default is
	// to timeout immediately.  Use a negative value to specify no timeout, but
	// be aware that this could result in an unlimited call, which could result
	// in a very unfriendly bill.
	Timeout int `json:"timeout,omitempty"`

	// CallerID specifies the Caller ID (name and number) to be set on the
	// newly-created channel.  This is optional but recommended.  The format is
	// `"Name" <number>`, but most every component is optional.
	//
	// Examples:
	//
	//   - "Jane" <100>
	//
	//   - <102>
	//
	//   - 8005558282
	//
	CallerID string `json:"callerId,omitempty"`

	// CEP (Context/Extension/Priority) is the location in the Asterisk dialplan
	// into which the newly created channel should be dropped.  All of these are
	// required if the CEP is used.  Exactly one of CEP or App/AppArgs must be
	// specified.
	Context   string `json:"context,omitempty"`
	Extension string `json:"extension,omitempty"`
	Priority  int64  `json:"priority,omitempty"`

	// The Label is the string form of Priority, if there is such a label in the
	// dialplan.  Like CEP, Label may not be used if an ARI App is specified.
	// If both Label and Priority are specified, Label will take priority.
	Label string `json:"label,omitempty"`

	// App specifies the ARI application and its arguments into which
	// the newly-created channel should be placed.  Exactly one of CEP or
	// App/AppArgs is required.
	App string `json:"app,omitempty"`

	// AppArgs defines the arguments to supply to the ARI application, if one is
	// defined.  It is optional but only applicable for Originations which
	// specify an ARI App.
	AppArgs string `json:"appArgs,omitempty"`

	// Formats describes the (comma-delimited) set of codecs which should be
	// allowed for the created channel.  This is an optional parameter, and if
	// an Originator is specified, this should be left blank so that Asterisk
	// derives the codecs from that Originator channel instead.
	//
	// Ex. "ulaw,slin16".
	//
	// The list of valid codecs can be found with Asterisk command "core show codecs".
	Formats string `json:"formats,omitempty"`

	// ChannelID specifies the unique ID to be used for the channel to be
	// created.  It is optional, and if not specified, a time-based UUID will be
	// generated.
	ChannelID string `json:"channelId,omitempty"` // Optionally assign channel id

	// OtherChannelID specifies the unique ID of the second channel to be
	// created.  This is only valid for the creation of Local channels, which
	// are always generated in pairs.  It is optional, and if not specified, a
	// time-based UUID will be generated (again, only if the Origination is of a
	// Local channel).
	OtherChannelID string `json:"otherChannelId,omitempty"`

	// Originator is the channel for whom this Originate request is being made, if there is one.
	// It is used by Asterisk to set the right codecs (and possibly other parameters) such that
	// when the new channel is bridged to the Originator channel, there should be no transcoding.
	// This is a purely optional (but helpful, where applicable) field.
	Originator string `json:"originator,omitempty"`

	// Variables describes the set of channel variables to apply to the new channel.  It is optional.
	Variables map[string]string `json:"variables,omitempty"`
}
