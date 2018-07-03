package celerity

import (
	"net/http"
)

// Scope - A group of routes and subgroups used to represent the routing
// structure for the serve.r
type Scope struct {
	Path          RoutePath
	Scopes        []*Scope
	Routes        []Route
	Middleware    []MiddlewareHandler
	PreMiddleware []MiddlewareHandler
}

// NewScope - Initializes a new scope
func newScope(path string) *Scope {
	return &Scope{
		Path:   RoutePath(path),
		Scopes: []*Scope{},
		Routes: []Route{},
	}
}

// Scope - Creates a new sub scope
func (s *Scope) Scope(path string) *Scope {
	ss := newScope(path)
	s.Scopes = append(s.Scopes, ss)
	return ss
}

// GET creates a route for a GET method.
func (s *Scope) GET(path string, handler RouteHandler) Route {
	return s.Route(GET, path, handler)
}

// POST creates a route for a POST method.
func (s *Scope) POST(path string, handler RouteHandler) Route {
	return s.Route(POST, path, handler)
}

// PUT creates a route for a PUT method.
func (s *Scope) PUT(path string, handler RouteHandler) Route {
	return s.Route(PUT, path, handler)
}

// PATCH creates a route for a PATCH method.
func (s *Scope) PATCH(path string, handler RouteHandler) Route {
	return s.Route(PATCH, path, handler)
}

// DELETE creates a route for a DELETE method.
func (s *Scope) DELETE(path string, handler RouteHandler) Route {
	return s.Route(DELETE, path, handler)
}

// Route - Create a new route within the scope
func (s *Scope) Route(method, path string, handler RouteHandler) Route {
	r := Route{
		Path:    RoutePath(path),
		Method:  method,
		Handler: handler,
	}
	s.Routes = append(s.Routes, r)
	return r
}

// Use - Use a middleware function
func (s *Scope) Use(mf ...MiddlewareHandler) {
	s.Middleware = append(s.Middleware, mf...)
}

// Pre - Registers middleware to be executed when the scope is first matched.
// This happens before the router searches for routes.
func (s *Scope) Pre(mf ...MiddlewareHandler) {
	s.PreMiddleware = append(s.PreMiddleware, mf...)
}

// Match - Check if the scope can handle the incomming url
func (s *Scope) Match(req *http.Request, path string) bool {
	ok, rPath := s.Path.Match(path)
	if !ok {
		return false
	}
	for _, r := range s.Routes {
		if r.Match(req.Method, rPath) {
			return true
		}
	}
	for _, ss := range s.Scopes {
		if ss.Match(req, rPath) {
			return true
		}
	}
	return false
}

func notFoundHandler(c Context) Response {
	return NewErrorResponse(http.StatusNotFound, "The requested resource was not found")
}

func (s *Scope) handleWithMiddleware(c Context, path string, middleware []MiddlewareHandler) Response {
	ok, rPath := s.Path.Match(path)

	middleware = append(middleware, s.Middleware...)

	if !ok {
		var h RouteHandler
		h = notFoundHandler
		for i := len(s.Middleware); i > 0; i-- {
			h = s.Middleware[i-1](h)
		}
		return h(c)
	}

	ph := func(c Context) Response {

		for _, r := range s.Routes {
			if r.Match(c.Request.Method, rPath) {
				c.SetParams(r.Path.GetURLParams(path))
				var h RouteHandler
				h = r.Handler
				for i := len(s.Middleware); i > 0; i-- {
					h = s.Middleware[i-1](h)
				}
				return h(c)
			}
		}
		for _, ss := range s.Scopes {
			if ss.Match(c.Request, rPath) {
				return ss.handleWithMiddleware(c, rPath, middleware)
			}
		}

		var h RouteHandler
		h = notFoundHandler
		for i := len(s.Middleware); i > 0; i-- {
			h = s.Middleware[i-1](h)
		}
		return h(c)
	}

	for i := len(s.PreMiddleware) - 1; i >= 0; i-- {
		ph = s.PreMiddleware[i](ph)
	}

	return ph(c)

}

// Handle - Handle an incomming URL
func (s *Scope) Handle(c Context) Response {
	return s.handleWithMiddleware(c, c.Request.URL.Path, []MiddlewareHandler{})
}
