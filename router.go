package celerity

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//MiddlewareHandler is a function that can be used in scopes and
//routes to transform the context before the route is processed.
type MiddlewareHandler func(RouteHandler) RouteHandler

// Router - The server router stores all routes, groups, and determines what to
// call when a given path is invoked.
type Router struct {
	Root *Scope
}

// NewRouter - Initailize a new router
func NewRouter() *Router {
	return &Router{
		Root: newScope("/"),
	}
}

// Handle - Process the inccomming URL and execute the first matching route
func (r *Router) Handle(c Context, req *http.Request) Response {
	return r.Root.Handle(c)
}

// RoutePath - A path for a route or a group.
type RoutePath string

// Match - Matches the routepath aganst an incomming path
func (rp RoutePath) Match(path string) (bool, string) {
	if rp == "/" {
		return true, path
	}
	if rp[0] == '/' {
		rp = rp[1:]
	}
	if path == "" {
		path = "/"
	}
	if path[0] == '/' {
		path = path[1:]
	}
	pathTokens := strings.Split(path, "/")
	rpTokens := strings.Split(string(rp), "/")
	if len(rpTokens) > len(pathTokens) {
		return false, path
	}
	for idx, t := range rpTokens {
		if t == "*" {
			return true, ""
		}
		if t == "" {
			continue
		}
		if t[0] == ':' {
			continue
		}
		if t != pathTokens[idx] {
			return false, path
		}
	}
	return true, strings.Join(pathTokens[len(rpTokens):], "/")
}

// GetURLParams - Returns a map of url param/values based on the path given.
func (rp RoutePath) GetURLParams(path string) map[string]string {
	println(path)
	if path[0] != '/' {
		path = "/" + path
	}
	params := map[string]string{}
	pathTokens := strings.Split(path, "/")
	rpTokens := strings.Split(string(rp), "/")
	for idx, t := range rpTokens {
		if t == "" {
			continue
		}
		if t[0] == ':' {
			params[t[1:]] = pathTokens[idx]
		}
	}
	return params
}

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

	serveFile(c.Writer, fpath, fname)
	return Response{Handled: true}
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
	serveFile(c.Writer, fpath, fname)
	return Response{Handled: true}
}

// RoutePath returns the routepath for the route
func (l *LocalPathRoute) RoutePath() RoutePath {
	return l.Path
}

func serveFile(w http.ResponseWriter, froot, fpath string) {
	fs := FSAdapter.RootPath(froot)
	f, err := fs.Open(fpath)
	if os.IsNotExist(err) {
		w.WriteHeader(404)
		w.Write([]byte("The file does not exists"))
		return

	}
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	fileHeader := make([]byte, 512)
	f.Read(fileHeader)
	fstat, err := f.Stat()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	fsize := strconv.FormatInt(fstat.Size(), 10)
	fname := filepath.Base(fpath)
	contentType := http.DetectContentType(fileHeader)
	switch filepath.Ext(fname) {
	case ".css":
		contentType = "text/css"
	}

	// w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fsize)
	f.Seek(0, 0)
	io.Copy(w, f)
}
