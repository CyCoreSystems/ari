package native

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rotisserie/eris"
)

// MaxIdleConnections is the maximum number of idle web client
// connections to maintain.
var MaxIdleConnections = 20

// RequestTimeout describes the maximum amount of time to wait
// for a response to any request.
var RequestTimeout = 2 * time.Second

// RequestError describes an error with an error Code.
type RequestError interface {
	error
	Code() int
}

type requestError struct {
	statusCode int
	text       string
}

// Error returns the request error as a string.
func (e *requestError) Error() string {
	return e.text
}

// Code returns the status code from the request.
func (e *requestError) Code() int {
	return e.statusCode
}

// CodeFromError extracts and returns the code from an error, or
// 0 if not found.
func CodeFromError(err error) int {
	if reqerr, ok := err.(RequestError); ok {
		return reqerr.Code()
	}

	return 0
}

func maybeRequestError(resp *http.Response) RequestError {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 2xx response: All good.
		return nil
	}

	return &requestError{
		text:       "Non-2XX response: " + resp.Status,
		statusCode: resp.StatusCode,
	}
}

// MissingParams is an error message response emitted when a request
// does not contain required parameters
type MissingParams struct {
	// Message
	Type   string   `json:"type"`
	Params []string `json:"params"` // List of missing parameters which are required
}

// get calls the ARI server with a GET request
func (c *Client) get(url string, resp interface{}) error {
	url = c.Options.URL + url

	return c.makeRequest("GET", url, resp, nil)
}

// post calls the ARI server with a POST request.
func (c *Client) post(requestURL string, resp interface{}, req interface{}) error {
	url := c.Options.URL + requestURL
	return c.makeRequest("POST", url, resp, req)
}

// put calls the ARI server with a PUT request.
func (c *Client) put(url string, resp interface{}, req interface{}) error {
	url = c.Options.URL + url

	return c.makeRequest("PUT", url, resp, req)
}

// del calls the ARI server with a DELETE request
func (c *Client) del(url string, resp interface{}, req string) error {
	url = c.Options.URL + url
	if req != "" {
		url = url + "?" + req
	}

	return c.makeRequest("DELETE", url, resp, nil)
}

func (c *Client) makeGenericRequest(method, url string, req interface{}, contentType string) (*http.Response, error) {
	var (
		reqBody io.Reader
		err     error
	)
	if req != nil {
		reqBody, err = structToRequestBody(req)
		if err != nil {
			return nil, eris.Wrap(err, "failed to marshal request")
		}
	}

	var r *http.Request

	if r, err = http.NewRequest(method, url, reqBody); err != nil {
		return nil, eris.Wrap(err, "failed to create request")
	}

	r.Header.Set("Content-Type", "application/json")

	if c.Options.Username != "" {
		r.SetBasicAuth(c.Options.Username, c.Options.Password)
	}

	return c.httpClient.Do(r)
}

func (c *Client) makeRequest(method, url string, resp interface{}, req interface{}) (err error) {
	ret, err := c.makeGenericRequest(method, url, req, "application/json")
	if err != nil {
		return eris.Wrap(err, "failed to make request")
	}

	defer ret.Body.Close() //nolint:errcheck

	if resp != nil {
		err = json.NewDecoder(ret.Body).Decode(resp)
		if err != nil {
			return eris.Wrap(err, "failed to decode response")
		}
	}

	return maybeRequestError(ret)
}

func structToRequestBody(req interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)

	if req != nil {
		if err := json.NewEncoder(buf).Encode(req); err != nil {
			return nil, err
		}
	}

	return buf, nil
}
