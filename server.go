package celerity

import (
	"net/http"
)

// Server - Main server instance
type Server struct {
	Router          *Router
	ResponseAdapter ResponseAdapter
}

// NewServer - Initialize a new server
func NewServer() *Server {
	return &Server{
		ResponseAdapter: &JSONResponseAdapter{},
		Router:          NewRouter(),
	}
}

// ServeHTTP - Serves the HTTP request. Complies with http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := RequestContext(r)
	c.SetQueryParamsFromURL(r.URL)
	resp := s.Router.Handle(c, r)
	b, err := s.ResponseAdapter.Process(c, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, vs := range c.Response.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

// Pre - Register prehandle middleware for the root scope.
func (s *Server) Pre(mw MiddlewareHandler) {
	s.Router.Root.Pre(mw)
}

// Use - Use a middleware in the root scope.
func (s *Server) Use(mw MiddlewareHandler) {
	s.Router.Root.Use(mw)
}

// Start the server.
func (s *Server) Start(host string) error {
	return http.ListenAndServe(host, s)
}

// Scope creates a new scope from the root scope
func (s *Server) Scope(path string) *Scope {
	return s.Router.Root.Scope(path)
}

// Route - Set a route on the root scope.
func (s *Server) Route(method, path string, h RouteHandler) {
	s.Router.Root.Route(method, path, h)
}

// GET creates a route for a GET method.
func (s *Server) GET(path string, handler RouteHandler) Route {
	return s.Router.Root.Route(GET, path, handler)
}

// POST creates a route for a POST method.
func (s *Server) POST(path string, handler RouteHandler) Route {
	return s.Router.Root.Route(POST, path, handler)
}

// PUT creates a route for a PUT method.
func (s *Server) PUT(path string, handler RouteHandler) Route {
	return s.Router.Root.Route(PUT, path, handler)
}

// PATCH creates a route for a PATCH method.
func (s *Server) PATCH(path string, handler RouteHandler) Route {
	return s.Router.Root.Route(PATCH, path, handler)
}

// DELETE creates a route for a DELETE method.
func (s *Server) DELETE(path string, handler RouteHandler) Route {
	return s.Router.Root.Route(DELETE, path, handler)
}
