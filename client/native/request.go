package native

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	//Message
	Type   string   `json:"type"`
	Params []string `json:"params"` // List of missing parameters which are required
}

// Get calls the ARI server with a GET request
func Get(conn *Conn, url string, ret interface{}) error {

	finalURL := conn.Options.URL + url

	httpReq, err := buildRequest(conn, "GET", finalURL, "", nil)
	if err != nil {
		return err
	}

	resp, err := conn.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if ret != nil {
		if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
			return err
		}
	}

	return maybeRequestError(resp)
}

// Post calls the ARI server with a POST request.
func Post(conn *Conn, requestURL string, ret interface{}, req interface{}) error {

	finalURL := conn.Options.URL + requestURL

	requestBody, contentType, err := structToRequestBody(req)
	if err != nil {
		return err
	}

	httpReq, err := buildRequest(conn, "POST", finalURL, contentType, requestBody)
	if err != nil {
		return err
	}

	resp, err := conn.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if ret != nil {
		if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
			return err
		}
	}

	return maybeRequestError(resp)
}

// Put calls the ARI server with a PUT request.
func Put(conn *Conn, url string, ret interface{}, req interface{}) error {

	finalURL := conn.Options.URL + url

	requestBody, contentType, err := structToRequestBody(req)
	if err != nil {
		return err
	}

	httpReq, err := buildRequest(conn, "PUT", finalURL, contentType, requestBody)
	if err != nil {
		return err
	}

	resp, err := conn.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if ret != nil {
		if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
			return err
		}
	}

	return maybeRequestError(resp)
}

// Delete calls the ARI server with a DELETE request
func Delete(conn *Conn, url string, ret interface{}, req string) error {

	finalURL := conn.Options.URL + url
	if req != "" {
		finalURL = finalURL + "?" + req
	}

	httpReq, err := buildRequest(conn, "DELETE", finalURL, "", nil)
	if err != nil {
		return err
	}

	resp, err := conn.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if ret != nil {
		if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
			return err
		}
	}

	return maybeRequestError(resp)
}

func buildRequest(conn *Conn, method string, finalURL string, contentType string, body io.Reader) (*http.Request, error) {

	if contentType == "" {
		contentType = "application/json"
	}

	ret, err := http.NewRequest(method, finalURL, body)
	if err != nil {
		return nil, err
	}

	if conn.Options.Username != "" {
		ret.SetBasicAuth(conn.Options.Username, conn.Options.Password)
	}

	ret.Header.Set("Content-Type", contentType)

	return ret, nil
}

func structToRequestBody(req interface{}) (io.Reader, string, error) {
	buf := bytes.NewBuffer([]byte(""))
	if req != nil {
		if err := json.NewEncoder(buf).Encode(req); err != nil {
			return nil, "", err
		}
	}

	return buf, "application/json", nil
}
