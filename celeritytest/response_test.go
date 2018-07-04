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

func TestGetLength(t *testing.T) {
	{
		l := mockResponse.GetLength("data.people")
		if l != 2 {
			t.Errorf("length was not correct: %d", l)
		}
	}
}

func TestExists(t *testing.T) {
	{
		if !mockResponse.Exists("data.people") {
			t.Error("data.people should exist")
		}
		if mockResponse.Exists("data.nothing") {
			t.Error("data.nothing should not exist")
		}
	}
}

func TestExtract(t *testing.T) {
	r := struct {
		Success bool `json:"success"`
	}{}
	mockResponse.Extract(&r)
	if r.Success != true {
		t.Error("values not extracted properly")
	}
}

func TestExtractAt(t *testing.T) {
	r := []struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Age       int    `json:"age"`
	}{}
	mockResponse.ExtractAt("data.people", &r)
	if len(r) != 2 {
		t.Fatal("incorrect number of items in the array")
	}
	if r[0].FirstName != "Alice" {
		t.Error("values not extracted properly")
	}
}

func TestGetResult(t *testing.T) {
	if mockResponse.GetResult("success").Bool() != true {
		t.Error("bool value not returned correctly.")
	}
	if mockResponse.GetResult("data.people.0.firstName").String() != "Alice" {
		t.Error("string value not returned correctly.")
	}
	if mockResponse.GetResult("data.people.1.age").Int() != 21 {
		t.Error("int value not returned correctly.")
	}
}

func TestValidation(t *testing.T) {
	{
		s := struct {
			Success bool   `json:"success" validate:"required"`
			Error   string `json:"error" validate:"isdefault"`
			Data    struct {
				People []struct {
					FirstName string `json:"firstName" validate:"required,alpha"`
					LastName  string `json:"lastName" validate:"required,alpha"`
					Age       int    `json:"age" validate:"numeric,required"`
				} `json:"people" validate:"required,len=2,dive"`
			} `json:"data" validate:"required"`
		}{}
		if err := mockResponse.Validate(&s); err != nil {
			t.Error(err.Error())
		}
	}
	{
		s := struct {
			People []struct {
				FirstName string `json:"firstName" validate:"required,alpha"`
				LastName  string `json:"lastName" validate:"required,alpha"`
				Age       int    `json:"age" validate:"numeric,required"`
			} `json:"people" validate:"required,len=2,dive"`
		}{}
		if err := mockResponse.ValidateAt("data", &s); err != nil {
			t.Error(err.Error())
		}
	}
}
