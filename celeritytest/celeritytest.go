package celeritytest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/5Sigma/celerity"
)

//Post creates a TestServer for the given celerity.Server and makes a POST
//request against it using the TestServer.Post function.
func Post(s *celerity.Server, path string, data []byte) (*Response, error) {
	svr := &TestServer{s}
	return svr.Post(path, data)
}

//Get creates a TestServer for the given celerity.Server and makes a GET
//request against it using the TestServer.Get function.
func Get(s *celerity.Server, path string) (*Response, error) {
	svr := &TestServer{s}
	return svr.Get(path)
}

//Request creates a TestServer for the given celerity.Server and makes a
//request against it using the TestServer.Request function.
func Request(s *celerity.Server, opts RequestOptions) (*Response, error) {
	svr := &TestServer{s}
	return svr.Request(opts)
}

// TestServer can be used to make calls against a managed test
// version of the http server. This is internally used by the Request, Get, and
// Post package level functions.
type TestServer struct {
	Server *celerity.Server
}

// RequestOptions are used by the TestServer.Request function can
// be used with a RequestOptions structure when more advanced request
// customization is needed. Such as configuring headers.
type RequestOptions struct {
	Method string
	Path   string
	Header http.Header
	Data   []byte
}

// Post makes a POST request against the test server. This function is called
// by the package level Post function in cases where you want to make a one off
// request.
func (ts *TestServer) Post(path string, data []byte) (*Response, error) {
	reqOpts := RequestOptions{
		Method: celerity.POST,
		Path:   path,
		Data:   data,
	}
	return ts.Request(reqOpts)
}

// Get - Makes a GET request against the test server. This function is called
// by the package level Get function in cases where you want to make a one off
// request.
func (ts *TestServer) Get(path string) (*Response, error) {
	reqOpts := RequestOptions{
		Method: celerity.GET,
		Path:   path,
	}
	return ts.Request(reqOpts)
}

// Request makes a request against the test server. This function is called by
// the package level Request function for one off requests.  This function can
// be used for more customization when making requests than the Get and Post
// functions provide.
func (ts *TestServer) Request(reqOpts RequestOptions) (*Response, error) {
	httpServer := httptest.NewServer(ts.Server)
	defer httpServer.Close()

	url := fmt.Sprintf("%s%s", httpServer.URL, reqOpts.Path)
	req, err := http.NewRequest(reqOpts.Method, url, bytes.NewBuffer(reqOpts.Data))
	if err != nil {
		return nil, err
	}
	req.Header = reqOpts.Header
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: res.StatusCode,
		Data:       string(bbody),
		Header:     res.Header,
	}
	return response, err
}
