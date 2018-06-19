package middleware

import "github.com/5Sigma/celerity"

// JSONMiddleware - Adds content type json headers
func JSONMiddleware() celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			c.Response.Headers["Content-Type"] = "application/json"
			return next(c)
		}
	}
}
