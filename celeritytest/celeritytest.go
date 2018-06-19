package celeritytest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/5Sigma/celerity"
)

// TestServer - A test server instance
type TestServer struct {
	Server *celerity.Server
}

// RequestOptions - options for the request to the server
type RequestOptions struct {
	Method string
	Path   string
	Header http.Header
	Data   []byte
}

//Post - Make a GET request against a given server
func Post(s *celerity.Server, path string, data []byte) (*Response, error) {
	svr := &TestServer{s}
	return svr.Post(path, data)
}

//Post - Make a POST request to a test server and return a query response
func (ts *TestServer) Post(path string, data []byte) (*Response, error) {
	httpServer := httptest.NewServer(ts.Server)
	defer httpServer.Close()

	res, err := http.Post(httpServer.URL+path, "application/json", bytes.NewBuffer(data))
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
	}
	return response, err
}

//Get - Make a GET request against a given server
func Get(s *celerity.Server, path string) (*Response, error) {
	svr := &TestServer{s}
	return svr.Get(path)
}

//Get - Make a GET request to a test server and return a query response
func (ts *TestServer) Get(path string) (*Response, error) {
	httpServer := httptest.NewServer(ts.Server)
	defer httpServer.Close()

	res, err := http.Get(httpServer.URL + path)
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
	}
	return response, err
}

// Request - Make a request against the server using a Request object
func Request(s *celerity.Server, opts RequestOptions) (*Response, error) {
	svr := &TestServer{s}
	return svr.Request(opts)
}

// Request - Make a request against the server using a Request object
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
	}
	return response, err
}
