package generic

// Conn is a connection to a  ARI server
type Conn interface {
	// Get calls the ARI server with a GET request
	Get(url string, parts []interface{}, ret interface{}) error

	// Post calls the ARI server with a POST request.
	Post(url string, parts []interface{}, ret interface{}, req interface{}) error

	// Put calls the ARI server with a PUT request.
	Put(url string, parts []interface{}, ret interface{}, req interface{}) error

	// Delete calls the ARI server with a DELETE request
	Delete(url string, parts []interface{}, ret interface{}, req string) error
}
