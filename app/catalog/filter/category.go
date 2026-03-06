package filter

import (
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/database"
)

// CategoryFilter is a Filter that filters products by a single or multiple
// categories. A user will provide a comma-separated list of category names
// as a query parameter.
//
// An example of such query can be seen below.
//
// <scheme>://<host>:<port>/catalog?category=category1,category2
type CategoryFilter struct {
	Categories []string
}

func (cf *CategoryFilter) Parse(r *http.Request) error {
	data := r.URL.Query().Get("category")

	if len(data) > 0 {
		for _, category := range strings.Split(data, ",") {
			if !cf.Validate("category", category) {
				continue
			}

			cf.Categories = append(cf.Categories, category)
		}
	}

	return nil
}

func (cf *CategoryFilter) Validate(parameter string, value any) bool {
	// assuming the values are URL encoded, but again the Query() function parses
	// any wrongdoing by an ill-intentioned user.
	return len(value.(string)) > 0
}

func (cf *CategoryFilter) Translate() database.Filter {
	return &database.CategoryFilter{
		Categories: cf.Categories,
	}
}
