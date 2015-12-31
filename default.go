package ari

var (
	DefaultBaseUri  = "http://216.66.0.147:8088/ari"
	DefaultWsUri    = "ws://216.66.0.147:8088/ari/events"
	DefaultUsername = "ariTest"
	DefaultSecret   = "ariDev"
)

var DefaultClient, _ = NewClient("default", DefaultBaseUri, DefaultWsUri, DefaultUsername, DefaultSecret)
