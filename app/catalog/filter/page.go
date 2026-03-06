package filter

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

// PageFilter is a filter for pagination parameters present on a URL query string
type PageFilter struct {
	Page int
	Size int
}

func (pf *PageFilter) Parse(r *http.Request) error {
	parameters := r.URL.Query()
	m := map[string]int{}

	m["offset"] = 0
	m["limit"] = 10

	for parameter, _ := range m {
		if parameters.Has(parameter) {
			value, err := strconv.Atoi(parameters.Get(parameter))
			if err != nil {
				return fmt.Errorf("failed to parse value %v for the parameter %s", parameters.Get(parameter), parameter)
			}

			if !pf.Validate(parameter, value) {
				return fmt.Errorf("parameter %s is not valid for the given value %d", parameter, value)
			} else {
				m[parameter] = value
			}
		}
	}

	pf.Page = m["offset"]
	pf.Size = m["limit"]

	return nil
}

// Validate verifies if a given parameter satisfies the rules defined for it.
func (pf *PageFilter) Validate(parameter string, value any) bool {
	switch parameter {
	case "offset":
		return value.(int) >= 0
	case "limit":
		return value.(int) >= 1 && value.(int) <= 100
	}

	return false
}

func (pf *PageFilter) Translate() database.Filter {
	return &database.PageFilter{
		Page: pf.Page,
		Size: pf.Size,
	}
}
