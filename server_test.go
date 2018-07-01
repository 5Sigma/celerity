package celerity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	validator "gopkg.in/go-playground/validator.v9"
)

func TestServeHTTP(t *testing.T) {
	server := New()
	server.Route("GET", "/foo", func(c Context) Response {
		return c.R(map[string]string{"test": "test"})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response: %s", err.Error())
	}
	defer res.Body.Close()

	jsRes := struct {
		RequestID string `json:"requestId" validate:"len=32"`
		Data      struct {
			Test string `json:"test" validate:"required,eq=test"`
		} `validate:"required" json:"data"`
		Error string `json:"eq="`
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

func TestURLParamHandling(t *testing.T) {
	server := New()
	server.Route("GET", "/foo/:id", func(c Context) Response {
		return c.R(map[string]interface{}{"id": c.URLParams.Int("id")})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo/13")
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
			ID int `json:"id" validate:"eq=13"`
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

func TestQueryParamHandling(t *testing.T) {
	server := New()
	server.Route("GET", "/foo", func(c Context) Response {
		return c.R(map[string]interface{}{"name": c.QueryParams.String("name")})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo?name=alice")
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
			Name string `json:"name" validate:"eq=alice"`
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

func TestNotFound(t *testing.T) {
	server := New()
	server.Route("GET", "/foo/:id", func(c Context) Response {
		return c.R(map[string]interface{}{"id": c.URLParams.Int("id")})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/bad/path")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	defer res.Body.Close()
	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response: %s", err.Error())
	}

	jsRes := struct {
		Error string      `json:"error" validate:"eq=The requested resource was not found"`
		Data  interface{} `json:"data" validate:"isdefault"`
	}{}

	err = json.Unmarshal(bbody, &jsRes)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if err := validator.New().Struct(jsRes); err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != 404 {
		t.Errorf("Status code should be 404 was %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestRewrite(t *testing.T) {
	server := New()
	server.Route("GET", "/users/:id/profile", func(c Context) Response {
		return c.R(map[string]interface{}{"id": c.URLParams.Int("id")})
	})

	server.Rewrite(RewriteRules{
		"/people/(.*)/profile": "/users/$1/profile",
	})

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

func TestDataExtration(t *testing.T) {
	server := New()
	server.Route(POST, "/foo", func(c Context) Response {
		req := struct {
			Name string `json:"name"`
		}{}
		if err := c.Extract(&req); err != nil {
			return c.Fail(err)
		}
		return c.R(map[string]interface{}{"name": req.Name})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	payload := []byte(`{ "name": "alice" }`)

	res, err := http.Post(ts.URL+"/foo", "application/json", bytes.NewBuffer(payload))
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
			Name string `json:"name" validate:"eq=alice"`
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

func TestServerMethodAliases(t *testing.T) {
	{
		svr := New()
		svr.GET("/get", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(GET, "/get", nil)
		r := svr.Router.Root.Handle(c, req)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.PUT("/put", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(PUT, "/put", nil)
		r := svr.Router.Root.Handle(c, req)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.DELETE("/delete", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(DELETE, "/delete", nil)
		r := svr.Router.Root.Handle(c, req)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.PATCH("/patch", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(PATCH, "/patch", nil)
		r := svr.Router.Root.Handle(c, req)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.POST("/post", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(POST, "/post", nil)
		r := svr.Router.Root.Handle(c, req)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
}

func BenchmarkRouteProcessing(b *testing.B) {
	server := New()
	server.Route("GET", "/foo/:id", func(c Context) Response {
		return c.R(map[string]interface{}{"id": c.URLParams.Int("id")})
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	for n := 0; n < b.N; n++ {
		res, _ := http.Get(ts.URL + "/foo/13")
		res.Body.Close()
	}
}

func BenchmarkMiddleware(b *testing.B) {
	server := New()
	for i := 0; i < 50; i++ {
		server.Router.Root.Use(func(next RouteHandler) RouteHandler {
			return func(c Context) Response {
				c.Set(fmt.Sprintf("prop%d", i), "123")
				return next(c)
			}
		})

	}

	ts := httptest.NewServer(server)
	defer ts.Close()

	for n := 0; n < b.N; n++ {
		res, _ := http.Get(ts.URL + "/foo/13")
		res.Body.Close()
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
