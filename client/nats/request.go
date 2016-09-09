package nats

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// RequestTimeout describes the maximum amount of time to wait
// for a response to any request.
var RequestTimeout = 2 * time.Second

func convertURL(requestURL string, method string, items ...interface{}) string {

	r := make([]interface{}, len(items))
	for i := range items {
		r[i] = ""
	}

	requestURL = fmt.Sprintf(requestURL, r...)

	new := ""
	for _, s := range strings.Split(requestURL, "/") {
		if s == "" {
			continue
		}
		new = new + "." + s
	}

	if len(items) == 0 {
		return new[1:] + "." + method
	}

	var itemsStr []string
	for _, i := range items {

		if i.(string)[0] == '?' {
			u, err := url.ParseQuery(i.(string)[1:])
			if err != nil {
				panic(err)
			}

			for key, val := range u {
				for _, v := range val {
					itemsStr = append(itemsStr, key)
					itemsStr = append(itemsStr, v)
				}
			}
			continue
		}

		itemsStr = append(itemsStr, fmt.Sprintf("%v", i))
	}

	return new[1:] + "." + method + "." + strings.Join(itemsStr, ".")
}

// Get calls the ARI server with a GET request
func (conn *Conn) Get(url string, items []interface{}, ret interface{}) error {
	url = convertURL(url, "get", items...)
	err := conn.c.Request(url, "", &ret, RequestTimeout)
	return err
}

// Post calls the ARI server with a POST request.
func (conn *Conn) Post(url string, items []interface{}, ret interface{}, req interface{}) error {
	url = convertURL(url, "post", items...)
	err := conn.c.Request(url, &req, &ret, RequestTimeout)
	return err
}

// Put calls the ARI server with a PUT request.
func (conn *Conn) Put(url string, items []interface{}, ret interface{}, req interface{}) error {
	url = convertURL(url, "put", items...)
	err := conn.c.Request(url, &req, &ret, RequestTimeout)
	return err
}

// Delete calls the ARI server with a DELETE request
func (conn *Conn) Delete(url string, items []interface{}, ret interface{}, req string) error {
	url = convertURL(url, "delete", items...)
	err := conn.c.Request(url, &req, &ret, RequestTimeout)
	return err
}
