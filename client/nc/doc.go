package nc

// Documentation is the list of nats endpoints, their request and response types, and their descriptions.
var Documentation = []struct {
	Endpoint    string
	Request     string
	Response    string
	Description string
}{
	{"ari.applications.all", "ignored", "[]string", "Get all applications"},
	{"ari.applications.data.>", "ignored", "ari.ApplicationData", "Get the data for the given applications"},
	{"ari.applications.subscribe.>", "string", "ignored", "Subscribe the app to the event source"},
	{"ari.applications.unsubscribe.>", "string", "ignored", "Unsubscribe the app from the event source"},

	{"ari.asterisk.reload.>", "ignored", "ignored", "Reload an Asterisk module"},
	{"ari.asterisk.info", "ignored", "*ari.AsteriskInfo", "Get information about the Asterisk server"},

	{"ari.asterisk.variables.get.>", "ignored", "string", "Get the global asterisk variable"},
	{"ari.asterisk.variables.set.>", "string", "ignored", "Set the global asterisk variable"},

	{"ari.bridges.all", "ignored", "[]string", "Get all the bridges"},
	{"ari.bridges.data.>", "ignored", "ari.BridgeData", "Get the bridge data for the bridge"},
	{"ari.bridges.addChannel.>", "string", "ignored", "Add the channel to the bridge"},
	{"ari.bridges.removeChannel.>", "string", "ignored", "Remove the channel from the bridge"},
	{"ari.bridges.play.>", "nc.PlayRequest", "ignored", "Play the media URI to the bridge"},

	{"ari.channels.all", "ignored", "[]string", "List all the channels"},
	{"ari.channels.create", "ari.OriginateRequest", "string", "Create a new channel, returning the channel ID"},
	{"ari.channels.data.>", "ignored", "ari.ChannelData", "Get the channel data"},
	{"ari.channels.hangup.>", "string", "ignored", "Hangup the channel, using the passed reason"},
	{"ari.channels.ring.>", "ignored", "ignored", "Ring the channel"},
	{"ari.channels.stopstring.>", "ignored", "ignored", "Stop ringing the channel"},
	{"ari.channels.hold.>", "ignored", "ignored", "Put the channel on hold"},
	{"ari.channels.stophold.>", "ignored", "ignored", "Stop the channel on hold"},
	{"ari.channels.mute.>", "string", "ignored", "Mute the channel in the givin direction"},
	{"ari.channels.unmute.>", "string", "ignored", "Unmute the channel in the givin direction"},
	{"ari.channels.silence.>", "ignored", "ignored", "Play silence on the channel"},
	{"ari.channels.stopsilence.>", "ignored", "ignored", "Stop silence on the channel"},
	{"ari.channels.senddtmf.>", "string", "ignored", "Send the DTMF to the channel"},
	{"ari.channels.moh.>", "string", "ignored", "Play the given music on hold to the channel"},
	{"ari.channels.stopmoh.>", "ignored", "ignored", "Stop all music on hold on the channel"},
	{"ari.channels.play.>", "nc.PlayRequest", "ignored", "Play the given mediaURI on the channel"},

	{"ari.devices.list", "ignored", "[]string", "List the devices"},
	{"ari.devices.data.>", "ignored", "ari.DeviceStateData", "Get the device state"},
	{"ari.devices.update.>", "string", "ignored", "Update the device state"},
	{"ari.devices.delete.>", "ignored", "ignored", "Delete the device state"},

	{"ari.mailboxes.all", "ignored", "[]string", "List the mailboxes"},
	{"ari.mailboxes.data.>", "ignored", "ari.MailboxData", "Get the mailbox data"},
	{"ari.mailboxes.update.>", "ignored", "nc.UpdateMailboxRequest", "Update the mailbox state"},
	{"ari.mailboxes.delete.>", "ignored", "delete", "Delete the mailbox state"},

	{"ari.playback.data.>", "ignored", "ari.PlaybackData", "Get the playback data"},
	{"ari.playback.control.>", "string", "ignored", "Send the control command to the playback"},
	{"ari.playback.stop.>", "ignored", "ignored", "Stop the playback"},
}
