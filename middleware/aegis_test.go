package middleware

import (
	"net/http"
	"testing"

	"github.com/5Sigma/celerity"
	"github.com/5Sigma/celerity/celeritytest"
)

type MockAegisAdapter struct {
}

// ValidateSession - Session validation
func (a *MockAegisAdapter) ValidateSession(c celerity.Context, token string) bool {
	return token == "123"
}

func TestSessionValidation(t *testing.T) {
	server := celerity.New()

	server.Router.Route("GET", "/foo", func(c celerity.Context) celerity.Response {
		return c.R(nil)
	})

	aeAdapter := &MockAegisAdapter{}
	aeConfig := AegisConfig{Adapter: aeAdapter}
	server.Use(Aegis(aeConfig))
	{
		res, _ := celeritytest.Get(server, "/foo")

		if res.StatusCode != 401 {
			t.Errorf("Status was %d", res.StatusCode)
		}
		if ok, _ := res.AssertBool("success", false); !ok {
			t.Error("Success flag not set to false")
		}
	}
	{
		opts := celeritytest.RequestOptions{
			Path: "/foo",
			Header: http.Header{
				"Authorization": []string{"123"},
			},
			Method: "GET",
		}
		res, _ := celeritytest.Request(server, opts)
		if res.StatusCode != 200 {
			t.Errorf("Status was %d", res.StatusCode)
		}
	}

}
