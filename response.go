package celerity

import (
	"errors"
	"net/http"
)

// Response - The response object reurned by an endpoint.
type Response struct {
	StatusCode int
	Data       interface{}
	Error      error
	Meta       map[string]interface{}
	Header     http.Header
	raw        []byte
	Handled    bool
}

// NewResponse - Create a new response object
func NewResponse() Response {
	return Response{
		Meta:       map[string]interface{}{},
		Header:     http.Header{},
		StatusCode: 200,
	}
}

// NewErrorResponse - Return a new error response
func NewErrorResponse(status int, message string) Response {
	return Response{
		StatusCode: status,
		Error:      errors.New(message),
	}
}

// StatusText returns the text version of the StatusCode
func (r *Response) StatusText() string {
	return http.StatusText(r.StatusCode)
}

// Success returns true if the response was marked succcessful and if an error
// is not present
func (r *Response) Success() bool {
	return r.Error == nil
}

// IsRaw determens if the response is a raw response
func (r *Response) IsRaw() bool {
	return len(r.raw) > 0
}

// SetRaw sets the responses raw output
func (r *Response) SetRaw(b []byte) {
	r.raw = b
	r.Data = nil
}

// Raw returns the raw data for the request
func (r *Response) Raw() []byte {
	return r.raw
}

// Status sets the status code for the response
func (r Response) Status(code int) Response {
	r.StatusCode = code
	return r
}

// Respond sets teh response data
func (r Response) Respond(data interface{}) Response {
	r.Data = data
	return r
}

// R aliases Respond
func (r Response) R(data interface{}) Response {
	return r.Respond(data)
}

// S aliases Status
func (r Response) S(code int) Response {
	return r.Status(code)
}

// MetaValue sets the metadata key for the response
func (r Response) MetaValue(k string, v interface{}) Response {
	r.Meta[k] = v
	return r
}
