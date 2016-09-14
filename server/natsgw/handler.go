package natsgw

// Reply is a function which, when called, replies to the request via the
// response object or error.
type Reply func(interface{}, error)

// A Handler is a handler which provides the subject, the raw body, and the reply function
type Handler func(subj string, request []byte, reply Reply)

// NamedHandler returns a handler which is given the rest of string after the prefix
func NamedHandler(prefixLen int, h Handler) Handler {
	return func(subj string, request []byte, reply Reply) {
		name := subj[prefixLen:]
		h(name, request, reply)
	}
}
