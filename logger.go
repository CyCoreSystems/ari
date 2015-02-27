package ari

import (
	"io/ioutil"
	"log"
)

// Logger
var Logger *log.Logger

func init() {
	// Null logger, by default
	Logger = log.New(ioutil.Discard, "restclient", log.LstdFlags|log.Lshortfile)
}
