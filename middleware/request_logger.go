package middleware

import (
	"html/template"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/5Sigma/celerity"
)

// RequestLoggerConfig configures the request logger middleware and can be
// passed to RequestLoggerWithConfig.
type RequestLoggerConfig struct {
	LogToFiles      bool
	ConsoleOut      io.Writer
	ConsoleTemplate template.Template
}

// RequestLoggerData is used to format the output for console and file logging
// operations.
type RequestLoggerData struct {
	Now         time.Time
	URL         *url.URL
	RequestTime time.Duration
}

// NewRequestLoggerData creates and hydraites a RequestLoggerData structure to
// be used to format the output line for console and file logs.
func NewRequestLoggerData(ctx celerity.Context) RequestLoggerData {
	return RequestLoggerData{
		Now: time.Now(),
		URL: ctx.Request.URL,
	}
}

// NewLoggerConfig creates a default configuration for RequestLoggerConfig.
func NewLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		LogToFiles: false,
		ConsoleOut: os.Stdout,
	}
}

// ConsoleOutput generates the console line output for the request
func (rlc *RequestLoggerConfig) ConsoleOutput(ctx celerity.Context) {
	rlc.ConsoleTemplate.Execute(rlc.ConsoleOut, NewRequestLoggerData(ctx))
}

// RequestLogger creates a new logger middleware with sane defaults.
func RequestLogger() celerity.MiddlewareHandler {
	return RequestLoggerWithConfig(NewLoggerConfig())
}

// RequestLoggerWithConfig creates a new Request logger with the given config.
func RequestLoggerWithConfig(rlc RequestLoggerConfig) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			r := next(c)
			rlc.ConsoleOutput(c)
			return r
		}
	}
}
