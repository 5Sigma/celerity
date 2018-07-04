package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5Sigma/celerity"

	validator "gopkg.in/go-playground/validator.v9"
)

func TestRewrite(t *testing.T) {
	server := celerity.New()
	server.Route("GET", "/users/:id/profile", func(c celerity.Context) celerity.Response {
		return c.R(map[string]interface{}{"id": c.URLParams.Int("id")})
	})

	server.Pre(Rewrite(RewriteRules{
		"/people/(.*)/profile": "/users/$1/profile",
	}))

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/people/3/profile")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	defer res.Body.Close()
	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response: %s", err.Error())
	}

	jsRes := struct {
		Data struct {
			ID int `json:"id" validate:"eq=3"`
		}
	}{}

	err = json.Unmarshal(bbody, &jsRes)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := validator.New().Struct(jsRes); err != nil {
		t.Error(err.Error())
	}
}

func TestRewriteRulesMatch(t *testing.T) {
	rr := RewriteRules{
		"/people/(.*)/profile": "/users/$1/profile",
	}
	{
		ok, newPath := rr.Match("/people/3/profile")
		if !ok {
			t.Error("did not match path")
		}
		if newPath != "/users/3/profile" {
			t.Errorf("transformed path not correct: %s", newPath)
		}
	}
	{
		ok, _ := rr.Match("/peoples/3/profile")
		if ok {
			t.Error("should not match bad path")
		}
	}
}
