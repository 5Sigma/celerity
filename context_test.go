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

	if len(c.RequestID) != 32 {
		t.Errorf("request id invalid: %s", c.RequestID)
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

func TestGet(t *testing.T) {
	c := NewContext()
	c.Set("test", "123")
	if c.Get("test").(string) != "123" {
		t.Errorf("get returned incorrect value")
	}
	if c.Get("1231") != nil {
		t.Error("get for nonexistant key should return nil")
	}

}

func TestBody(t *testing.T) {
	payload := []byte(`
		{
			"foo": {
				"foos": "bar"
			}
		}
	`)
	req, _ := http.NewRequest("GET", "/test", bytes.NewReader(payload))
	c := RequestContext(req)
	v := c.ExtractValue("foo.foos").String()
	if v != "bar" {
		t.Errorf("bar value not found: %s", c.Body())
	}
}

func TestChainResponses(t *testing.T) {
	c := NewContext()
	data := map[string]string{"name": "jim"}
	r := c.Respond(data).Status(201).MetaValue("test", "123")
	if r.StatusCode != 201 {
		t.Errorf("Status code not set to 201 (%d)", r.StatusCode)
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

	if r.Meta["test"].(string) != "123" {
		t.Error("meta not properly set")
	}
}
