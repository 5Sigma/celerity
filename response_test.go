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

func TestRaw(t *testing.T) {
	r := NewResponse()
	if r.IsRaw() {
		t.Error("response should not be raw")
	}
	r.SetRaw([]byte("test"))
	if !r.IsRaw() {
		t.Error("response should be raw")
	}
	if string(r.Raw()) != "test" {
		t.Errorf("raw data not correct: %s", string(r.Raw()))
	}
}
