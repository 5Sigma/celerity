package celerity

import "testing"

func TestRouteMatch(t *testing.T) {
	{
		r := Route{
			Path: "/users",
		}
		if !r.Match("/users") {
			t.Error("Did not match valid path")
		}
		if r.Match("/bad") {
			t.Error("Did match invalid path")
		}
	}
}
