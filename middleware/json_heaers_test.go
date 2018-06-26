package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5Sigma/celerity"
)

func TestJSONHeaders(t *testing.T) {
	server := celerity.New()
	server.GET("/foo", func(c celerity.Context) celerity.Response {
		return c.R(nil)
	})

	server.Router.Root.Use(JSONMiddleware())

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content type header not set")
	}

	res.Body.Close()
}
