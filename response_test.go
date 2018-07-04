package celerity

import "testing"

func TestNewResponse(t *testing.T) {
	r := NewResponse()
	if r.Meta == nil {
		t.Error("meta not initialized")
	}
}

func TestNewErrorResponse(t *testing.T) {
	r := NewErrorResponse(404, "not found")
	if r.StatusCode != 404 {
		t.Errorf("status code should be 404 was %d", r.StatusCode)
	}
	if r.Error.Error() != "not found" {
		t.Errorf("error message incorrect: %s", r.Error.Error())
	}
	if r.StatusText() != "Not Found" {
		t.Errorf("status text not correct: %s", r.StatusText())
	}
	if r.Success() != false {
		t.Errorf("success returned true for errored response")
	}
}
