package celerity

import (
	"net/http"
)

// Scope - A group of routes and subgroups used to represent the routing
// structure for the serve.r
type Scope struct {
	Path       RoutePath
	Scopes     []*Scope
	Routes     []Route
	Middleware []MiddlewareHandler
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

func (s *Scope) handleWithMiddleware(c Context, req *http.Request, path string, middleware []MiddlewareHandler) Response {
	ok, rPath := s.Path.Match(path)
	if !ok {
		return NewErrorResponse(http.StatusNotFound, "The requested resource was not found")
	}

	middleware = append(middleware, s.Middleware...)

	for _, r := range s.Routes {
		if r.Match(req.Method, rPath) {
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
		if ss.Match(req, rPath) {
			return ss.handleWithMiddleware(c, req, rPath, middleware)
		}
	}
	return NewErrorResponse(http.StatusNotFound, "The requested resource was not found")
}

// Handle - Handle an incomming URL
func (s *Scope) Handle(c Context, req *http.Request) Response {
	return s.handleWithMiddleware(c, req, req.URL.Path, []MiddlewareHandler{})
}
