package middleware

import (
	"testing"

	"github.com/5Sigma/celerity"
	"github.com/5Sigma/celerity/celeritytest"
	"github.com/5Sigma/vox"
)

func TestConsoleOutput(t *testing.T) {
	svr := celerity.New()
	svr.GET("/foo", func(c celerity.Context) celerity.Response {
		return c.R(nil)
	})

	config := NewLoggerConfig()
	v := vox.New()
	pipeline := v.Test()
	config.Log = v

	svr.Use(RequestLoggerWithConfig(config))

	celeritytest.Get(svr, "/foo")

	if len(pipeline.All()) < 54 {
		t.Errorf("output length was not correct: %d, should be > 54\n%s",
			len(pipeline.Last()), pipeline.Last())
	}
}
