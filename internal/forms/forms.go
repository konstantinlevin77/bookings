package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

// Form is a type that embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// NewForm initializes a new form.
func NewForm(data url.Values) *Form {

	var errorsInstance errors = map[string][]string{}

	return &Form{
		data,
		errorsInstance,
	}
}

// Required ranges through all specified fields and adds an error if a field is blank.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can't be blank.")
		}
	}
}

// Has returns whether a field is available or not in the form.
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field can't be blank.")
		return false
	}
	return true
}

// MinLength checks for the minimum length for a field.
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) bool {

	x := f.Get(field)

	if !govalidator.IsEmail(x) {
		f.Errors.Add(field, "Invalid email address")
		return false
	}
	return true
}

// Valid returns true if there are no errors, otherwise false.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
