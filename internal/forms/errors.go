package forms

type errors map[string][]string

// Add adds an error message to the field given.
func (e errors) Add(field, message string) {

	e[field] = append(e[field], message)
}

// Get returns the first error message from the field given.
func (e errors) Get(field string) string {

	es := e[field]

	if len(es) == 0 {
		return ""
	}

	return es[0]
}
