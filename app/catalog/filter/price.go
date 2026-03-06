package filter

import (
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/shopspring/decimal"
)

var (
	// keeping some objects from being recreated whenever I'm validating a the price filter.
	operators = map[string]bool{
		"gt":  true,
		"gte": true,
		"lt":  true,
		"lte": true,
		"eq":  true,
	}

	zero = decimal.NewFromInt(0)
)

// PriceFilter is a Filter that filters products by price and an operator. Such
// filter can be used to filter products by price range, e.g. products with a
// price greater than or equal to a certain value.
//
// Example of a query parameter that conforms to this filter can be seen below.
//
//	<scheme>://<host>:<port>/catalog?target=5.50&operator=gt
//
// Possible operators values: gt, gte, lt, lte, eq
//
//	gt = greater than
//	gte = greater than or equal to
//	lt = less than
//	lte = less than or equal to
//	eq = equal to
type PriceFilter struct {
	Price    decimal.Decimal
	Operator string
}

func (pf *PriceFilter) Parse(r *http.Request) error {
	price := r.URL.Query().Get("price")
	operator := r.URL.Query().Get("operator")

	if len(price) <= 0 || len(operator) <= 0 {
		return nil
	}

	if !pf.Validate("price", price) {
		return fmt.Errorf("parameter price is not valid for the given value %s", price)
	}

	if !pf.Validate("operator", operator) {
		return fmt.Errorf("parameter operator is not valid for the given value %s", operator)
	}

	pf.Price, _ = decimal.NewFromString(price)
	pf.Operator = operator

	return nil
}

func (pf *PriceFilter) Validate(parameter string, value any) bool {
	switch parameter {
	case "price":
		price, err := decimal.NewFromString(value.(string))
		if err != nil {
			return false
		}

		return price.GreaterThanOrEqual(zero)
	case "operator":
		if ok, _ := operators[value.(string)]; ok {
			return true
		}
	}

	return false
}

func (pf *PriceFilter) Translate() database.Filter {
	return &database.PriceFilter{
		Price:    pf.Price,
		Operator: pf.Operator,
	}
}
