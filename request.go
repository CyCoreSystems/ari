package ari

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
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
		text:       "Non-2XX response: " + http.StatusText(resp.StatusCode),
		statusCode: resp.StatusCode,
	}
}

// MissingParams is an error message response emitted when a request
// does not contain required parameters
type MissingParams struct {
	Message
	Params []string `json:"params"` // List of missing parameters which are required
}

func (c *Client) assureHTTPClient() {
	if c.httpClient == nil {
		c.httpClient = gorequest.New().Timeout(RequestTimeout)
		if c.Options.Username != "" {
			c.httpClient = c.httpClient.SetBasicAuth(c.Options.Username, c.Options.Password)
		}
	}
}

// Get wraps gorequest.Get with the complete url
// It calls the ARI server with a GET request
func (c *Client) Get(url string, ret interface{}) error {
	c.assureHTTPClient()
	finalURL := c.Options.URL + url
	resp, body, errs := c.httpClient.Get(finalURL).EndBytes()
	if errs != nil {
		var errString string
		for _, e := range errs {
			errString += ": " + e.Error()
		}
		return fmt.Errorf("Error making request: %s", errString)
	}

	if ret != nil {
		err := json.Unmarshal(body, ret)
		if err != nil {
			return err
		}
	}
	return maybeRequestError(resp)
}

// Post is a shorthand for MakeRequest("POST",url,ret,req)
// It calls the ARI server with a POST request
// Uses gorequest.PostForm since ARI returns bad request otherwise
func (c *Client) Post(url string, ret interface{}, req interface{}) error {
	c.assureHTTPClient()
	finalURL := c.Options.URL + url
	r := c.httpClient.Post(finalURL).Type("form")
	if req != nil {
		r = r.SendStruct(req)
	}
	resp, body, errs := r.EndBytes()
	if errs != nil {
		var errString string
		for _, e := range errs {
			errString += ": " + e.Error()
		}
		return fmt.Errorf("Error making request: %s", errString)
	}

	if ret != nil {
		err := json.Unmarshal(body, ret)
		if err != nil {
			return err
		}
	}
	return maybeRequestError(resp)
}

// Put is a shorthand for MakeRequest("PUT",url,ret,req)
// It calls the ARI server with a PUT request
func (c *Client) Put(url string, ret interface{}, req interface{}) error {
	c.assureHTTPClient()
	finalURL := c.Options.URL + url
	r := c.httpClient.Put(finalURL).Type("form")
	if req != nil {
		r = r.Send(req)
	}

	resp, body, errs := r.EndBytes()
	if errs != nil {
		var errString string
		for _, e := range errs {
			errString += ": " + e.Error()
		}
		return fmt.Errorf("Error making request: %s", errString)
	}

	if ret != nil {
		err := json.Unmarshal(body, ret)
		if err != nil {
			return err
		}
	}
	return maybeRequestError(resp)
}

// Delete is a shorthand for MakeRequest("DELETE",url,nil,nil)
// It calls the ARI server with a DELETE request
func (c *Client) Delete(url string, ret interface{}, req interface{}) error {
	c.assureHTTPClient()
	finalURL := c.Options.URL + url

	r := c.httpClient.Delete(finalURL)
	if req != nil {
		r = r.Query(req)
	}

	resp, body, errs := r.EndBytes()
	if errs != nil {
		var errString string
		for _, e := range errs {
			errString += ": " + e.Error()
		}
		return fmt.Errorf("Error making request: %s", errString)
	}

	if ret != nil {
		err := json.Unmarshal(body, ret)
		if err != nil {
			return err
		}
	}
	return maybeRequestError(resp)
}
