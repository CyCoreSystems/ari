package ari

import (
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// Mute-related constants for use with the mute commands
// `MuteIn` mutes audio coming in to the channel from Asterisk
// `MuteOut` mutes audio coming out from the channel to Asterisk
// `MuteBoth` mutes audio in both directions
const (
	MuteIn   = "in"
	MuteOut  = "out"
	MuteBoth = "both"
)

// Channel describes a(n active) communication connection between Asterisk
// and an Endpoint
type Channel struct {
	Accountcode  string       `json:"accountcode"`
	Caller       CallerId     `json:"caller"`    // CallerId of the calling endpoint
	Connected    CallerId     `json:"connected"` // CallerId of the connected line
	Creationtime AsteriskDate `json:"creationtime"`
	Dialplan     DialplanCEP  `json:"dialplan"` // Current location in the dialplan
	Id           string       `json:"id"`       // Unique id for this channel (same as for AMI)
	Name         string       `json:"name"`     // Name of this channel (tech/name-id format)
	State        string       `json:"state"`    // State of the channel

	client *Client // Reference to the client which created or returned this channel
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

// String returns the stringified callerid
func (cid *CallerId) String() string {
	return cid.Name + "<" + cid.Number + ">"
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

// AttachClient attaches the provided ARI client to the
// channel
func (c *Channel) AttachClient(a *Client) {
	c.client = a
}

// GetClient returns the ARI client which created the channel
func (c *Channel) GetClient() *Client {
	return c.client
}

// GetID returns the ID of this channel
func (c *Channel) GetID() string {
	return c.Id
}

// IsReal indicates whether the channel is "real"
// (not a local channel)
func (c *Channel) IsReal() bool {
	return !strings.HasPrefix(c.Name, "Local/")
}

// Hangup hangs up the current channel
func (c *Channel) Hangup() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}
	return c.client.HangupChannel(c.Id, "normal")
}

// Busy hangs up the current channel with the "busy"
// cause code
func (c *Channel) Busy() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}
	return c.client.HangupChannel(c.Id, "busy")
}

// Congestion hangs up the current channel with the "congestion"
// cause code
func (c *Channel) Congestion() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}
	return c.client.HangupChannel(c.Id, "congestion")
}

// Continue causes the current channel to continue in
// the dialplan
func (c *Channel) Continue(context string, extension string, priority int64) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	req := ContinueChannelRequest{
		Context:   context,
		Extension: extension,
		Priority:  priority,
	}

	return c.client.ContinueChannel(c.Id, req)
}

// Answer answers the current channel
func (c *Channel) Answer() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.AnswerChannel(c.Id)
}

// IsAnswered checks the current state of the channel to
// see if it is "Up" (answered)
func (c *Channel) IsAnswered() (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("No client found in Channel")
	}
	updated, err := c.GetClient().GetChannel(c.GetID())
	if err != nil {
		return false, fmt.Errorf("Failed to get updated channel: %s", err.Error())
	}
	return strings.ToLower(updated.State) == "up", nil
}

// Ring indicates ringing to the channel
func (c *Channel) Ring() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.RingChannel(c.Id)
}

// StopRing stops ringing on the channel
func (c *Channel) StopRing() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.StopRinging(c.Id)
}

// SendDTMF sends DTMF to the channel
func (c *Channel) SendDTMF(dtmf string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	req := SendDTMFToChannelRequest{
		Dtmf: dtmf,
	}

	return c.client.SendDTMFToChannel(c.Id, req)
}

// Mute mutes the channel in the given direction
// (one of "in", "out", or "both")
func (c *Channel) Mute(dir string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.MuteChannel(c.Id, dir)
}

// Unmute stops muting of the channel in the given direction
// (one of "in", "out", or "both")
func (c *Channel) Unmute(dir string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.UnMuteChannel(c.Id, dir)
}

// Hold puts the channel on hold
func (c *Channel) Hold() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.HoldChannel(c.Id)
}

// Unhold retrieves the channel from hold
func (c *Channel) Unhold() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.StopHoldChannel(c.Id)
}

// MOH plays music on hold of the given class
// to the channel
func (c *Channel) MOH(mohClass string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.PlayMOHToChannel(c.Id, mohClass)
}

// StopMOH stops playing of music on hold to the channel
func (c *Channel) StopMOH() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.StopMohChannel(c.Id)
}

// Silence transmits silence (comfort noise) to the
// channel
func (c *Channel) Silence() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.PlaySilenceToChannel(c.Id)
}

// StopSilence stops transmission of silence to
// channel
func (c *Channel) StopSilence() error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.StopSilenceChannel(c.Id)
}

// Play initiates playback of the specified media uri
// to the channel, returning the Playback's Id
func (c *Channel) Play(mediaUri string) (string, error) {
	id := uuid.NewV1().String()
	err := c.PlayWithID(id, mediaUri)
	return id, err
}

// PlayWithID initiates playback of the specified media uri
// with the provided playbackID to the channel.
func (c *Channel) PlayWithID(id, mediaUri string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	var err error
	_, err = c.client.PlayToChannelById(c.Id, id, PlayMediaRequest{Media: mediaUri})
	return err
}

// Record starts recording the channel, returning
// the LiveRecording
func (c *Channel) Record(name string, opts *RecordingOptions) (*LiveRecording, error) {
	if opts == nil {
		opts = &RecordingOptions{}
	}
	if c.client == nil {
		return nil, fmt.Errorf("No client found in Channel")
	}

	return c.GetClient().RecordChannel(c.Id, opts.ToRequest(name))
}

// Get retrieves a channel variable from the channel
func (c *Channel) Get(name string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("No client found in Channel")
	}

	chanVar, err := c.client.GetChannelVariable(c.Id, name)
	if err != nil {
		// Asterisk (as of 13.9.1) returns an Internal Server Error
		// when requesting a PJSIP header which does not exist;
		// therefore, we treat 500 as simply an empty value.
		if rerr, ok := err.(RequestError); ok {
			if rerr.Code() == 500 {
				return "", nil
			}
		}
		return "", err
	}
	return chanVar.Value, nil
}

// Set sets a channel variable of the channel
func (c *Channel) Set(name string, value string) error {
	if c.client == nil {
		return fmt.Errorf("No client found in Channel")
	}

	return c.client.SetChannelVariable(c.Id, name, value)
}

// Snoop begins a snoop session and returns its id
// TODO: what is the channel being returned; do we need it?
func (c *Channel) Snoop() (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("No client found in Channel")
	}

	var err error
	id := uuid.NewV1().String()
	_, err = c.client.StartSnoopChannelById(c.Id, id, SnoopRequest{App: c.client.Options.Application})
	if err != nil {
		return "", fmt.Errorf("Failed to initiate snoop session: %s", err.Error())
	}
	return id, nil
}

//List all active channels in asterisk
//Equivalent to Get /channels
func (c *Client) ListChannels() ([]Channel, error) {
	var m []Channel
	err := c.Get("/channels", &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// NewOriginateRequest generates an originate request
// with a unique channel Id, destination equal to the
// current client, and an unlimited call timeout
func (c *Client) NewOriginateRequest(endpoint string) OriginateRequest {
	return OriginateRequest{
		Endpoint:  endpoint,
		Timeout:   -1,
		App:       c.Options.Application,
		ChannelId: uuid.NewV1().String(),
	}
}

// NewChannel is a shorthand for creating a new channel.
// It generates a unique Id, sets the destination to be
// the current application, and passes any variables
// through;  pass nil for vars if no variables are needed
func (c *Client) NewChannel(endpoint string, cid *CallerId, vars map[string]string) (Channel, error) {
	o := c.NewOriginateRequest(endpoint)
	if cid != nil {
		o.CallerId = cid.String()
	}
	if vars != nil {
		o.Variables = vars
	}
	return c.CreateChannelWithId(o.ChannelId, o)
}

// NewLocalChannel creates a new local channel and returns each side
func (c *Client) NewLocalChannel(endpoint string, cid *CallerId, vars map[string]string) (Channel, Channel, error) {
	o := c.NewOriginateRequest(endpoint)
	if cid != nil {
		o.CallerId = cid.String()
	}
	if vars != nil {
		o.Variables = vars
	}

	o.OtherChannelId = uuid.NewV1().String()

	ch, err := c.CreateChannelWithId(o.ChannelId, o)
	if err != nil {
		return ch, ch, err
	}

	ch2, err := c.GetChannel(o.OtherChannelId)
	if err != nil {
		return ch, ch2, err
	}

	return ch, ch2, nil
}

// CreateChannel originates a new call
func (c *Client) CreateChannel(req OriginateRequest) (Channel, error) {
	var m Channel

	err := c.Post("/channels", &m, &req)

	// Attach the client
	m.client = c

	return m, err
}

// CreateChannelWithId originates a new call with
// the given channel Id
func (c *Client) CreateChannelWithId(id string, req OriginateRequest) (Channel, error) {
	var m Channel

	if id == "" {
		return m, fmt.Errorf("No channel Id provided")
	}
	req.ChannelId = id

	err := c.Post("/channels/"+id, &m, &req)

	// Attach the client
	m.client = c

	return m, err
}

//Get a specific channel's details
//Equivalent to Get /channels/{channelId}
func (c *Client) GetChannel(channelId string) (Channel, error) {
	var m Channel
	err := c.Get("/channels/"+channelId, &m)
	if err != nil {
		return m, err
	}

	// Attach the client
	m.client = c

	return m, nil
}

//Exit application and continue execution in the dialplan
//Equivalent to Post /channels/{channelId}/continue
func (c *Client) ContinueChannel(channelId string, req ContinueChannelRequest) error {
	err := c.Post("/channels/"+channelId+"/continue", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Answer a channel
//Equivalent to Post /channels/{channelId}/answer
func (c *Client) AnswerChannel(channelId string) error {
	err := c.Post("/channels/"+channelId+"/answer", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Indicate ringing to a channel
//Equivalent to Post /channels/{channelId}/ring
func (c *Client) RingChannel(channelId string) error {
	err := c.Post("/channels/"+channelId+"/ring", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

//Send provided DTMF to a given channel
//Equivalent to Post /channels/{channelId}/dtmf
func (c *Client) SendDTMFToChannel(channelId string, req SendDTMFToChannelRequest) error {
	err := c.Post("/channels/"+channelId+"/dtmf", nil, &req)
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
	err = c.Post("/channels/"+channelId+"/mute", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Hold a channel
//Equivalent to Post /channels/{channelId}/hold
func (c *Client) HoldChannel(channelId string) error {
	err := c.Post("/channels/"+channelId+"/hold", nil, nil)
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

	err := c.Post("/channels/"+channelId+"/moh", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Play silence to a channel
//Equivalent to Post /channels/{channelId}/silence
func (c *Client) PlaySilenceToChannel(channelId string) error {
	err := c.Post("/channels/"+channelId+"/silence", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// PlayMedia is a wrapper to initiate a playback
// with the given Media URI, returning the playback id
func (c *Client) PlayMedia(channelId, mediaUri string) (string, error) {
	playbackId := uuid.NewV1().String()
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

	err := c.Post("/channels/"+channelId+"/play", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Specifiy media to playback to a channel
//Equivalent to Post /channels/{channelId}/play/{playbackId}
func (c *Client) PlayToChannelById(channelId string, playbackId string, req PlayMediaRequest) (Playback, error) {

	var m Playback

	err := c.Post("/channels/"+channelId+"/play/"+playbackId, &m, &req)
	return m, err
}

//Start a live recording
//Equivalent to Post /channels/{channelId}/record
func (c *Client) RecordChannel(channelId string, req *RecordRequest) (*LiveRecording, error) {
	var m LiveRecording

	m.client = c

	err := c.Post("/channels/"+channelId+"/record", &m, req)
	return &m, err
}

//Get the value of a channel variable or function
//Equivalent to Get /channels/{channelId}/variable
func (c *Client) GetChannelVariable(channelId string, variable string) (Variable, error) {
	var m Variable
	err := c.Get("/channels/"+channelId+"/variable?variable="+variable, &m)
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

	err := c.Post("/channels/"+channelId+"/variable", nil, &req)
	if err != nil {
		return err
	}
	return nil
}

//Start snooping
//Equivalent to Post /channels/{channelId}/snoop
func (c *Client) StartSnoopChannel(channelId string, req SnoopRequest) (Channel, error) {
	var m Channel

	err := c.Post("/channels/"+channelId+"/snoop", &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Start Snooping by specific ID
//Equivalent to Post /channels/{channelId}/snoop/{snoopId}
func (c *Client) StartSnoopChannelById(channelId string, snoopId string, req SnoopRequest) (Channel, error) {
	var m Channel

	err := c.Post("/channels/"+channelId+"/snoop/"+snoopId, &m, &req)
	if err != nil {
		return m, err
	}
	return m, nil
}

//Delete (i.e. hangup) a channel.
//Equivalent to DELETE /channels/{channelId}
func (c *Client) HangupChannel(channelId string, reason string) error {
	var req string
	if reason != "" {
		req = fmt.Sprintf("reason=%s", reason)
	}

	//send request
	return c.Delete("/channels/"+channelId, nil, req)
}

//Stop ringing indication on a channel if locally generated.
//Equivalent to DELETE /channels/{channelId}/ring
func (c *Client) StopRinging(channelId string) error {
	err := c.Delete("/channels/"+channelId+"/ring", nil, "")
	return err
}

//Unmute a channel
//Equivalent to DELETE /channels/{channelId}/mute
func (c *Client) UnMuteChannel(channelId string, direction string) error {

	err := c.CheckDirection(&direction)

	if err != nil {
		return err
	}

	var req string
	if direction != "" {
		req = fmt.Sprintf("direction=%s", direction)
	}

	err = c.Delete("/channels/"+channelId+"/mute", nil, req)
	return err
}

//Stop playing music on hold to a channel
//Equivalent to DELETE /channels/{channelId}/moh
func (c *Client) StopMohChannel(channelId string) error {
	err := c.Delete("/channels/"+channelId+"/moh", nil, "")
	return err
}

//Stop playing silence to a channel
//Equivalent to DELETE /channels/{channelId}/silence
func (c *Client) StopSilenceChannel(channelId string) error {
	err := c.Delete("/channels/"+channelId+"/silence", nil, "")
	return err
}

//Remove a channel from hold
//Equivalent to DELETE /channels/{channelId}/hold
func (c *Client) StopHoldChannel(channelId string) error {
	err := c.Delete("/channels/"+channelId+"/hold", nil, "")
	return err
}

//
//  Context-related items
//

// channelKey is the key type for contexts
type channelKey string

// NewChannelContext returns a context with the channel attached
func NewChannelContext(ctx context.Context, c *Channel) context.Context {
	return NewChannelContextWithKey(ctx, c, "_default")
}

// NewChannelContext returns a context with the channel attached
// as the given key
func NewChannelContextWithKey(ctx context.Context, c *Channel, name string) context.Context {
	return context.WithValue(ctx, channelKey(name), c)
}

// ChannelFromContext returns the default Channel stored in the context
func ChannelFromContext(ctx context.Context) (*Channel, bool) {
	return ChannelFromContextWithKey(ctx, "_default")
}

// ChannelFromContextWithKey returns the default Channel
// stored in the context
func ChannelFromContextWithKey(ctx context.Context, name string) (*Channel, bool) {
	c, ok := ctx.Value(channelKey(name)).(*Channel)
	return c, ok
}
