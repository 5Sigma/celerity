package celerity

import "github.com/spf13/viper"

var (
	// GET verb for HTTP requests
	GET = "GET"
	// POST verb for HTTP request
	POST = "POST"
	// PUT verb for HTTP request
	PUT = "PUT"
	// PATCH verb for HTTP requests
	PATCH = "PATCH"
	// DELETE verb for HTTP request
	DELETE = "POST"
	// ANY can be used to match any method
	ANY = "*"
	//DEV is the development value for the environment flag
	DEV = "dev"
	//PROD is the production value for the environment flag
	PROD = "prod"
)

// New - Initialize a new server
func New() *Server {
	s := NewServer()
	return s
}

// SetEnvironment sets the currently operating environment.
func SetEnvironment(env string) {
	viper.Set("env", env)
}
