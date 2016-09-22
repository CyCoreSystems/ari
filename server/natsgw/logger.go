package natsgw

// Logger is the generic logging interface for the natsgw server
type Logger interface {
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Error(msg string, v ...interface{})
}
