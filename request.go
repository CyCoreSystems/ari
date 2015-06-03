package ari

import "github.com/CyCoreSystems/restclient"

type RequestError restclient.Error

// MissingParams is an error message response emitted when a request
// does not contain required parameters
type MissingParams struct {
	Message
	Params []string `json:"params"` // List of missing parameters which are required
}

// AriGet wraps restclient.Get with the complete url
// It calls the ARI server with a GET request
func (c *Client) AriGet(url string, ret interface{}) RequestError {
	finalUrl := c.Url + url
	return restclient.Get(finalUrl, restclient.Auth{c.username, c.password}, ret)
}

// AriPost is a shorthand for MakeRequest("POST",url,ret,req)
// It calls the ARI server with a POST request
// Uses restclient.PostForm since ARI returns bad request otherwise
func (c *Client) AriPost(url string, ret interface{}, req interface{}) RequestError {
	finalUrl := c.Url + url
	return restclient.PostForm(finalUrl, restclient.Auth{c.username, c.password}, req, ret)
}

// AriPut is a shorthand for MakeRequest("PUT",url,ret,req)
// It calls the ARI server with a PUT request
func (c *Client) AriPut(url string, ret interface{}, req interface{}) RequestError {
	finalUrl := c.Url + url
	return restclient.Put(finalUrl, restclient.Auth{c.username, c.password}, req, ret)
}

// AriDelete is a shorthand for MakeRequest("DELETE",url,nil,nil)
// It calls the ARI server with a DELETE request
func (c *Client) AriDelete(url string, ret interface{}, req interface{}) RequestError {
	finalUrl := c.Url + url
	return restclient.Delete(finalUrl, restclient.Auth{c.username, c.password}, req, ret)
}
