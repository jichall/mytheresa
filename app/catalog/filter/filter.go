package filter

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

// Filter is an interface that defines the methods for parsing and validating
// filter parameters before using them in a database query, just to be safe and
// remove any noise/errors that an attacker might inject into the query parameters.
// It also defines the Translate method for converting the filter into a database
// filter for use in a database query.
type Filter interface {
	Parse(r *http.Request) error
	Validate(parameter string, value any) bool
	Translate() database.Filter
}

type FilterMap map[string]Filter
type FilterList []Filter

func (m FilterMap) List() FilterList {
	var filters []Filter

	for _, f := range m {
		filters = append(filters, f)
	}

	return filters
}

func (l FilterList) Translate() []database.Filter {
	var filters []database.Filter

	for _, f := range l {
		filters = append(filters, f.Translate())
	}

	return filters
}
