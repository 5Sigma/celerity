package middleware

import "github.com/5Sigma/celerity"

// AegisAdapter - Adapter interface for handling authentication
type AegisAdapter interface {
	ValidateSession(c celerity.Context, token string) bool
}

// AegisConfig - Description
type AegisConfig struct {
	Adapter AegisAdapter
}

//Aegis - Create an aegis middleware instance
func Aegis(config AegisConfig) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			sessionToken := c.Header().Get("Authorization")
			if !config.Adapter.ValidateSession(c, sessionToken) {
				return celerity.NewErrorResponse(401, "Session invalid")
			}
			return next(c)
		}
	}
}
