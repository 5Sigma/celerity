package middleware

import (
	"bytes"
	"testing"

	"github.com/5Sigma/celerity"
)

func TestConsoleOutput(t *testing.T) {
	svr := celerity.New()
	svr.GET("/foo", func(c celerity.Context) celerity.Response {

	})

	out := new(bytes.Buffer)
	if err != nil {
		return err
	}

	config := NewLoggerConfig()
	config.ConsoleOut = out
}
