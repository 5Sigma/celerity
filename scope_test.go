package celerity

import (
	"fmt"
	"net/http"
	"testing"
)

func emptyHandler(c Context) Response {
	return c.Respond(nil)
}

func TestScopeRoute(t *testing.T) {
	{
		req, _ := http.NewRequest(GET, "/users", nil)
		scope := newScope("/")
		scope.Route(GET, "/users", emptyHandler)
		if !scope.Match(req, req.URL.Path) {
			t.Error("scope did not match valid path")
		}
	}
	{
		req, _ := http.NewRequest(GET, "/bad", nil)
		scope := newScope("/")
		scope.Route(GET, "/users", emptyHandler)
		if scope.Match(req, req.URL.Path) {
			t.Error("scope did match invalid path")
		}
	}
	{
		scope := newScope("/")
		sub := scope.Scope("/users")
		sub.Route(GET, "/:id", emptyHandler)
		{
			req, _ := http.NewRequest(GET, "/users/1", nil)
			if !scope.Match(req, req.URL.Path) {
				t.Error("scope did not match valid path")
			}
		}
		{
			req, _ := http.NewRequest(GET, "/users/1/abd", nil)
			if scope.Match(req, req.URL.Path) {
				t.Error("scope did match invalid path")
			}
		}
	}
}

func TestScopeHandle(t *testing.T) {
	{
		scope := newScope("/")
		scope.Route("GET", "/users", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(GET, "/users", nil)
		c.Request = req
		r := scope.Handle(c)
		v, ok := r.Data.(string)
		if !ok {
			t.Error("Data result incorrect type")
			return
		} else if v != "test" {
			t.Error("Data result not correct")
		}
	}
}

func TestMethodRouting(t *testing.T) {
	{
		scope := newScope("/")
		scope.Route(GET, "/users", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(GET, "/users", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		scope := newScope("/")
		scope.Route(POST, "/users", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(GET, "/users", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 404 {
			t.Error("Non 404 response code for invalid method/path")
		}
	}
}

func TestMethodAliases(t *testing.T) {
	{
		scope := newScope("/")
		scope.GET("/get", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(GET, "/get", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		scope := newScope("/")
		scope.PUT("/put", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(PUT, "/put", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		scope := newScope("/")
		scope.DELETE("/delete", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(DELETE, "/delete", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		scope := newScope("/")
		scope.PATCH("/patch", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(PATCH, "/patch", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
	{
		scope := newScope("/")
		scope.POST("/post", func(c Context) Response {
			return c.R("test")
		})
		c := NewContext()
		req, _ := http.NewRequest(POST, "/post", nil)
		c.Request = req
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	}
}

func BenchmarkScopeRoute(b *testing.B) {
	root := newScope("/")
	scope := root
	for n := 0; n < 10; n++ {
		scope = root.Scope(fmt.Sprintf("/%d", n))
	}
	scope.Route("GET", "ep", emptyHandler)
	for n := 0; n < b.N; n++ {
		req, _ := http.NewRequest(GET, "/0/1/2/3/4/5/6/7/8/9/ep", nil)
		scope.Match(req, req.URL.Path)
	}
}

func TestScopeCollision(t *testing.T) {
	root := newScope("/")
	ss1 := root.Scope("/")
	ss2 := root.Scope("/")
	users := ss1.Scope("/users")
	ss2.Scope("/users")
	users.Route(GET, "/get", func(c Context) Response {
		return c.R("test")
	})
	c := NewContext()
	req, _ := http.NewRequest(GET, "/users/get", nil)
	c.Request = req
	r := root.Handle(c)
	v, ok := r.Data.(string)
	if !ok {
		t.Error("Data result incorrect type")
		return
	} else if v != "test" {
		t.Error("Data result not correct")
	}
}

func TestMiddleware(t *testing.T) {
	root := newScope("/")
	root.Use(func(next RouteHandler) RouteHandler {
		return func(c Context) Response {
			c.Set("prop", "123")
			return next(c)
		}
	})
	root.Route(GET, "/middleware", func(c Context) Response {
		if v, ok := c.Get("prop").(string); ok {
			return c.R(v)
		}
		return c.R("bad")
	})
	c := NewContext()
	req, _ := http.NewRequest(GET, "/middleware", nil)
	c.Request = req
	r := root.Handle(c)
	v, ok := r.Data.(string)
	if !ok {
		t.Error("Data result incorrect type")
		return
	} else if v != "123" {
		t.Error("Data result not correct")
	}
}
