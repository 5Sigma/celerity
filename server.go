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
	c := NewContext()
	c.Request = r
	c.SetQueryParamsFromURL(r.URL)
	resp := s.Router.Handle(c, r)
	b, err := s.ResponseAdapter.Process(c, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range c.Response.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

// Use - Use a middleware in the root scope.
func (s *Server) Use(mw MiddlewareHandler) {
	s.Router.Root.Use(mw)
}

// Start the server.
func (s *Server) Start(host string) error {
	return http.ListenAndServe(host, s)
}

// Route - Set a route on the root scope.
func (s *Server) Route(method, path string, h RouteHandler) {
	s.Router.Root.Route(method, path, h)
}
