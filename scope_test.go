package celerity

import (
	"fmt"
	"testing"
)

func emptyHandler(c Context) Response {
	return c.Respond(nil)
}

func TestScopeRoute(t *testing.T) {
	{
		scope := newScope("/")
		scope.Route(GET, "/users", emptyHandler)
		if !scope.Match("/users") {
			t.Error("scope did not match valid path")
		}
		if scope.Match("/bad") {
			t.Error("scope did match invalid path")
		}
	}
	{
		scope := newScope("/")
		sub := scope.Scope("/users")
		sub.Route(GET, "/:id", emptyHandler)
		if !scope.Match("/users/1") {
			t.Error("scope did not match valid path")
		}
		if scope.Match("/bad") {
			t.Error("scope did match invalid path")
		}
		if scope.Match("/users/1/bad") {
			t.Error("scope did match invalid path")
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
		r := scope.Handle(c, "/users")
		v, ok := r.Data.(string)
		if !ok {
			t.Error("Data result incorrect type")
			return
		} else if v != "test" {
			t.Error("Data result not correct")
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
		scope.Match("/0/1/2/3/4/5/6/7/8/9/ep")
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
	r := root.Handle(c, "/users/get")
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
	r := root.Handle(c, "/middleware")
	v, ok := r.Data.(string)
	if !ok {
		t.Error("Data result incorrect type")
		return
	} else if v != "123" {
		t.Error("Data result not correct")
	}
}
