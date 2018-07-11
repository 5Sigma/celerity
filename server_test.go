package celerity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

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

func TestErrorResponse(t *testing.T) {
	server := New()
	server.Route(GET, "/foo", func(c Context) Response {
		return c.Error(412, errors.New("error message"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/foo")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	defer res.Body.Close()
	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response: %s", err.Error())
	}

	jsRes := struct {
		Error string   `json:"error" validate:"eq=error message"`
		Data  struct{} `json:"data" validate:"isdefault"`
	}{}

	err = json.Unmarshal(bbody, &jsRes)
	if err != nil {
		t.Errorf("%s : %s", err.Error(), string(bbody))
		return
	}

	if err := validator.New().Struct(jsRes); err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != 412 {
		t.Error("status code invalid")
	}
}

func TestFailResponse(t *testing.T) {
	{
		server := New()
		server.Route(GET, "/foo", func(c Context) Response {
			return c.Fail(errors.New("error message"))
		})

		ts := httptest.NewServer(server)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/foo")
		if err != nil {
			t.Errorf("Error requesting url: %s", err.Error())
		}

		defer res.Body.Close()
		bbody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error reading response: %s", err.Error())
		}

		jsRes := struct {
			Error string   `json:"error" validate:"eq=error message"`
			Data  struct{} `json:"data" validate:"isdefault"`
		}{}

		err = json.Unmarshal(bbody, &jsRes)
		if err != nil {
			t.Errorf("%s : %s", err.Error(), string(bbody))
			return
		}

		if err := validator.New().Struct(jsRes); err != nil {
			t.Error(err.Error())
		}

		if res.StatusCode != 500 {
			t.Error("status code invalid")
		}
	}
	{
		viper.Set("env", "prod")
		server := New()
		server.Route(GET, "/foo", func(c Context) Response {
			return c.Fail(errors.New("error message"))
		})

		ts := httptest.NewServer(server)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/foo")
		if err != nil {
			t.Errorf("Error requesting url: %s", err.Error())
		}

		defer res.Body.Close()
		bbody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error reading response: %s", err.Error())
		}

		jsRes := struct {
			Error string   `json:"error" validate:"eq=the request could not be processed"`
			Data  struct{} `json:"data" validate:"isdefault"`
		}{}

		err = json.Unmarshal(bbody, &jsRes)
		if err != nil {
			t.Errorf("%s : %s", err.Error(), string(bbody))
			return
		}

		if err := validator.New().Struct(jsRes); err != nil {
			t.Error(err.Error())
		}

		if res.StatusCode != 500 {
			t.Error("status code invalid")
		}
	}
}

func TestWildcardRouting(t *testing.T) {
	t.Run("valid path", func(t *testing.T) {
		svr := New()
		svr.GET("/get/*", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/get/something/cool", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode != 200 {
			t.Errorf("Non 200 response: %d", r.StatusCode)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		svr := New()
		svr.GET("/get/*", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/not/something/cool", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode == 200 {
			t.Error("200 response code")
		}
	})
}

func TestServerMethodAliases(t *testing.T) {
	{
		svr := New()
		svr.GET("/get", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/get", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode != 200 {
			t.Errorf("Non 200 response code for valid method/path: %d", r.StatusCode)
		}
	}
	{
		svr := New()
		svr.PUT("/put", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(PUT, "/put", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.DELETE("/delete", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(DELETE, "/delete", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.PATCH("/patch", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(PATCH, "/patch", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		svr := New()
		svr.POST("/post", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(POST, "/post", nil)
		c := RequestContext(req)
		r := svr.Router.Root.Handle(c)
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

func TestStaticPathServing(t *testing.T) {
	server := New()
	adapter := NewMEMAdapter()
	server.FSAdapter = adapter
	server.ServePath("/test", "/public/files")

	adapter.MEMFS.MkdirAll("/outsideroot", 0755)
	afero.WriteFile(adapter.MEMFS, "/outsideroot/test.txt", []byte("outside root file"), 0755)
	afero.WriteFile(adapter.MEMFS, "/public/files/test.txt", []byte("public file"), 0755)

	t.Run("get valid file", func(t *testing.T) {
		ts := httptest.NewServer(server)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/test/test.txt")
		if err != nil {
			t.Errorf("Error requesting url: %s", err.Error())
		}

		if v := res.Header.Get("Content-Type"); v != "application/octet-stream" {
			t.Errorf("content type was %s", v)
		}

		if s := res.StatusCode; s != 200 {
			t.Errorf("status code was %d", s)
		}

		defer res.Body.Close()
		bbody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error reading response: %s", err.Error())
		}
		if string(bbody) != "public file" {
			t.Errorf("body was: %s", string(bbody))
		}
	})

	t.Run("get invalid file", func(t *testing.T) {
		ts := httptest.NewServer(server)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/test/test2.txt")
		if err != nil {
			t.Errorf("Error requesting url: %s", err.Error())
		}

		if s := res.StatusCode; s != 404 {
			t.Errorf("status code was %d", s)
		}
	})

	t.Run("get unrooted file", func(t *testing.T) {
		ts := httptest.NewServer(server)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/../outsideroot/test.txt")
		if err != nil {
			t.Errorf("Error requesting url: %s", err.Error())
		}

		if s := res.StatusCode; s != 404 {
			t.Errorf("status code was %d", s)
		}
	})
}

func TestFileServing(t *testing.T) {
	server := New()
	adapter := NewMEMAdapter()
	server.FSAdapter = adapter
	server.ServePath("/test", "/public/files")

	adapter.MEMFS.MkdirAll("/outsideroot", 0755)
	afero.WriteFile(adapter.MEMFS, "/public/files/test.txt", []byte("public file"), 0755)

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test/test.txt")
	if err != nil {
		t.Errorf("Error requesting url: %s", err.Error())
	}

	if s := res.StatusCode; s != 200 {
		t.Errorf("status code was %d", s)
	}
	defer res.Body.Close()
	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response: %s", err.Error())
	}
	if string(bbody) != "public file" {
		t.Errorf("body was: %s", string(bbody))
	}
}
