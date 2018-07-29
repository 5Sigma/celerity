package celerity

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/spf13/afero"
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
		req, _ := http.NewRequest(GET, "/users", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		v, ok := r.Data.(string)
		if !ok {
			t.Error("Data result incorrect type")
			return
		} else if v != "test" {
			t.Error("Data result not correct")
		}
	}
	{
		scope := newScope("/")
		scope.Route("GET", "/users", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/notfound", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 404 {
			t.Error("status code should be 404 for not found path")
		}
	}
}

func TestMethodRouting(t *testing.T) {
	{
		scope := newScope("/")
		scope.Route(GET, "/users", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/users", nil)
		c := RequestContext(req)
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
		req, _ := http.NewRequest(GET, "/users", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 404 {
			t.Error("Non 404 response code for invalid method/path")
		}
	}
}

func TestMethodAliases(t *testing.T) {

	t.Run("GET", func(t *testing.T) {
		scope := newScope("/")
		scope.GET("/get", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(GET, "/get", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	})

	t.Run("PUT", func(t *testing.T) {
		scope := newScope("/")
		scope.PUT("/put", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(PUT, "/put", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	})
	t.Run("DELETE", func(t *testing.T) {
		scope := newScope("/")
		scope.DELETE("/delete", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(DELETE, "/delete", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	})
	t.Run("PATCH", func(t *testing.T) {
		scope := newScope("/")
		scope.PATCH("/patch", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(PATCH, "/patch", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	})
	t.Run("POST", func(t *testing.T) {
		scope := newScope("/")
		scope.POST("/post", func(c Context) Response {
			return c.R("test")
		})
		req, _ := http.NewRequest(POST, "/post", nil)
		c := RequestContext(req)
		r := scope.Handle(c)
		if r.StatusCode != 200 {
			t.Error("Non 200 response code for valid method/path")
		}
	})
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
	req, _ := http.NewRequest(GET, "/users/get", nil)
	c := RequestContext(req)
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
	req, _ := http.NewRequest(GET, "/middleware", nil)
	c := RequestContext(req)
	r := root.Handle(c)
	v, ok := r.Data.(string)
	if !ok {
		t.Error("Data result incorrect type")
		return
	} else if v != "123" {
		t.Error("Data result not correct")
	}
}

func TestPre(t *testing.T) {
	s := newScope("/")
	mw := func(next RouteHandler) RouteHandler {
		return func(c Context) Response {
			return next(c)
		}
	}
	s.Pre(mw)
	if len(s.PreMiddleware) != 1 {
		t.Error("middleware not added to collection")
	}
}

func TestUse(t *testing.T) {
	s := newScope("/")
	mw := func(next RouteHandler) RouteHandler {
		return func(c Context) Response {
			return next(c)
		}
	}
	s.Use(mw)
	if len(s.Middleware) != 1 {
		t.Error("middleware not added to collection")
	}
}

func TestPanicRecovery(t *testing.T) {
	s := newScope("/")
	s.GET("/foo", func(c Context) Response {
		panic("uh oh")
	})
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	c := RequestContext(req)
	r := s.Handle(c)
	if r.StatusCode != 500 {
		t.Errorf("status code should be 500: %d", r.StatusCode)
	}
}

func TestServePath(t *testing.T) {
	s := newScope("/")
	adapter := NewMEMAdapter()
	FSAdapter = adapter
	afero.WriteFile(adapter.MEMFS, "/public/test.txt", []byte("public file"), 0755)
	s.ServePath("/test", "/public")
	t.Run("valid path", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://example.com/test/test.txt", nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		c := RequestContext(req)
		r := s.Handle(c)
		if r.StatusCode != 200 {
			t.Errorf("status code should be 200: %d", r.StatusCode)
		}
		if r.Filepath != "/test.txt" {
			t.Errorf("filepath not correct: %s", r.Filepath)
		}

		if r.Fileroot != "/public" {
			t.Errorf("fileroot not correct: %s", r.Fileroot)
		}
	})
}

func TestFixPath(t *testing.T) {
	p := "test"
	if fixPath(p) != "/test" {
		t.Errorf("path not prepended with slash: %s", fixPath(p))
	}
}

func TestDeepestRoutePriority(t *testing.T) {
	s := newScope("/")
	s.GET("/*", func(c Context) Response {
		return c.R("catch")
	})
	ss := s.Scope("/api")
	ss.GET("/endpoint", func(c Context) Response {
		return c.R("api")
	})
	req, err := http.NewRequest("GET", "http://example.com/api/endpoint", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	c := RequestContext(req)
	r := s.Handle(c)
	rStr := r.Data.(string)
	if rStr != "api" {
		t.Errorf("should get api response, got %s", rStr)
	}
}
