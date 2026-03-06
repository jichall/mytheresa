package filter

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryValidation(t *testing.T) {
	t.Run("single category filter", func(t *testing.T) {
		pf := &CategoryFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "category=electronics"},
		}

		err := pf.Parse(r)

		assert.NoError(t, err)
		assert.Equal(t, []string{"electronics"}, pf.Categories)
	})

	t.Run("multiple category filter", func(t *testing.T) {
		pf := &CategoryFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "category=electronics,books"},
		}

		err := pf.Parse(r)

		assert.NoError(t, err)
		assert.Equal(t, []string{"electronics", "books"}, pf.Categories)
	})
}
