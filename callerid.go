package ari

import "errors"

//NOTE: Direct translation from ARI client 2.0

// CallerID describes the name and number which
// identifies the caller to other endpoints
type CallerID struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// CallerIDFromString interprets the provided string
// as a CallerID.  Usually, this string will be of the following forms:
//   - "Name" <number>
//   - <number>
//   - "Name" number
func CallerIDFromString(src string) (*CallerID, error) {
	//TODO: implement complete callerid parser
	return nil, errors.New("CallerIDFromString not yet implemented")
}

// String returns the stringified callerid
func (cid *CallerID) String() string {
	return cid.Name + "<" + cid.Number + ">"
}
