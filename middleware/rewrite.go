package middleware

import (
	"regexp"
	"strings"

	"github.com/5Sigma/celerity"
)

// RewriteRules is a set of rules used to rewrite and transform the incoming
// url. See the Server.Rewrite function.
type RewriteRules map[string]string

// Match check if a path matches any of the rules in the ruleset. If it does it
// returns the transformed URL.
func (rr RewriteRules) Match(path string) (bool, string) {
	for k, v := range rr {
		re, err := regexp.Compile(k)
		if err != nil {
			continue
		}
		res := re.FindStringSubmatch(path)
		if len(res) > 0 {
			for _, s := range res[1:] {
				v = strings.Replace(v, "$1", s, -1)
			}
			return true, v
		}
	}
	return false, ""
}

// Rewrite returns a rewrite middleware with the given rewrite rules
func Rewrite(rules RewriteRules) celerity.MiddlewareHandler {
	return func(next celerity.RouteHandler) celerity.RouteHandler {
		return func(c celerity.Context) celerity.Response {
			rewrite, rewritePath := rules.Match(c.ScopedPath)
			if rewrite {
				c.ScopedPath = rewritePath
			}
			return next(c)
		}
	}
}
