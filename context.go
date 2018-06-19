package celerity

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Context - A request context object
type Context struct {
	Method      string
	Request     *http.Request
	RequestID   string
	URLParams   Params
	QueryParams Params
	properties  map[string]interface{}
	Response    Response
}

// NewContext - Create a new context object
func NewContext() Context {
	return Context{
		URLParams:   Params(map[string]string{}),
		QueryParams: Params(map[string]string{}),
		properties:  map[string]interface{}{},
		Response:    NewResponse(),
	}
}

// Header - Request headers
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Set - Set an arbitrary value in the context
func (c *Context) Set(key string, v interface{}) {
	c.properties[key] = v
}

// Get - Return an arbitrary value from the context that was set with the Set
// function.
func (c *Context) Get(key string) interface{} {
	if v, ok := c.properties[key]; ok {
		return v
	}
	return nil
}

// R - Alias for Respond
func (c *Context) R(obj interface{}) Response {
	return c.Respond(obj)
}

// Respond - Respond with an object
func (c *Context) Respond(obj interface{}) Response {
	c.Response.StatusCode = 200
	c.Response.Data = obj
	return c.Response
}

// Fail - Returns a 500 internal server erorr and outputs the passed error
// message.
func (c *Context) Fail(err error) Response {
	c.Response.StatusCode = 500
	c.Response.Data = nil
	c.Response.Error = err
	return c.Response
}

// Extract - Unmarshal request data into a structure.
func (c *Context) Extract(obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(obj)
	return err
}

//Params - Stores key value params for URL parameters and query paramters. It
//offers several helper methods for getting results for a key.
type Params map[string]string

// SetParams - Sets the url parameters for the request.
func (c *Context) SetParams(params map[string]string) {
	c.URLParams = Params(params)
}

// SetQueryParamsFromURL - Sets the query parameters for the request.
func (c *Context) SetQueryParamsFromURL(u *url.URL) {
	m := map[string]string{}
	for k, v := range u.Query() {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	c.QueryParams = Params(m)
}

// String - Returns the string value for a parameter key or ""
func (p Params) String(key string) string {
	if v, ok := p[key]; ok {
		return v
	}
	return ""
}

// Int - Returns the int value for a parameter key or ""
func (p Params) Int(key string) int {
	i, err := strconv.Atoi(p[key])
	if err != nil {
		return -1
	}
	return i
}
