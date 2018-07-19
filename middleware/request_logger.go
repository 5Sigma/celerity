package middleware

import (
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/5Sigma/celerity"
	"github.com/5Sigma/vox"
)

// RequestLoggerConfig configures the request logger middleware and can be
// passed to RequestLoggerWithConfig.
type RequestLoggerConfig struct {
	ConsoleTemplate string
	Log             *vox.Vox
}

// RequestLoggerData is used to format the output for console and file logging
// operations.
type RequestLoggerData struct {
	Now      time.Time
	URL      *url.URL
	Duration time.Duration
	C        *celerity.Context
	R        *celerity.Response
}

// NewLoggerConfig creates a default configuration for RequestLoggerConfig.
func NewLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		ConsoleTemplate: `
{{ .C.RequestID }} - [{{.Now.Format "1/2/2006 15:04:05"}}] - {{ .C.Request.Method }} {{.URL.Path}} ({{ .Duration }}) - {{ .R.StatusCode }}  {{ .R.StatusText }}
{{ if eq .R.Success false -}}
	ERROR: {{ .R.Error }}
{{- end }}
`,
	}
}

// RequestLogger creates a new logger middleware with sane defaults.
func RequestLogger() celerity.MiddlewareHandler {
	return RequestLoggerWithConfig(NewLoggerConfig())
}

// RequestLoggerWithConfig creates a new Request logger with the given config.
func RequestLoggerWithConfig(rlc RequestLoggerConfig) celerity.MiddlewareHandler {
	if rlc.Log == nil {
		rlc.Log = vox.New()
	}
	tmpl, err := template.New("console").Parse(strings.TrimSpace(rlc.ConsoleTemplate) + "\n")
	if err != nil {
		return func(next celerity.RouteHandler) celerity.RouteHandler {
			return func(c celerity.Context) celerity.Response {
				return next(c)
			}
		}
	}
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			start := time.Now()
			r := next(c)
			reqTime := time.Since(start)
			data := &RequestLoggerData{
				URL:      c.Request.URL,
				Now:      time.Now(),
				Duration: reqTime,
				C:        &c,
				R:        &r,
			}
			tmpl.Execute(rlc.Log, data)
			return r
		}
	}
}
