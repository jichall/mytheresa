package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	handler := NewCatalogHandler(&CatalogHandlerOpts{Repository: &mockrepo{}, Logger: slog.New(slog.NewJSONHandler(io.Discard, nil))})

	t.Run("test get products", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog", nil)
		response := Response{}

		handler.HandleGet(recorder, request)

		decoder := json.NewDecoder(recorder.Body)

		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, response.Products)
		assert.Equal(t, 10, len(response.Products))
		assert.Equal(t, "PRD009", response.Products[9].Code)
	})

	t.Run("test get products with pagination", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=0&limit=3", nil)
		response := Response{}

		handler.HandleGet(recorder, request)

		decoder := json.NewDecoder(recorder.Body)

		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(response.Products))
		assert.Equal(t, "PRD000", response.Products[0].Code)
	})

	t.Run("test get products with category filtering", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		// always fill in the page filter as the mock object I built would panic otherwise
		request := httptest.NewRequest("GET", "/catalog?category=electronics&offset=0&limit=5", nil)
		response := Response{}

		handler.HandleGet(recorder, request)

		decoder := json.NewDecoder(recorder.Body)

		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 5, len(response.Products))

		for _, product := range response.Products {
			assert.Equal(t, "electronics", product.Category)
		}
	})

	t.Run("test get products with wrong offset parameter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=-10&limit=1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("test get products with wrong limit parameter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=0&limit=-1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})
}

// I wish I could use mockery, that would make this test much cleaner
type mockrepo struct{}

func (m *mockrepo) GetAll(ctx context.Context) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{ID: uint(i), Code: fmt.Sprintf("PRD%03d", i), Price: decimal.NewFromInt(int64(i)), Variants: []models.Variant{}}
	}

	return data, nil
}

func (m *mockrepo) GetPaged(ctx context.Context, page, size int) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{ID: uint(i), Code: fmt.Sprintf("PRD%03d", i), Price: decimal.NewFromInt(int64(i)), Variants: []models.Variant{}}
	}

	offset := page * size
	if offset >= len(data) {
		return data[len(data)-1-size : len(data)-1], nil
	}

	return data[offset : offset+size], nil
}

func (m *mockrepo) GetWithFilters(ctx context.Context, filters []database.Filter) ([]models.Product, error) {
	data := make([]models.Product, 30)

	for i := 0; i < 30; i++ {
		data[i] = models.Product{
			ID:         uint(i),
			Code:       fmt.Sprintf("PRD%03d", i),
			Price:      decimal.NewFromInt(int64(i)),
			CategoryID: 0,
		}

		if i >= 20 {
			data[i].CategoryID = 1
			data[i].Category = models.Category{
				ID:   1,
				Code: "electronics",
				Name: "Electronics",
			}
		}
	}

	database.Sort(filters)

	// a bit of type punning in this case, very cumbersome but I already now what
	// filters this endpoint creates
	mf := map[string]database.Filter{}

	for _, f := range filters {
		switch v := f.(type) {
		case *database.CategoryFilter:
			if len(v.Categories) > 0 {
				mf["category"] = v
			}
		case *database.PriceFilter:
			if len(v.Operator) > 0 {
				mf["price"] = v
			}
		case *database.PageFilter:
			mf["page"] = v
		}
	}

	// if a category filter is present, build the logic and filter data
	if category, ok := mf["category"].(*database.CategoryFilter); ok {
		filtered := data[:0]
		allowed := map[string]bool{}
		for _, c := range category.Categories {
			allowed[c] = true
		}

		for _, p := range data {
			if _, ok := allowed[p.Category.Code]; ok {
				filtered = append(filtered, p)
			}
		}

		data = filtered
	}

	page := mf["page"].(*database.PageFilter)

	offset := page.Page * page.Size
	if offset >= len(data) {
		return []models.Product{}, nil
	}

	return data[offset : offset+page.Size], nil
}

func (m *mockrepo) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	return nil, nil
}

func (m *mockrepo) Save(ctx context.Context, data *models.Product) error {
	return nil
}
