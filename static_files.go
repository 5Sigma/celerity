package celerity

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

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
