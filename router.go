package celerity

import "net/http"

var (
	// GET verb for HTTP requests
	GET = "GET"
	// POST verb for HTTP request
	POST = "POST"
	// PUT verb for HTTP request
	PUT = "POST"
	// DELETE verb for HTTP request
	DELETE = "POST"
	// ANY can be used to match any method
	ANY = "*"
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
	return r.Root.Handle(c, req)
}

// Route - Create a route on the root scope
func (r *Router) Route(method, path string, handler RouteHandler) {
	r.Root.Route(method, path, handler)
}
