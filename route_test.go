package celerity

import "testing"

func TestRouteMatch(t *testing.T) {
	{
		r := Route{
			Method: GET,
			Path:   "/users",
		}
		if ok, _ := r.Match(GET, "/users"); !ok {
			t.Error("Did not match valid path")
		}
		if ok, _ := r.Match(GET, "/bad"); ok {
			t.Error("Did match invalid path")
		}
	}
	{
		r := Route{
			Method: POST,
			Path:   "/users",
		}
		if ok, _ := r.Match(GET, "/users"); ok {
			t.Error("should not match incorrect method")
		}
	}
}
