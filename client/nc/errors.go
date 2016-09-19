package nc

import (
	"errors"
	"fmt"
)

type wrappedError struct {
	Message string
	Err     error
}

func (err *wrappedError) Cause() error {
	return err.Err
}

func (err *wrappedError) Error() string {
	return err.Message + ": " + err.Err.Error()
}

// remote error, used wrap the error response before sending

type remoteError struct {
	Subj string
	Err  error
}

func (r *remoteError) Cause() error {
	return r.Err
}

func (r *remoteError) Remote() bool {
	return true
}

func (r *remoteError) Error() string {
	return fmt.Sprintf("Remote Error in Endpoint '%s': %v", r.Subj, r.Err)
}

type codedError struct {
	err  error
	code int
}

func (err *codedError) Error() string {
	return err.err.Error()
}

func (err *codedError) Code() int {
	return err.code
}

type causer interface {
	Cause() error
}

type coded interface {
	Code() int
}

// ErrorToMap converts an error type to a key-value map
func ErrorToMap(err error, parent string) map[string]interface{} {
	data := make(map[string]interface{})
	if parent == err.Error() {
		// NOTE: this is done because of how errors.Wrap works, internally,
		// to build a stacktrace. We end up with duplicate
		// entries in the tree of errors.
		if c, ok := err.(causer); ok {
			return ErrorToMap(c.Cause(), parent)
		}
	}
	data["message"] = err.Error()
	if c, ok := err.(coded); ok {
		data["code"] = c.Code()
	}
	if c, ok := err.(causer); ok {
		data["cause"] = ErrorToMap(c.Cause(), data["message"].(string))
	}
	return data
}

// MapToError converts a JSON parsed map to an error type
func MapToError(i map[string]interface{}) error {
	msg, _ := i["message"].(string)
	code, codeOK := i["code"].(int)
	cause, causeOK := i["cause"].(map[string]interface{})

	err := errors.New(msg)

	if codeOK {
		err = &codedError{err, code}
	}

	if causeOK {
		causeError := MapToError(cause)
		l := len(msg) - len(causeError.Error())
		msg = msg[:l-2]
		err = &wrappedError{msg, causeError}
	}

	return err
}
