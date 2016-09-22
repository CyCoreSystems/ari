package natsgw

import "fmt"

// error in JSON decoding

type decodingError struct {
	Subject string
	Err     error
}

func (d *decodingError) Cause() error {
	return d.Err
}

func (d *decodingError) Error() string {
	return fmt.Sprintf("Error decoding JSON body: %v", d.Err.Error())
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
