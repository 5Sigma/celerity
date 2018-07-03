package middleware

import (
	"strings"

	"github.com/5Sigma/celerity"
)

// CORSConfig configures the cors middleware and can be passed to
// CORSWithConfig.
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	Age              int
}

// CORSWithConfig returns a CORS middleware handler with a specific
// configuration.
func CORSWithConfig(config CORSConfig) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			if c.Request.Method != celerity.OPTIONS {
				return next(c)
			}

			origins := strings.Join(config.AllowOrigins, ",")
			methods := strings.Join(config.AllowMethods, ",")
			headers := strings.Join(config.AllowHeaders, ",")
			c.Response.Header.Set("Access-Control-Allow-Origin", origins)
			if methods != "" {
				c.Response.Header.Set("Access-Control-Allow-Methods", methods)
			}
			if headers != "" {
				c.Response.Header.Set("Access-Control-Allow-Headers", methods)
			}
			return next(c)
		}
	}
}

// CORS returns a CORS middleware with sane defaults
func CORS() celerity.MiddlewareHandler {
	config := CORSConfig{
		AllowOrigins: []string{"*"},
	}
	return CORSWithConfig(config)
}
