package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValid(t *testing.T) {
	req := httptest.NewRequest("GET", "/nil", nil)
	form := NewForm(req.PostForm)

	isvalid := form.Valid()
	if !isvalid {
		t.Errorf("Form should be valid")
	}
}

func TestRequire(t *testing.T) {
	req := httptest.NewRequest("get", "/somewhere", nil)
	form := NewForm(req.PostForm)

	form.Require("a", "b", "c")
	isvalid := form.Valid()
	if isvalid {
		t.Errorf("there should be not valid, beacuse a,b,c is empty")
	}

	form = NewForm(req.PostForm)
	values := url.Values{}
	values.Add("a", "a")
	values.Add("b", "b")
	values.Add("c", "c")
	req.PostForm = values
	isvalid = form.Valid()
	if !isvalid {
		t.Errorf("there should be valid,a,b,c is not empty. but get not valid")
	}

}

func TestHas(t *testing.T) {
	postData := url.Values{}
	form := NewForm(postData)

	isHas := form.Has("a")
	if isHas {
		t.Error("does't has a but has a")
	}

	postData.Add("a", " a is a")
	form = NewForm(postData)
	isHas = form.Has("a")
	if !isHas {
		t.Error("has but not has")
	}
}

func TestMimimum(t *testing.T) {
	postData := url.Values{}
	form := NewForm(postData)

	form.Minimum("a", 3)

	postData.Add("a", "aa")
	if form.Valid() {
		t.Error("field length below 3 but got pass")
	}

	postData.Add("b", "bbb")
	form.Minimum("b", 3)
	bError := form.Errors.Get("b")
	if bError != "" {
		t.Error("field lenth is satisfied but got error")
	}

}

func TestIsEmail(t *testing.T) {
	postData := url.Values{}
	form := NewForm(postData)

	postData.Add("success", "a@a.com")
	postData.Add("false", "x")
	postData.Add("false2", "")

	form.IsEmail("success")
	form.IsEmail("false")
	form.IsEmail("false2")

	successError := form.Errors.Get("success")
	if successError != "" {
		t.Error("good email but get error")
	}

	if form.Valid() {
		t.Error("form has error but success")
	}

	false2Error := form.Errors.Get("false2")
	if false2Error == "" {
		t.Error("false2 is error but success")

	}
}
