package celerity

import (
	"net/http"
	"regexp"
	"strings"
)

// Server - Main server instance
type Server struct {
	Router          *Router
	ResponseAdapter ResponseAdapter
	RewriteRules    RewriteRules
}

// NewServer - Initialize a new server
func NewServer() *Server {
	return &Server{
		ResponseAdapter: &JSONResponseAdapter{},
		Router:          NewRouter(),
		RewriteRules:    RewriteRules{},
	}
}

// ServeHTTP - Serves the HTTP request. Complies with http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext()
	rewrite, rewritePath := s.RewriteRules.Match(r.URL.Path)
	if rewrite {
		r.URL.Path = rewritePath
	}
	c.Request = r
	c.SetQueryParamsFromURL(r.URL)
	resp := s.Router.Handle(c, r.URL.Path)
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

/*
Rewrite rewrites urls based on on ruleset. The rewrite function accepts a
RewriteRule map that shold contain matching patterns and transformed URLs. The
matching patterns use RegEx to match incomming URLs and replace the URL using
the map value. Regex capture groups are provided to the tranformed url with
$1,$2, etc.

This function is also additive so subsequent calls will be added to the rule
set.

		svr := celerity.Server.New()
		svr.Rewrite(celerity.RewriteRules{
			"/people/(.*)/profile": "/users/$1/profile",
		})
*/
func (s *Server) Rewrite(rules RewriteRules) {
	for k, v := range rules {
		s.RewriteRules[k] = v
	}
}

// RewriteRules is a set of rules used to rewrite and transform the incoming
// url. See the Server.Rewrite function.
type RewriteRules map[string]string

// Match check if a path matches any of the rules in the ruleset. If it does it
// returns the transformed URL.
func (rr RewriteRules) Match(path string) (bool, string) {
	for k, v := range rr {
		re, err := regexp.Compile(k)
		if err != nil {
			continue
		}
		res := re.FindStringSubmatch(path)
		if len(res) > 0 {
			for _, s := range res[1:] {
				v = strings.Replace(v, "$1", s, -1)
			}
			return true, v
		}
	}
	return false, ""
}
