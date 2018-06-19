package celeritytest

import (
	"github.com/tidwall/gjson"
)

// Response - A response form the test server
type Response struct {
	StatusCode int
	Data       string
}

// AssertString - Assert a string value in the data
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

// AssertBool - Assert a string value in the data
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

// AssertInt - Assert a int value in the data
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
