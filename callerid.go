package ari

//NOTE: Direct translation from ARI client 2.0

// CallerID describes the name and number which
// identifies the caller to other endpoints
type CallerID struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// String returns the stringified callerid
func (cid *CallerID) String() string {
	return cid.Name + "<" + cid.Number + ">"
}
