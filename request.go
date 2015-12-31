package ari

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

var MaxIdleConnections int = 20
var RequestTimeout time.Duration = 2 * time.Second

type RequestError interface {
	error
	Code() int
}

type requestError struct {
	statusCode int
	text       string
}

func (e *requestError) Error() string {
	return e.text
}

func (e *requestError) Code() int {
	return e.statusCode
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

func (c *Client) assureHttpClient() {
	if c.httpClient == nil {
		c.httpClient = gorequest.New().Timeout(RequestTimeout)
		if c.username != "" {
			c.httpClient = c.httpClient.SetBasicAuth(c.username, c.password)
		}
	}
}

// AriGet wraps restclient.Get with the complete url
// It calls the ARI server with a GET request
func (c *Client) AriGet(url string, ret interface{}) error {
	c.assureHttpClient()
	finalUrl := c.Url + url
	resp, body, errs := c.httpClient.Get(finalUrl).EndBytes()
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

// AriPost is a shorthand for MakeRequest("POST",url,ret,req)
// It calls the ARI server with a POST request
// Uses restclient.PostForm since ARI returns bad request otherwise
func (c *Client) AriPost(url string, ret interface{}, req interface{}) error {
	c.assureHttpClient()
	finalUrl := c.Url + url
	r := c.httpClient.Post(finalUrl).Type("form")
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

// AriPut is a shorthand for MakeRequest("PUT",url,ret,req)
// It calls the ARI server with a PUT request
func (c *Client) AriPut(url string, ret interface{}, req interface{}) error {
	c.assureHttpClient()
	finalUrl := c.Url + url
	r := c.httpClient.Put(finalUrl).Type("form")
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

// AriDelete is a shorthand for MakeRequest("DELETE",url,nil,nil)
// It calls the ARI server with a DELETE request
func (c *Client) AriDelete(url string, ret interface{}, req interface{}) error {
	c.assureHttpClient()
	finalUrl := c.Url + url

	r := c.httpClient.Delete(finalUrl)
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
