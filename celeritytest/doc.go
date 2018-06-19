/*
Package celeritytest provides helpers for testing Celerity based applications.

Its purpose is to make it easy to write full integration tests for endpoints
in the server. You can easily make requests against a testing server and get
back a response object which contains various helper methods to parse the
output.


Making Request

Internally celeritytest boots up a test server and and makes a request to
the given endpoint.  The return value is a celeritytest. Response object
which has helper methods to validate the JSON response given from the
server.

		func TestExample(t *testing.T) {
			svr := celerity.New()
			svr.Route(celerity.GET, "/foo", func(c celerity.Context) celerity.Response {
				d := map[string]string{"firstName": "alice"}
				return c.R(d)
			})

			r, _ := celeritytest.Get(svr, "/foo")

			if ok, v := r.AssertString("firstName", "alice"); !ok {
				t.Errorf("first name not valid: %s", v)
			}
		}

Using RequestOptions

For more complicated requests the Request function can be used. This function
accepts a RequestOptions stuct that allows for more configuration.


		func TestExample(t *testing.T) {
			svr := celerity.New()
			svr.Route(celerity.GET, "/foo", func(c celerity.Context) celerity.Response {
				d := map[string]string{"firstName": "alice"}
				return c.R(d)
			})

			opts := RequestOptions {
				Path: "/foo",
				Method: celerity.GET,
				Headers: http.Header(map[string]string{"Authorization": "1234567"}),
			}

			r, _ := celeritytest.Request(svr, opts)

			if ok, v := r.AssertString("firstName", "alice"); !ok {
				t.Errorf("first name not valid: %s", v)
			}
		}

*/
package celeritytest
