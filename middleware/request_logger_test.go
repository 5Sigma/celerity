package middleware

import (
	"bytes"
	"testing"

	"github.com/5Sigma/celerity"
	"github.com/5Sigma/celerity/celeritytest"
)

func TestConsoleOutput(t *testing.T) {
	svr := celerity.New()
	svr.GET("/foo", func(c celerity.Context) celerity.Response {
		return c.R(nil)
	})

	out := new(bytes.Buffer)
	config := NewLoggerConfig()
	config.ConsoleOut = out

	svr.Use(RequestLoggerWithConfig(config))

	celeritytest.Get(svr, "/foo")

	if out.Len() != 39 {
		t.Errorf("output length was not correct: %d, should be 39", out.Len())
	}
}
