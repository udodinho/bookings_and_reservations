package form

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid instead of valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("forms shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "whatever", nil)

	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("forms shows invalid when required fields are present")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("z", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("z")
	if isError == "" {
		t.Error("should have error, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some values")
	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("form shows minlength 0f 100 for field when it should not")
	}

	postedValues = url.Values{}
	postedValues.Add("latest_field", "some1234")
	form = New(postedValues)

	form.MinLength("latest_field", 1)
	if !form.Valid() {
		t.Error("form shows minlength for 1 is not met when it should")
	}

	isError = form.Errors.Get("latest_field")
	if isError != "" {
		t.Error("should not have error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("z")
	if form.Valid() {
		t.Error("form shows is email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "jd@gmail.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email for field when it should not")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "jd@gmail")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got a valid email for an invalid email address")
	}
}
