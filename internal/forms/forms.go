package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func NewForm(values url.Values) *Form {
	return &Form{
		values,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Require(fields ...string) {
	for _, field := range fields {
		value := f.Values.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors[field] = append(f.Errors[field], "This field cannot be blank")
		}
	}
}

func (f *Form) Minimum(field string, limit int) {
	value := f.Get(field)
	if len(value) < limit {
		f.Errors[field] = append(f.Errors[field], fmt.Sprintf(" %s must be at least %d characters long ", field, limit))
	}
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
