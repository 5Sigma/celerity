package celerity

import "errors"

// Response - The response object reurned by an endpoint.
type Response struct {
	StatusCode int
	Data       interface{}
	Error      error
	Meta       map[string]interface{}
	Headers    map[string]string
}

// NewResponse - Create a new response object
func NewResponse() Response {
	return Response{
		Meta:    map[string]interface{}{},
		Headers: map[string]string{},
	}
}

// NewErrorResponse - Return a new error response
func NewErrorResponse(status int, message string) Response {
	return Response{
		StatusCode: status,
		Error:      errors.New(message),
	}
}
