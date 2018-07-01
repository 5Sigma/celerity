package celeritytest

import (
	"encoding/json"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/tidwall/gjson"
)

// Response is returend when a request is made against the test server It
// contains helper methods to validate the resulting JSON and check things
// like the HTTP status.
type Response struct {
	StatusCode int
	Data       string
	validator  *validator.Validate
}

// AssertString checks a string value in the returning JSON at a given path.
//
// 		r := celeritytest.Get(svr, "/foo")
// 		if ok, v := r.AssertString("data.firstName", "alice"); !ok {
// 			t.Errrof("first name not returned correctly: %s", v)
// 		}
func (r *Response) AssertString(path, value string) (bool, string) {
	v := gjson.Get(r.Data, path)
	if !v.Exists() {
		return false, ""
	}
	if v.String() != value {
		return false, v.String()
	}
	return true, v.String()
}

// AssertBool checks a boolean value in the returning JSON at a given path.
//
// 		r := celeritytest.Get(svr, "/foo")
// 		if ok, _ := r.AssertBool("data.isAdmin", ); !ok {
// 			t.Errrof("admin should be true")
// 		}
func (r *Response) AssertBool(path string, value bool) (bool, bool) {
	v := gjson.Get(r.Data, path)
	if !v.Exists() {
		return false, false
	}
	if v.Bool() != value {
		return false, v.Bool()
	}
	return true, v.Bool()
}

// AssertInt checks an integer value in the returning JSON at a given path.
//
// 		r := celeritytest.Get(svr, "/foo")
// 		if ok, v := r.AssertString("data.age", 19); !ok {
// 			t.Errrof("age was not returned correctly: %d", v)
// 		}
func (r *Response) AssertInt(path string, value int) (bool, int) {
	v := gjson.Get(r.Data, path)
	if !v.Exists() {
		return false, 0
	}
	if v.Int() != int64(value) {
		return false, int(v.Int())
	}
	return true, int(v.Int())
}

// GetLength returns the length of an array in a at a given JSON path.
func (r *Response) GetLength(path string) int {
	return len(gjson.Get(r.Data, path).Array())
}

// Exists checks if a value exists at a given JSON path
func (r *Response) Exists(path string) bool {
	return gjson.Get(r.Data, path).Exists()
}

// Extract unmarshals the JSON into a struct
func (r *Response) Extract(obj interface{}) error {
	return json.Unmarshal([]byte(r.Data), &obj)
}

// ExtractAt Unmarshals JSON at a path into a struct.
func (r *Response) ExtractAt(path string, obj interface{}) error {
	raw := gjson.Get(r.Data, path).Raw
	return json.Unmarshal([]byte(raw), &obj)
}

// GetResult returns a result object at a given path.
func (r *Response) GetResult(path string) gjson.Result {
	return gjson.Get(r.Data, path)
}

// Validate validates the the response data against a validation structure
func (r *Response) Validate(vs interface{}) error {
	if r.validator == nil {
		r.validator = validator.New()
	}
	err := json.Unmarshal([]byte(r.Data), vs)
	if err != nil {
		return err
	}
	return r.validator.Struct(vs)
}
