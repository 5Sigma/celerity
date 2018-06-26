package celerity

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestNewContext(t *testing.T) {
	c := NewContext()
	if c.URLParams == nil {
		t.Error("URLParams not initialized")
	}
	if c.QueryParams == nil {
		t.Error("QueryParams not initialized")
	}
}

func TestRespond(t *testing.T) {
	c := NewContext()
	data := map[string]string{"name": "jim"}
	r := c.Respond(data)
	if r.StatusCode != 200 {
		t.Errorf("Status code not set to 200 (%d)", r.StatusCode)
	}
	if r.Error != nil {
		t.Errorf("Error not empty: %s", r.Error.Error())
	}
	if v, ok := r.Data.(map[string]string); ok {
		if v["name"] != "jim" {
			t.Errorf("name not set correctly in data: %s", v["name"])
		}
	} else {
		t.Error("response data not properly set")
	}
}

func TestFail(t *testing.T) {
	c := NewContext()
	r := c.Fail(errors.New("some error"))
	if r.Error == nil {
		t.Error("error object not set on response")
	} else if r.Error.Error() != "some error" {
		t.Error("Incorrect error object")
	}
	if r.StatusCode != 500 {
		t.Errorf("Status code not 500: %d", r.StatusCode)
	}
}

func TestFailProduction(t *testing.T) {
	c := NewContext()
	r := c.Fail(errors.New("some error"))
	if r.Error == nil {
		t.Error("error object not set on response")
	} else if r.Error.Error() != "some error" {
		t.Error("Incorrect error object")
	}
	if r.StatusCode != 500 {
		t.Errorf("Status code not 500: %d", r.StatusCode)
	}
}

func TestURLParams(t *testing.T) {
	c := NewContext()
	c.SetParams(map[string]string{"n": "1"})
	if c.URLParams.String("n") != "1" {
		t.Error("string parameter not retrieved correctly")
	}
	if c.URLParams.Int("n") != 1 {
		t.Error("int parameter not retrieved correctly")
	}
	if c.URLParams.String("notfound") != "" {
		t.Error("nonexistant string parameter not retrieved correctly")
	}
	if c.URLParams.Int("notfound") != -1 {
		t.Error("nonexistant int parameter not retrieved correctly")
	}

}

func TestQueryParams(t *testing.T) {
	url, _ := url.Parse("http://example.com?user=1&name=alice")
	c := NewContext()
	c.SetQueryParamsFromURL(url)
	if v := c.QueryParams.String("name"); v != "alice" {
		t.Errorf("name param not correct: %s", v)
	}
	if v := c.QueryParams.Int("user"); v != 1 {
		t.Errorf("name param not correct: %d", v)
	}
}

func TestHeader(t *testing.T) {
	c := NewContext()
	if r, err := http.NewRequest("GET", "/", bytes.NewBuffer(nil)); err != nil {
		t.Fatal(err.Error())
	} else {
		c.Request = r
	}
	c.Request.Header.Set("test-header", "header-value")
	if c.Header("test-header") != "header-value" {
		t.Error("header value incorrect")
	}
}
