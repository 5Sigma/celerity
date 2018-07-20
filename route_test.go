package celerity

import (
	"testing"

	"github.com/spf13/afero"
)

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
