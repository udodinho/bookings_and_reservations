package form

import (
	"net/http"
	"net/url"
	"strings"
)

// Form creates a custom form struct, embeds an url.Values object.
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and not empty.
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)
	if formField == "" {
		f.Errors.Add(field, "This field cannot be empty")
		return false
	}
	return true
}

// Valid returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}
	}
}
