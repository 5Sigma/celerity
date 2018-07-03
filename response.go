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
}

// NewResponse - Create a new response object
func NewResponse() Response {
	return Response{
		Meta:   map[string]interface{}{},
		Header: http.Header{},
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
