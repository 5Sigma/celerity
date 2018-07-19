package middleware

import (
	"io"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/5Sigma/celerity"
)

// RequestLoggerConfig configures the request logger middleware and can be
// passed to RequestLoggerWithConfig.
type RequestLoggerConfig struct {
	ConsoleOut      io.Writer
	ConsoleTemplate string
}

// RequestLoggerData is used to format the output for console and file logging
// operations.
type RequestLoggerData struct {
	Now         time.Time
	URL         *url.URL
	RequestTime time.Duration
	Context     *celerity.Context
	Response    *celerity.Response
}

// NewLoggerConfig creates a default configuration for RequestLoggerConfig.
func NewLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		ConsoleOut: os.Stdout,
		ConsoleTemplate: `
		{{ .Context.RequestID }} - [{{.Now.Format "1/2/2006 15:04:05"}}] - {{ .Context.Request.Method }} {{.URL.Path}} ({{ .RequestTime }}) - {{ .Response.StatusCode }}  {{ .Response.StatusText }}
{{ if eq .Response.Success false -}}
	ERROR: {{ .Response.Error }}
{{- end }}
`,
	}
}

// ConsoleOutput generates the console line output for the request
func (rlc *RequestLoggerConfig) ConsoleOutput(ctx celerity.Context, r celerity.Response) {
}

// RequestLogger creates a new logger middleware with sane defaults.
func RequestLogger() celerity.MiddlewareHandler {
	return RequestLoggerWithConfig(NewLoggerConfig())
}

// RequestLoggerWithConfig creates a new Request logger with the given config.
func RequestLoggerWithConfig(rlc RequestLoggerConfig) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			start := time.Now()
			r := next(c)
			reqTime := time.Since(start)

			data := &RequestLoggerData{
				URL:         c.Request.URL,
				Now:         time.Now(),
				RequestTime: reqTime,
				Context:     &c,
				Response:    &r,
			}

			t, err := template.New("console").Parse(strings.TrimSpace(rlc.ConsoleTemplate) + "\n")
			if err != nil {
				return r
			}
			t.Execute(c.Log, data)

			return r
		}
	}
}
