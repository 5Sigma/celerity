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
	ExposeHeaders    []string
	Age              int
}

// CORSWithConfig returns a CORS middleware handler with a specific
// configuration.
func CORSWithConfig(config CORSConfig) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {

			origins := strings.Join(config.AllowOrigins, ",")
			methods := strings.Join(config.AllowMethods, ",")
			headers := strings.Join(config.AllowHeaders, ",")
			eHeaders := strings.Join(config.ExposeHeaders, ",")
			c.Response.Header.Set("Access-Control-Allow-Origin", origins)
			if methods != "" {
				c.Response.Header.Set("Access-Control-Allow-Methods", methods)
			}
			if headers != "" {
				c.Response.Header.Set("Access-Control-Allow-Headers", headers)
			}
			if config.Age > 0 {
				c.Response.Header.Set("Access-Control-Max-Age", string(config.Age))
			}
			if config.AllowCredentials {
				c.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			}
			if eHeaders != "" {
				c.Response.Header.Set("Access-Control-Expose-Headers", eHeaders)
			}

			if c.Request.Method == celerity.OPTIONS {
				return c.Response
			}
			return next(c)
		}
	}
}

// CORS returns a CORS middleware with sane defaults
//
// CORS middleware can be invoked with a configuration using CORSWithConfig.
// This middleware will add CORS headers to all requests matched by the scope.
// This middleware must be used as Preroute middleware using the Scope.Pre
// function, instead of Scope.Use.
func CORS() celerity.MiddlewareHandler {
	config := CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}
	return CORSWithConfig(config)
}
