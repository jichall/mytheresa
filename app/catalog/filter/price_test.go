package filter

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPriceValidation(t *testing.T) {
	t.Run("invalid price (price < 0)", func(t *testing.T) {
		pf := &PriceFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "price=-11&operator=gt"},
		}

		err := pf.Parse(r)

		assert.Equal(t, errors.New("parameter price is not valid for the given value -11"), err)
	})

	t.Run("invalid price operator", func(t *testing.T) {
		pf := &PriceFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "price=5.50&operator=BAD"},
		}

		err := pf.Parse(r)

		assert.Equal(t, errors.New("parameter operator is not valid for the given value BAD"), err)
	})

	t.Run("default values", func(t *testing.T) {
		pf := &PriceFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: ""},
		}

		err := pf.Parse(r)

		assert.Nil(t, err)
		assert.Equal(t, decimal.Decimal{}, pf.Price)
		assert.Equal(t, "", pf.Operator)
	})

	t.Run("valid price filter", func(t *testing.T) {
		pf := &PriceFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "price=5.50&operator=gt"},
		}

		err := pf.Parse(r)
		assert.Nil(t, err)

		p, _ := decimal.NewFromString("5.50")

		assert.Equal(t, p, pf.Price)
		assert.Equal(t, "gt", pf.Operator)
	})
}
