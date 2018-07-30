package celerity

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/5Sigma/vox"
	"github.com/tidwall/gjson"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// Context A request context object
type Context struct {
	Method      string
	Request     *http.Request
	RequestID   string
	URLParams   Params
	QueryParams Params
	properties  map[string]interface{}
	Log         *vox.Vox
	Response    Response
	Env         string
	ScopedPath  string
	data        []byte
	Writer      http.ResponseWriter
	Server      *Server
}

// NewContext Create a new context object
func NewContext() Context {

	return Context{
		URLParams:   Params(map[string]string{}),
		QueryParams: Params(map[string]string{}),
		properties:  map[string]interface{}{},
		Response:    NewResponse(),
		Env:         viper.GetString("env"),
		RequestID:   strings.Replace(uuid.New().String(), "-", "", -1),
	}
}

// RequestContext - Creates a new context and sets its request.
func RequestContext(r *http.Request) Context {
	c := NewContext()
	c.Request = r
	c.ScopedPath = r.URL.Path
	return c
}

// Header Request headers
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Set Set an arbitrary value in the context
func (c *Context) Set(key string, v interface{}) {
	c.properties[key] = v
}

// Get Return an arbitrary value from the context that was set with the Set
// function.
func (c *Context) Get(key string) interface{} {
	if v, ok := c.properties[key]; ok {
		return v
	}
	return nil
}

// R Alias for Respond
func (c *Context) R(obj interface{}) Response {
	return c.Respond(obj)
}

// Respond Respond with an object
func (c *Context) Respond(obj interface{}) Response {
	c.Response.StatusCode = 200
	c.Response.Data = obj
	return c.Response
}

// F is an alias for the Fail function
func (c *Context) F(err error) Response {
	return c.Fail(err)
}

// Body returns the request body
func (c *Context) Body() []byte {
	if len(c.data) == 0 {
		buf, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return []byte{}
		}
		c.data = buf
		return buf
	}
	return c.data
}

// Fail is used for unrecoverable and internal errors. In a production
// environment the error is not passed to the client.
// message.
func (c *Context) Fail(err error) Response {
	c.Response.StatusCode = 500
	c.Response.Data = nil
	if viper.GetString("env") == PROD {
		c.Response.Error = errors.New("the request could not be processed")
	} else {
		c.Response.Error = err
	}

	return c.Response
}

// E is an alias for the Error function.
func (c *Context) E(status int, err error) Response {
	return c.Error(status, err)
}

// Error - Returns a erorr and outputs the passed error message.
func (c *Context) Error(status int, err error) Response {
	c.Response.StatusCode = status
	c.Response.Data = nil
	c.Response.Error = err
	return c.Response
}

// Raw returns a response configured to output a raw []byte resposne. This
// resposne will also skip the response transformation adapter.
func (c *Context) Raw(b []byte) Response {
	c.Response.SetRaw(b)
	return c.Response
}

// Extract - Unmarshal request data into a structure.
func (c *Context) Extract(obj interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	err := decoder.Decode(obj)
	return err
}

// ExtractValue extracts a value from the request body at a specific JSON node.
func (c *Context) ExtractValue(path string) gjson.Result {
	return gjson.Get(string(c.Body()), path)
}

//Params - Stores key value params for URL parameters and query parameters. It
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
