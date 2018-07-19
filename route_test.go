package celerity

import "testing"

func TestBasicRouteMatch(t *testing.T) {
	{
		r := &BasicRoute{
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
		r := &BasicRoute{
			Method: POST,
			Path:   "/users",
		}
		if ok, _ := r.Match(GET, "/users"); ok {
			t.Error("should not match incorrect method")
		}
	}
}
