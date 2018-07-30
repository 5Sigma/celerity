package celerity

import (
	"testing"

	"github.com/spf13/afero"
)

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

func TestLocalFileRoute(t *testing.T) {
	adapter := NewMEMAdapter()
	FSAdapter = adapter
	r := &LocalFileRoute{
		Path:      "/public/test.txt",
		LocalPath: "/files/test.txt",
	}

	t.Run("without file", func(t *testing.T) {
		if ok, _ := r.Match(GET, "/public/test.txt"); ok {
			t.Error("should not match non existant file")
		}
	})
	t.Run("with file", func(t *testing.T) {
		afero.WriteFile(adapter.MEMFS, "/files/test.txt", []byte("public file"), 0755)
		if ok, _ := r.Match(GET, "/public/test.txt"); !ok {
			t.Error("should match existing file")
		}
	})
}

func TestLocalPathRoute(t *testing.T) {
	adapter := NewMEMAdapter()
	FSAdapter = adapter
	r := &LocalPathRoute{
		Path:      "/public/",
		LocalPath: "/files/",
	}
	t.Run("without file", func(t *testing.T) {
		if ok, _ := r.Match(GET, "/public/test.txt"); ok {
			t.Error("should not match non existant file")
		}
	})
	t.Run("with file", func(t *testing.T) {
		afero.WriteFile(adapter.MEMFS, "/files/test.txt", []byte("public file"), 0755)
		if ok, _ := r.Match(GET, "/public/test.txt"); !ok {
			t.Error("should match existing file")
		}
	})
}
