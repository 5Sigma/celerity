package middleware

import (
	"testing"

	"github.com/5Sigma/celerity"
	"github.com/5Sigma/celerity/celeritytest"
)

func TestCORS(t *testing.T) {
	svr := celerity.New()
	svr.GET("/foo", func(c celerity.Context) celerity.Response {
		return c.R(nil)
	})
	svr.Pre(CORS())

	reqOpts := celeritytest.RequestOptions{
		Path:   "/foo",
		Method: celerity.OPTIONS,
	}
	resp, _ := celeritytest.Request(svr, reqOpts)
	origins := resp.Header.Get("Access-Control-Allow-Origin")
	if len(origins) == 0 {
		t.Error("allow origins not set")
	}

	if origins != "*" {
		t.Errorf("allow origins should be * by default: %s", origins)
	}
}
