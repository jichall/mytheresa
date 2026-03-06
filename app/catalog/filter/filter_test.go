package filter

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFilterTranslation(t *testing.T) {

	t.Run("test price filter translation", func(t *testing.T) {
		pf := &PriceFilter{
			Price:    decimal.NewFromInt(10),
			Operator: "gt",
		}

		expected := &database.PriceFilter{
			Price:    decimal.NewFromInt(10),
			Operator: "gt",
		}

		assert.Equal(t, expected, pf.Translate())
	})

	t.Run("test category filter translation", func(t *testing.T) {
		cf := &CategoryFilter{
			Categories: []string{"electronics"},
		}

		expected := &database.CategoryFilter{
			Categories: []string{"electronics"},
		}

		assert.Equal(t, expected, cf.Translate())
	})

	t.Run("test page filter translation", func(t *testing.T) {
		pf := &PageFilter{
			Page: 1,
			Size: 10,
		}

		expected := &database.PageFilter{
			Page: 1,
			Size: 10,
		}

		assert.Equal(t, expected, pf.Translate())
	})
}
