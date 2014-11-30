package ari

import "code.google.com/p/go-uuid/uuid"

// Channel describes a(n active) communication connection between Asterisk
// and an Endpoint
type Channel struct {
	Accountcode  string       `json:"accountcode"`
	Caller       CallerId     `json:"caller"`    // CallerId of the calling endpoint
	Connected    CallerId     `json:"connected"` // CallerId of (TODO: what?)
	Creationtime AsteriskDate `json:"creationtime"`
	Dialplan     DialplanCEP  `json:"dialplan"` // Current location in the dialplan
	Id           string       `json:"id"`       // Unique id for this channel (same as for AMI)
	Name         string       `json:"name"`     // Name of this channel (tech/name-id format)
	State        string       `json:"state"`    // State of the channel
}

// OriginateRequest is the basic structure for all channel creation methods
type OriginateRequest struct {
	Endpoint string `json:"endpoint"`           // Endpoint to use (tech/resource notation)
	Timeout  int    `json:"timeout,omitempty"`  // Dial Timeout in seconds (-1 = no limit)
	CallerId string `json:"callerId,omitempty"` // CallerId to set for outgoing call

	// One set of:
	Context   string `json:"context,omitempty"` // Drop the channel into the dialplan
	Extension string `json:"extension,omitempty"`
	Priority  int64  `json:"priority,omitempty"`
	// OR
	App     string `json:"app,omitempty"`     // Associate channel to Stasis (Ari) application
	AppArgs string `json:"appArgs,omitempty"` // Arguments to the application
	// End OneSetOf

	// Channel Id declarations
	ChannelId      string `json:"channelId,omitempty"`      // Optionally assign channel id
	OtherChannelId string `json:"otherChannelId,omitempty"` // Optionally assign second channel's id (only for local channels)

	Variables map[string]string `json:"variables,omitempty"` // Channel variables to set
}

// CallerId describes the name and number which identifies the caller
// to other endpoints
type CallerId struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// DialplanCEP describes a location in the dialplan (context,extension,priority)
type DialplanCEP struct {
	Context  string `json:"context"`
	Exten    string `json:"exten"`
	Priority int64  `json:"priority"` //int64 derived from Java's 'long'
}

// Variable describes the value of a channel variable
//    NOTE: the variable name is not included, so it must be
//    tracked with the request
type Variable struct {
	Value string `json:"value"`
}

//priority is of type long = int64
//Request structure for ContinueChannel. All fields are optional.
//The three fields mirror the construction of a channel by Dialplan, as the function returns a channel to the dialplan.
type ContinueChannelRequest struct {
	Context   string `json:"context,omitempty"`
	Extension string `json:"extension,omitempty"`
	Priority  int64  `json:"priority,omitempty"`
}

//Structure for snooping a channel. Only App is required.
type SnoopRequest struct {
	Spy     string `json:"spy,omitempty"`     //Direction of audio to spy on, default is 'none'
	Whisper string `json:"whisper,omitempty"` //Direction of audio to whisper into, default is 'none'
	App     string `json:"app"`               //Application that the snooping channel is placed into
	AppArgs string `json:"appArgs,omitempty"` //The application arguments to pass to the Stasis application

	//Only necessary for 'StartSnoopChannel'
	SnoopId string `json:"snoopId,omitempty"` //Unique ID to assign to snooping channel
}

//SendDTMPFToChannel Request structure. All fields are optional.
type SendDTMFToChannelRequest struct {
	Dtmf     string `json:"dtmf,omitempty"`
	Before   int    `json:"before,omitempty"`
	Between  int    `json:"between,omitempty"`
	Duration int    `json:"duration,omitempty"`
	After    int    `json:"after,omitempty"`
}

//List all active channels in asterisk
//Equivalent to Get /channels
func (c *Client) ListChannels() ([]Channel, error) {
	var m []Channel
	err := c.AriGet("/channels", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// CreateChannel originates a new call
func (c *Client) CreateChannel(req OriginateRequest) (Channel, error) {
	var m Channel

	err := c.AriPost("/channels", &m, &req)
	return m, err
}

//Get a specific channel's details
//Equivalent to Get /channels/{channelId}
func (c *Client) GetChannel(channelId string) (Channel, error) {
	var m Channel
	err := c.AriGet("/channels/"+channelId, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Exit application and continue execution in the dialplan
//Equivalent to Post /channels/{channelId}/continue
func (c *Client) ContinueChannel(channelId string, req ContinueChannelRequest) error {
	err := c.AriPost("/channels/"+channelId+"/continue", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Answer a channel
//Equivalent to Post /channels/{channelId}/answer
func (c *Client) AnswerChannel(channelId string) error {
	err := c.AriPost("/channels/"+channelId+"/answer", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Indicate ringing to a channel
//Equivalent to Post /channels/{channelId}/ring
func (c *Client) ChannelRing(channelId string) error {
	err := c.AriPost("/channels/"+channelId+"/ring", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Send provided DTMF to a given channel
//Equivalent to Post /channels/{channelId}/dtmf
func (c *Client) SendDTMFToChannel(channelId string, req SendDTMFToChannelRequest) error {
	err := c.AriPost("/channels/"+channelId+"/dtmf", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Mute a channel
//Equivalent to Post /channels/{channelId}/mute
//Viable options are "in," "out," or "both"
func (c *Client) MuteChannel(channelId string, direction string) error {
	err := c.CheckDirection(&direction)

	if err != nil {
		return err
	}

	type request struct {
		Direction string `json:"direction,omitempty"`
	}

	req := request{direction}
	err = c.AriPost("/channels/"+channelId+"/mute", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Hold a channel
//Equivalent to Post /channels/{channelId}/hold
func (c *Client) HoldChannel(channelId string) error {
	err := c.AriPost("/channels/"+channelId+"/hold", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Play music on hold to a channel
//Equivalent to Post /channels/{channelId}/moh
func (c *Client) PlayMOHToChannel(channelId string, mohClass string) error {
	type request struct {
		MohClass string `json:"mohClass,omitempty"`
	}
	req := request{mohClass}

	err := c.AriPost("/channels/"+channelId+"/moh", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Play silence to a channel
//Equivalent to Post /channels/{channelId}/silence
func (c *Client) PlaySilenceToChannel(channelId string) error {
	err := c.AriPost("/channels/"+channelId+"/silence", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// PlayMedia is a wrapper to initiate a playback
// with the given Media URI, returning the playback id
func (c *Client) PlayMedia(channelId, mediaUri string) (string, error) {
	playbackId := uuid.New()
	req := PlayMediaRequest{
		Media: mediaUri,
	}
	_, err := c.PlayToChannelById(channelId, playbackId, req)
	return playbackId, err
}

//Start playback of media to a channel
//Equivalent to Post /channels/{channelId}/play
func (c *Client) PlayToChannel(channelId string, req PlayMediaRequest) (Playback, error) {
	var m Playback

	err := c.AriPost("/channels/"+channelId+"/play", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Specifiy media to playback to a channel
//Equivalent to Post /channels/{channelId}/play/{playbackId}
func (c *Client) PlayToChannelById(channelId string, playbackId string, req PlayMediaRequest) (Playback, error) {

	var m Playback

	err := c.AriPost("/channels/"+channelId+"/play/"+playbackId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Start a live recording
//Equivalent to Post /channels/{channelId}/record
func (c *Client) RecordChannel(channelId string, req RecordRequest) (LiveRecording, error) {
	var m LiveRecording

	err := c.AriPost("/channels/"+channelId+"/record", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Get the value of a channel variable or function
//Equivalent to Get /channels/{channelId}/variable
func (c *Client) GetChannelVariable(channelId string, variable string) (Variable, error) {
	var m Variable
	err := c.AriGet("/channels/"+channelId+"/variable?variable="+variable, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Set the value of a variable
//Equivalent to Post /channels/{channelId}/variable
func (c *Client) SetChannelVariable(channelId string, variable string, value string) error {
	type request struct {
		Variable string `json:"variable"`
		Value    string `json:"value,omitempty"`
	}

	req := request{variable, value}

	err := c.AriPost("/channels/"+channelId+"/variable", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Start snooping
//Equivalent to Post /channels/{channelId}/snoop
func (c *Client) StartSnoopChannel(channelId string, req SnoopRequest) (Channel, error) {
	var m Channel

	err := c.AriPost("/channels/"+channelId+"/snoop", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Start Snooping by specific ID
//Equivalent to Post /channels/{channelId}/snoop/{snoopId}
func (c *Client) StartSnoopChannelById(channelId string, snoopId string, req SnoopRequest) (Channel, error) {
	var m Channel

	err := c.AriPost("/channels/"+channelId+"/snoop/"+snoopId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Delete (i.e. hangup) a channel.
//Equivalent to DELETE /channels/{channelId}
func (c *Client) HangupChannel(channelId string, reason string) error {
	//Request structure for hanging up a channel. Reason is not required.
	type request struct {
		Reason string `json:"reason,omitempty"`
	}

	req := request{reason}

	//send request
	err := c.AriDelete("/channels/"+channelId, nil, &req)
	return err
}

//Stop ringing indication on a channel if locally generated.
//Equivalent to DELETE /channels/{channelId}/ring
func (c *Client) StopRinging(channelId string) error {
	err := c.AriDelete("/channels/"+channelId+"/ring", nil, nil)
	return err
}

//Unmute a channel
//Equivalent to DELETE /channels/{channelId}/mute
func (c *Client) UnMuteChannel(channelId string, direction string) error {

	err := c.CheckDirection(&direction)

	if err != nil {
		return err
	}

	type request struct {
		Direction string `json:"direction,omitempty"`
	}

	req := request{direction}
	err = c.AriDelete("/channels/"+channelId+"/mute", nil, &req)
	return err
}

//Stop playing music on hold to a channel
//Equivalent to DELETE /channels/{channelId}/moh
func (c *Client) StopMohChannel(channelId string) error {
	err := c.AriDelete("/channels/"+channelId+"/moh", nil, nil)
	return err
}

//Stop playing silence to a channel
//Equivalent to DELETE /channels/{channelId}/silence
func (c *Client) StopSilenceChannel(channelId string) error {
	err := c.AriDelete("/channels/"+channelId+"/silence", nil, nil)
	return err
}

//Remove a channel from hold
//Equivalent to DELETE /channels/{channelId}/hold
func (c *Client) StopHoldChannel(channelId string) error {
	err := c.AriDelete("/channels/"+channelId+"/hold", nil, nil)
	return err
}
