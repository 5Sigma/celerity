package celerity

import "testing"

func TestRouteMatch(t *testing.T) {
	{
		r := Route{
			Method: GET,
			Path:   "/users",
		}
		if !r.Match(GET, "/users") {
			t.Error("Did not match valid path")
		}
		if r.Match(GET, "/bad") {
			t.Error("Did match invalid path")
		}
	}
	{
		r := Route{
			Method: POST,
			Path:   "/users",
		}
		if r.Match(GET, "/users") {
			t.Error("should not match incorrect method")
		}
	}
}
