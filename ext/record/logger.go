package record

import "github.com/inconshreveable/log15"

// Logger defaults to a discard handler (null output).
// If you wish to enable logging, you can set your own
// handler like so:
//
//	ari.Logger.SetHandler(log15.StderrHandler)
var Logger = log15.New()

func init() {
	// Null logger, by default
	Logger.SetHandler(log15.DiscardHandler())
}
