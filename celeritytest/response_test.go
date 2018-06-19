package celeritytest

import "testing"

var mockResponse = Response{
	StatusCode: 200,
	Data: `
		{
			"success": true,
			"error": "",
			"data": {
				"people": [
					{
						"firstName": "Alice",
						"lastName": "Alisson",
						"age": 19
					},
					{
						"firstName": "Beverly",
						"lastName": "Beaver",
						"age": 21
					}
				]
			}
		}
	`,
}

func TestAssertString(t *testing.T) {
	{
		ok, v := mockResponse.AssertString("data.people.0.firstName", "something bad")
		if ok {
			t.Error("assertion should be false")
		}
		if v != "Alice" {
			t.Errorf("did not respond with current value: %s", v)
		}
	}
	{
		ok, _ := mockResponse.AssertString("data.people.1.firstName", "Beverly")
		if !ok {
			t.Error("assertion should be true")
		}
	}
}

func TestAssertInt(t *testing.T) {
	{
		ok, v := mockResponse.AssertInt("data.people.0.age", 16)
		if ok {
			t.Error("assertion should be false")
		}
		if v != 19 {
			t.Errorf("did not respond with current value: %d", v)
		}
	}
	{
		ok, _ := mockResponse.AssertInt("data.people.0.age", 19)
		if !ok {
			t.Error("assertion should be true")
		}
	}
}

func TestAssertBool(t *testing.T) {
	{
		ok, _ := mockResponse.AssertBool("success", false)
		if ok {
			t.Error("assertion should be false")
		}
	}
	{
		ok, _ := mockResponse.AssertBool("success", true)
		if !ok {
			t.Error("assertion should be true")
		}
	}
}
