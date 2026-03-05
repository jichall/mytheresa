package catalog

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterValidation(t *testing.T) {

	t.Run("invalid offset", func(t *testing.T) {
		pf := &ProductFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "offset=-4&limit=11"},
		}

		err := pf.parse(r)

		assert.Equal(t, errors.New("parameter offset is not valid for the given value -4"), err)
	})

	t.Run("invalid limit", func(t *testing.T) {
		pf := &ProductFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: "offset=0&limit=0"},
		}

		err := pf.parse(r)

		assert.Equal(t, errors.New("parameter limit is not valid for the given value 0"), err)
	})

	t.Run("default values", func(t *testing.T) {
		pf := &ProductFilter{}
		r := &http.Request{
			URL: &url.URL{RawQuery: ""},
		}

		err := pf.parse(r)

		assert.Nil(t, err)
		assert.Equal(t, 0, pf.Page)
		assert.Equal(t, 10, pf.Limit)
	})
}
