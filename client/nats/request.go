package nats

import (
	"fmt"
	"strings"
	"time"
)

// RequestTimeout describes the maximum amount of time to wait
// for a response to any request.
var RequestTimeout = 2 * time.Second

// Get calls the ARI server with a GET request
func (conn *Conn) Get(url string, items []interface{}, ret interface{}) error {

	r := make([]interface{}, len(items))
	for i := range items {
		r[i] = ""
	}

	url = fmt.Sprintf(url, r...)

	newUrl := ""
	for _, s := range strings.Split(url, "/") {
		if s == "" {
			continue
		}

		newUrl = newUrl + "." + s
	}

	err := conn.c.Request(newUrl+".get", "", &ret, RequestTimeout)
	return err
}

// Post calls the ARI server with a POST request.
func (conn *Conn) Post(url string, items []interface{}, ret interface{}, req interface{}) error {
	err := conn.c.Request(url+".post", &req, &ret, RequestTimeout)
	return err
}

// Put calls the ARI server with a PUT request.
func (conn *Conn) Put(url string, items []interface{}, ret interface{}, req interface{}) error {
	err := conn.c.Request(url+".put", &req, &ret, RequestTimeout)
	return err
}

// Delete calls the ARI server with a DELETE request
func (conn *Conn) Delete(url string, items []interface{}, ret interface{}, req string) error {
	err := conn.c.Request(url+".delete", &req, &ret, RequestTimeout)
	return err
}
