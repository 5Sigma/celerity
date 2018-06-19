package celerity

var (
	// GET - GET verb for HTTP requests
	GET = "GET"
	// POST - POST verb for HTTP request
	POST = "POST"
	// PUT - PUT verb for HTTP request
	PUT = "POST"
	// DELETE - DELETE verb for HTTP request
	DELETE = "POST"
)

//MiddlewareHandler - A middleware function that can be used in scopes and
//routes
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
func (r *Router) Handle(c Context, path string) Response {
	return r.Root.Handle(c, path)
}

// Route - Create a route on the root scope
func (r *Router) Route(method, path string, handler RouteHandler) {
	r.Root.Route(method, path, handler)
}
