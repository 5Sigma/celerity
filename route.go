package celerity

import (
	"fmt"
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
