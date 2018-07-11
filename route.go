package celerity

// Route - A basic route in the server.
type Route struct {
	Path       RoutePath
	Method     string
	Handler    RouteHandler
	Middleware []MiddlewareHandler
}

// RouteHandler - The handler function that gets called when a route is invoked.
type RouteHandler func(Context) Response

// Match - Matches the routes path against the incomming url
func (r *Route) Match(method, path string) (bool, string) {
	ok, xtra := r.Path.Match(path)
	return (ok && method == r.Method && xtra == ""), xtra
}
