package celeritytest

import (
	"net/http"
	"testing"

	"github.com/5Sigma/celerity"
)

func TestGet(t *testing.T) {
	s := celerity.New()
	s.Route(celerity.GET, "/foo", func(c celerity.Context) celerity.Response {
		return c.R(map[string]string{"param1": "123"})
	})

	r, err := Get(s, "/foo")
	if err != nil {
		t.Error(err.Error())
	}

	if ok, v := r.AssertString("data.param1", "123"); !ok {
		t.Errorf("param1 was %s", v)
	}
}

func TestPost(t *testing.T) {
	s := celerity.New()
	s.Route(celerity.POST, "/foo", func(c celerity.Context) celerity.Response {
		req := struct {
			Param1 string `json:"param1"`
		}{}
		c.Extract(&req)
		return c.R(req)
	})

	payload := []byte(`
		{
			"param1": "123"
		}
	`)
	r, err := Post(s, "/foo", payload)
	if err != nil {
		t.Error(err.Error())
	}

	if ok, v := r.AssertString("data.param1", "123"); !ok {
		t.Errorf("param1 was %s", v)
	}
}

func TestRequest(t *testing.T) {
	s := celerity.New()
	s.Route(celerity.POST, "/foo", func(c celerity.Context) celerity.Response {
		req := struct {
			Param1 string `json:"param1"`
		}{}
		c.Extract(&req)
		return c.R(req)
	})

	opts := RequestOptions{
		Method: "POST",
		Path:   "/foo",
		Data:   []byte(`{ "param1": "123"}`),
	}
	r, err := Request(s, opts)
	if err != nil {
		t.Error(err.Error())
	}

	if ok, v := r.AssertString("data.param1", "123"); !ok {
		t.Errorf("param1 was %s", v)
	}
}

func TestHeaders(t *testing.T) {
	s := celerity.New()
	s.Route(celerity.POST, "/foo", func(c celerity.Context) celerity.Response {
		req := struct {
			Param1 string `json:"param1"`
		}{Param1: c.Header("Test-Header")}
		return c.R(req)
	})

	opts := RequestOptions{
		Method: "POST",
		Path:   "/foo",
		Data:   []byte(`{ "param1": "123"}`),
		Header: http.Header{"Test-Header": []string{"test-value"}},
	}
	r, err := Request(s, opts)
	if err != nil {
		t.Error(err.Error())
	}

	if ok, v := r.AssertString("data.param1", "test-value"); !ok {
		t.Errorf("param1 was %s", v)
	}
}
