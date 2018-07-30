package celerity

import (
	"fmt"
	"os"
	"path/filepath"
)

// Route is an interface that can be implemented by any objects wishing to
// process url calls.
type Route interface {
	Match(method, path string) (bool, string)
	Handle(Context) Response
	RoutePath() RoutePath
}

// BasicRoute - A basic route in the server.
type BasicRoute struct {
	Path    RoutePath
	Method  string
	Handler RouteHandler
}

// RouteHandler - The handler function that gets called when a route is invoked.
type RouteHandler func(Context) Response

// Match - Matches the routes path against the incomming url
func (r *BasicRoute) Match(method, path string) (bool, string) {
	ok, xtra := r.Path.Match(path)
	return (ok && method == r.Method && xtra == ""), xtra
}

// Handle processes the request by passing it to the RouteHandler function
func (r *BasicRoute) Handle(c Context) Response {
	return r.Handler(c)
}

// RoutePath returns the RoutePath for the route.
func (r *BasicRoute) RoutePath() RoutePath {
	return r.Path
}

func (r *BasicRoute) String() string {
	return fmt.Sprint(r.Method, "\t", r.Path)
}

// LocalFileRoute will handle serving a single file from the local file system
// to a given path.
type LocalFileRoute struct {
	Path      RoutePath
	LocalPath string
}

// Match checks if the incoming path matches the route path and if the local
// file exists.
func (l *LocalFileRoute) Match(method string, path string) (bool, string) {
	ok, xtra := l.Path.Match(path)
	fs := FSAdapter.RootPath(filepath.Dir(l.LocalPath))
	if !ok || method != GET || xtra != "" {
		return false, path
	}
	fname := "/" + filepath.Base(l.LocalPath)
	if _, err := fs.Stat(fname); os.IsNotExist(err) {
		return false, path
	}
	return true, xtra
}

// Handle sets the response to return a local file.
func (l *LocalFileRoute) Handle(c Context) Response {
	fname := "/" + filepath.Base(l.LocalPath)
	fpath := filepath.Dir(l.LocalPath)
	return c.File(fpath, fname)
}

// RoutePath returns the route path for the route.
func (l *LocalFileRoute) RoutePath() RoutePath {
	return l.Path
}

// LocalPathRoute handles serving any file under a path. If the requested file
// exists it will be served if not the router will continue processing
// routes.
type LocalPathRoute struct {
	Path      RoutePath
	LocalPath string
}

// Match checks if a file exists under the local path
func (l *LocalPathRoute) Match(method string, path string) (bool, string) {
	fs := FSAdapter.RootPath(l.LocalPath)
	if len(path) < len(l.Path) {
		return false, path
	}
	fname := "/" + path[len(l.Path):]
	stat, err := fs.Stat(fname)
	if err != nil {
		return false, path
	}
	if stat.Mode().IsDir() {
		return false, path
	}
	return true, ""
}

// Handle sets the resposne up to serve the local file.
func (l *LocalPathRoute) Handle(c Context) Response {
	fname := c.ScopedPath[len(l.Path):]
	fpath := l.LocalPath
	return c.File(fpath, fname)
}

// RoutePath returns the routepath for the route
func (l *LocalPathRoute) RoutePath() RoutePath {
	return l.Path
}
