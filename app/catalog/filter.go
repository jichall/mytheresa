package catalog

import (
	"fmt"
	"net/http"
	"strconv"
)

type ProductFilter struct {
	Page  int
	Limit int
}

func (pf *ProductFilter) parse(r *http.Request) error {
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

			if !pf.valid(parameter, value) {
				return fmt.Errorf("parameter %s is not valid for the given value %d", parameter, value)
			} else {
				m[parameter] = value
			}
		}
	}

	pf.Page = m["offset"]
	pf.Limit = m["limit"]

	return nil
}

// valid verifies if a given parameter satisfies the rules defined for it
func (pf *ProductFilter) valid(parameter string, value int) bool {
	switch parameter {
	case "offset":
		return value >= 0
	case "limit":
		return value >= 1 && value <= 100
	}

	return false
}
