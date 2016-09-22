package ari

var (
	DefaultUsername = "ariTest"
	DefaultSecret   = "ariDev"
)

var DefaultClient = NewClient(&Options{
	Application: "default",
	Username:    "ariTest",
	Password:    "ariDev",
})
