package ari

var (
	DefaultBaseUri  = "http://216.66.0.147:8088/ari"
	DefaultWsUri    = "ws://216.66.0.147:8088/ari/events"
	DefaultUsername = "enswitch"
	DefaultSecret   = "enswitchDev"
)

var DefaultClient, _ = NewClient("default", DefaultBaseUri, DefaultWsUri, DefaultUsername, DefaultSecret)
