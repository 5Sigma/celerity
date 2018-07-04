package celerity

import "strings"

// RoutePath - A path for a route or a group.
type RoutePath string

// Match - Matches the routepath aganst an incomming path
func (rp RoutePath) Match(path string) (bool, string) {
	if rp == "/" {
		return true, path
	}
	if rp[0] == '/' {
		rp = rp[1:]
	}
	if path == "" {
		path = "/"
	}
	if path[0] == '/' {
		path = path[1:]
	}
	pathTokens := strings.Split(path, "/")
	rpTokens := strings.Split(string(rp), "/")
	if len(rpTokens) > len(pathTokens) {
		return false, path
	}
	for idx, t := range rpTokens {
		if t == "*" {
			return true, ""
		}
		if t == "" {
			continue
		}
		if t[0] == ':' {
			continue
		}
		if t != pathTokens[idx] {
			return false, path
		}
	}
	return true, strings.Join(pathTokens[len(rpTokens):], "/")
}

// GetURLParams - Returns a map of url param/values based on the path given.
func (rp RoutePath) GetURLParams(path string) map[string]string {
	params := map[string]string{}
	pathTokens := strings.Split(path, "/")
	rpTokens := strings.Split(string(rp), "/")
	for idx, t := range rpTokens {
		if t == "" {
			continue
		}
		if t[0] == ':' {
			params[t[1:]] = pathTokens[idx]
		}
	}
	return params
}
