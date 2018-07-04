package celerity

import "testing"

func TestRoutePathMatch(t *testing.T) {
	{
		var rp RoutePath = "/users/:id/pages/:pageid"

		{
			ok, extra := rp.Match("/users/1/pages/23")
			if !ok {
				t.Error("RoutePath did not match a correct path")
			}
			if extra != "" {
				t.Errorf("should not have any remaining path: %s", extra)
			}
		}

		{
			ok, _ := rp.Match("/users/1/badpath/23")
			if ok {
				t.Error("RoutePath did match a bad path")
			}
		}

		{
			ok, extra := rp.Match("/users/1/pages/23/extra/path")
			if !ok {
				t.Error("Did not match a path with extra data")
			}
			if extra != "extra/path" {
				t.Errorf("did not correctly return remaining path: %s", extra)
			}
		}
	}
	{
		var rp RoutePath = "/"
		{
			ok, xtra := rp.Match("/users")
			if !ok {
				t.Error("Did not match valid path")
			}
			if xtra != "/users" {
				t.Errorf("Extra path not correct: %s", xtra)
			}
		}
	}
}

func TestRoutePathGetURLParams(t *testing.T) {
	var rp RoutePath = "/users/:id/pages/:pageid"
	{
		params := rp.GetURLParams("/users/22/pages/12")
		if params["id"] != "22" {
			t.Errorf("id param should be 22 was %s", params["id"])
		}
	}
}

func TestWildCardMatch(t *testing.T) {
	{
		var rp RoutePath = "*"
		ok, _ := rp.Match("/some/route/to/test")
		if !ok {
			t.Error("wldcard path did not match")
		}
	}
	{
		var rp RoutePath = "/some/route/*"
		ok, _ := rp.Match("/some/route/to/test")
		if !ok {
			t.Error("wldcard path did not match")
		}
	}
}
